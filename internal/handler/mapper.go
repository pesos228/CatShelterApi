package handler

import (
	"api/catshelter/internal/domain"
	"api/catshelter/internal/handler/dto"
)

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
			Id:   cat.Id,
			Name: cat.Name,
			Age:  cat.Age,
		}
	}
	return catResponses
}

func mapUserToUserInfoResponse(user *domain.User, roles []string) *dto.UserInfoResponse {
	return &dto.UserInfoResponse{
		Id:    user.Id,
		Name:  user.Name,
		Login: user.Login,
		Roles: mapRolesToRolesResponse(roles),
		Cats:  mapCatsToCatResponses(user.Cats),
	}
}
