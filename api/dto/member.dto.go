package dto

import (
	"crm-glonass/data/models"
	"time"
)

type MemberResponse struct {
	Email      string                `json:"email,omitempty"`
	FirstName  string                `json:"firstName,omitempty"`
	LastName   string                `json:"lastName,omitempty"`
	MiddleName string                `json:"middleName,omitempty"`
	Birthday   time.Time             `json:"birthday,omitempty"`
	Phone      string                `json:"phone,omitempty"`
	Location   models.MemberLocation `json:"location,omitempty"`
	Role       []models.MemberRole   `json:"role,omitempty"`
	CreateAt   time.Time             `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt  time.Time             `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
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
	Status       string                `json:"-" default:"wait" bson:"status,omitempty"`
	IsTotp       bool                  `json:"-" default:"false" bson:"isTotp,omitempty"`
	FileQRCode   string                `json:"-"`
	SecretQrCode string                `json:"-"`
	CreateAt     time.Time             `json:"-"`
	UpdatedAt    time.Time             `json:"-"`
}

type MemberUpdate struct {
	ID         string                `json:"id,omitempty"`
	Email      string                `json:"email,omitempty"`
	Password   string                `json:"password,omitempty"`
	FirstName  string                `json:"firstName,omitempty"`
	LastName   string                `json:"lastName,omitempty"`
	MiddleName string                `json:"middleName,omitempty"`
	Birthday   time.Time             `json:"birthday,omitempty"`
	Phone      string                `json:"phone,omitempty"`
	Location   models.MemberLocation `json:"location,omitempty"`
	Role       models.MemberRole     `json:"-"`
	Verified   bool                  `json:"-"`
	CreatedAt  string                `json:"-"`
	UpdatedAt  string                `json:"-"`
}

type MemberAuth struct {
	Email    string `json:"email,omitempty"  example:"user@comecord.com" bson:"email"`
	Password string `json:"password,omitempty" example:"calista78Batista" bson:"password"`
}
