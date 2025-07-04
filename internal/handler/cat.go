package handler

import (
	"api/catshelter/internal/domain"
	"api/catshelter/internal/handler/dto"
	"api/catshelter/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type CatHandler struct {
	catService service.CatService
}

func (c *CatHandler) LonelyCats(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	cats, paginationInfo, err := c.catService.FindLonelyCats(r.Context(), page, pageSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := &dto.CatsPaginatedResponse{
		Data:       mapCatsToCatResponses(cats),
		Pagination: *paginationInfo,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *CatHandler) AddCat(w http.ResponseWriter, r *http.Request) {
	var newCatRequest dto.CatRequest
	err := json.NewDecoder(r.Body).Decode(&newCatRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}
	err = c.catService.AddCat(r.Context(), newCatRequest.Name, int(newCatRequest.Age))
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("New cat successfully created"))
}

func NewCatHandler(catService *service.CatService) *CatHandler {
	return &CatHandler{
		catService: *catService,
	}
}
