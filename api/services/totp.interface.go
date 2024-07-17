package services

import "crm-glonass/api/dto"

type TotpInterface interface {
	GenerateTotp(payload *dto.TotpRequest) (*dto.TotpResponse, error)
}
