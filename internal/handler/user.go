package handler

import (
	"api/catshelter/internal/custom_middleware/heplers"
	"api/catshelter/internal/domain"
	"api/catshelter/internal/handler/dto"
	"api/catshelter/internal/repository"
	"api/catshelter/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	response := mapUserToUserInfoResponse(userWithCats, userRoles)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) AdoptCat(w http.ResponseWriter, r *http.Request) {
	var cat dto.AdoptCatRequest
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

func (h *UserHandler) AboutUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "User id is missing in URL", http.StatusBadRequest)
		return
	}

	user, err := h.userService.FindByIdWithAll(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("User with id '%s' not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("DB error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	userRolesStrings := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		userRolesStrings = append(userRolesStrings, role.Name)
	}

	response := mapUserToUserInfoResponse(user, userRolesStrings)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) AddRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "User id is missing in URL", http.StatusBadRequest)
		return
	}

	var roleName dto.AddRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&roleName); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err := h.userService.AddRole(r.Context(), id, roleName.Name)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("User with id '%s' not found", id), http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrRoleNotFound) {
			http.Error(w, fmt.Sprintf("Role with name '%s' not found", roleName.Name), http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("New role successfully added"))
}

func (h *UserHandler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "User id is missing in URL", http.StatusBadRequest)
		return
	}

	var roleName dto.AddRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&roleName); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	err := h.userService.RemoveRole(r.Context(), id, roleName.Name)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("User with id '%s' not found", id), http.StatusNotFound)
			return
		}
		if errors.Is(err, repository.ErrRoleNotFound) {
			http.Error(w, fmt.Sprintf("Role with name '%s' not found", roleName.Name), http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, domain.ErrCannotRemoveLastRole) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Role successfully removed"))
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}
