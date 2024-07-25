package dto

import (
	"drivers-service/data/models"
	"time"
)

type MemberLocationResponse struct {
	Address  string `json:"address" bson:"address,omitempty"`
	City     string `json:"city" bson:"city,omitempty"`
	Postcode string `json:"postcode" bson:"postcode,omitempty"`
	Country  string `json:"country" bson:"country,omitempty"`
}
type MemberRoleResponse struct {
	Name        string   `json:"name,omitempty" bson:"name,omitempty"`
	Permissions []string `json:"permissions,omitempty" bson:"permissions,omitempty"`
}

type MemberResponse struct {
	ID         string                 `json:"id" bson:"_id,omitempty"`
	FirstName  string                 `json:"firstName,omitempty"`
	LastName   string                 `json:"lastName,omitempty"`
	MiddleName string                 `json:"middleName,omitempty"`
	Birthday   time.Time              `json:"birthday,omitempty"`
	Phone      string                 `json:"phone,omitempty"`
	Location   MemberLocationResponse `json:"location,omitempty"`
	Role       []MemberRoleResponse   `json:"role" bson:"role"`
	CreatedAt  time.Time              `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt,omitempty" bson:"updatedAt"`
}

type MemberRegistration struct {
	ID           string                `json:"-" bson:"_id,omitempty"`
	FirstName    string                `json:"-"`
	LastName     string                `json:"-"`
	MiddleName   string                `json:"-"`
	Birthday     time.Time             `json:"-"`
	Email        string                `json:"email,omitempty"  binding:"min=6,email" example:"user@comecord.com"`
	Password     string                `json:"password,omitempty" binding:"required,password,min=6" example:"calista78Batista"`
	Phone        string                `json:"phone,omitempty"  example:"+7 (999) 999-99-99"`
	Location     models.MemberLocation `json:"-" default:"{}"`
	Role         []models.MemberRole   `json:"-" default:"[]"`
	Verified     bool                  `json:"-"`
	Verification string                `json:"-"`
	Status       string                `json:"-" default:"wait" bson:"status,omitempty"`
	IsTotp       bool                  `json:"-" default:"false" bson:"isTotp,omitempty"`
	FileQRCode   string                `json:"-"`
	SecretQrCode string                `json:"-"`
	CreatedAt    time.Time             `json:"-" bson:"createdAt"`
	UpdatedAt    time.Time             `json:"-"  bson:"updatedAt"`
}

type MemberUpdate struct {
	FirstName  string                 `json:"firstName"  example:"Alexander"`
	LastName   string                 `json:"lastName"  example:"Hunter"`
	MiddleName string                 `json:"middleName"  example:"-"`
	Birthday   time.Time              `json:"birthday"  example:"1972-01-01T00:00:00Z"`
	Location   MemberLocationResponse `json:"location"  bson:"location"`
	UpdatedAt  time.Time              `json:"-"  bson:"updatedAt"`
}

type MemberAuth struct {
	Email    string `json:"email,omitempty"  example:"user@comecord.com" bson:"email"`
	Password string `json:"password,omitempty" example:"calista78Batista" bson:"password"`
	Code     string `json:"code" example:"123456" bson:"code"`
}
