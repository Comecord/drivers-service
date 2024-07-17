package services

import "drivers-service/api/dto"

type TotpInterface interface {
	GenerateTotp(payload *dto.TotpRequest) (*dto.TotpResponse, error)
}
