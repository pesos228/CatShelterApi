package heplers

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
)

func UserIdFromContext(ctx context.Context) (string, bool) {
	userId, ok := loadValueFromClaims(ctx, "user_id")
	if !ok {
		return "", false
	}
	userIdString, ok := userId.(string)
	if !ok {
		return "", false
	}
	return userIdString, true
}

func UserRolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := loadValueFromClaims(ctx, "roles")
	if !ok {
		return nil, false
	}
	rolesSlice, ok := roles.([]interface{})
	if !ok {
		return nil, false
	}
	stringRoles := make([]string, 0, len(rolesSlice))
	for _, role := range rolesSlice {
		roleStr, ok := role.(string)
		if !ok {
			return nil, false
		}
		stringRoles = append(stringRoles, roleStr)
	}

	return stringRoles, true
}

func loadValueFromClaims(ctx context.Context, value string) (interface{}, bool) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, false
	}
	val, ok := claims[value]
	if !ok {
		return nil, false
	}

	return val, true
}
