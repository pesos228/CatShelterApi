package handler

import (
	"api/catshelter/internal/handler/dto"
	"api/catshelter/internal/repository"
	"api/catshelter/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type AuthHandler struct {
	tokenService service.TokenService
	authService  service.AuthService
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Login, req.Password, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokens, err := h.tokenService.CreateSession(r.Context(), user)
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError)
		return
	}

	h.setAuthCookies(w, tokens)
	w.Write([]byte("Registration successful"))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		http.Error(w, "Incorrect login or password", http.StatusBadRequest)
		return
	}

	tokens, err := h.tokenService.CreateSession(r.Context(), user)
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError)
		return
	}

	h.setAuthCookies(w, tokens)
	w.Write([]byte("Login successful"))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	err = h.tokenService.DeleteRefreshToken(r.Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, repository.ErrRefreshTokenNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("DB error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	h.clearAuthCookies(w)
	w.Write([]byte("Successfully logged out"))
}

func (h *AuthHandler) LogoutEverywhere(w http.ResponseWriter, r *http.Request) {

	//TODO: Надо короче сделать так чтобы из middleware брали id пользователя

	userId, ok := "awdawd", true
	if !ok {
		http.Error(w, "User id not found", http.StatusBadRequest)
		return
	}

	err := h.tokenService.DeleteAllRefreshTokens(r.Context(), userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.clearAuthCookies(w)
	w.Write([]byte("Logout every where successfully"))
}

func (h *AuthHandler) UpdateSession(w http.ResponseWriter, r *http.Request) {
	accessToken, err := r.Cookie("access_token")
	if err == nil {
		if accessToken.Expires.After(time.Now()) {
			http.Error(w, "'access_token' is not expried", http.StatusBadRequest)
			return
		}
	}

	refreshToken, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "'refresh_token' cookie not found", http.StatusBadRequest)
		return
	}

	sessionTokens, err := h.tokenService.UpdateSession(r.Context(), refreshToken.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.setAuthCookies(w, sessionTokens)
	w.Write([]byte("Session successfully updated"))
}

func (h *AuthHandler) clearAuthCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})
}

func (h *AuthHandler) setAuthCookies(w http.ResponseWriter, tokens *service.SessionTokens) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken.Token,
		Expires:  tokens.AccessToken.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken.Token,
		Expires:  tokens.RefreshToken.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}
