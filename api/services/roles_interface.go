package services

import (
	"drivers-service/api/dto"
)

type RoleInterface interface {
	Create(role *dto.Role) error
	List() ([]dto.RoleList, error)
}
