package custom_middleware

import (
	"api/catshelter/internal/custom_middleware/heplers"
	"net/http"
	"strings"
)

func RoleRequired(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				userRoles, _ := heplers.UserRolesFromContext(r.Context())

				hasRole := false
				for _, role := range userRoles {
					if strings.EqualFold(role, requiredRole) {
						hasRole = true
						break
					}
				}

				if !hasRole {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}
