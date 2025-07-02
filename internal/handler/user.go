package handler

import (
	"api/catshelter/internal/domain"
	"api/catshelter/internal/handler/dto"
	"api/catshelter/internal/middleware/heplers"
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

func mapRolesToRolesResponse(roles []string) []dto.RoleResponse {
	roleResponses := make([]dto.RoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = dto.RoleResponse{Name: role}
	}
	return roleResponses
}

func mapCatsToCatResponses(cats []*domain.Cat) []dto.CatResponse {
	catResponses := make([]dto.CatResponse, len(cats))
	for i, cat := range cats {
		catResponses[i] = dto.CatResponse{
			Name: cat.Name,
			Age:  cat.Age,
		}
	}
	return catResponses
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
