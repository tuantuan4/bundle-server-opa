package dto

import "opa-test/models"

type PermissionDto struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Method string `json:"method"`
}

func ToPermissionDto(permission models.Permission) PermissionDto {
	return PermissionDto{
		Name:   permission.Name,
		Url:    permission.Url,
		Method: permission.Method,
	}
}
