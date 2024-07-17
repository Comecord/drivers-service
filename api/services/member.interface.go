package services

import "drivers-service/api/dto"

type MemberInterface interface {
	Register(createDTO *dto.MemberRegistration) error
	Login(*dto.MemberAuth) (*dto.TokenDetail, error)
	Update(updateDTO *dto.MemberUpdate) (*dto.MemberResponse, error)
}
