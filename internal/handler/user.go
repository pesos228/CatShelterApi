package handler

import (
	"api/catshelter/internal/custom_middleware/heplers"
	"api/catshelter/internal/handler/dto"
	"api/catshelter/internal/repository"
	"api/catshelter/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func (h *UserHandler) AboutMe(w http.ResponseWriter, r *http.Request) {
	userId, ok := heplers.UserIdFromContext(r.Context())
	if !ok {
		w.Write([]byte("Hello, anonymous!"))
		return
	}
	userWithCats, err := h.userService.FindByIdWithCats(r.Context(), userId)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("User with id '%s' not found", userId), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Internal server error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	userRoles, _ := heplers.UserRolesFromContext(r.Context())

	response := &dto.UserInfoResponse{
		Id:    userWithCats.Id,
		Name:  userWithCats.Name,
		Login: userWithCats.Login,
		Roles: mapRolesToRolesResponse(userRoles),
		Cats:  mapCatsToCatResponses(userWithCats.Cats),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) AdoptCat(w http.ResponseWriter, r *http.Request) {
	var cat dto.ShelterCatRequest
	err := json.NewDecoder(r.Body).Decode(&cat)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %s", err.Error()), http.StatusBadRequest)
		return
	}

	userId, ok := heplers.UserIdFromContext(r.Context())
	if !ok {
		http.Error(w, "User id not found", http.StatusBadRequest)
		return
	}

	err = h.userService.AdoptCat(r.Context(), cat.Id, userId)
	if err != nil {
		if errors.Is(err, repository.ErrCatNotFound) || errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Congratulations, you've adopted a cat"))
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
