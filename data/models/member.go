package models

import (
	"time"
)

type MemberLocation struct {
	Address  string `json:"address,omitempty" bson:"address,omitempty"`
	City     string `json:"city,omitempty" bson:"city,omitempty"`
	Postcode string `json:"postcode" bson:"postcode,omitempty"`
	Country  string `json:"country" bson:"country,omitempty"`
}

type MemberRole struct {
	ID          string   `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string   `json:"name,omitempty" bson:"name,omitempty"`
	Permissions []string `json:"permissions,omitempty" bson:"permissions,omitempty"`
}

type Member struct {
	ID           string         `json:"id,omitempty" bson:"_id,omitempty"`
	Email        string         `json:"email,omitempty" bson:"email,omitempty"`
	Password     string         `json:"password,omitempty" bson:"password,omitempty"`
	FirstName    string         `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName     string         `json:"lastName,omitempty" bson:"lastName,omitempty"`
	MiddleName   string         `json:"middleName,omitempty" bson:"middleName,omitempty"`
	Birthday     time.Time      `json:"birthday,omitempty" bson:"birthday,omitempty"`
	Phone        string         `json:"phone,omitempty" bson:"phone,omitempty"`
	Location     MemberLocation `json:"location,omitempty" bson:"location,omitempty"`
	Role         []MemberRole   `json:"role,omitempty" default:"[]" bson:"role,omitempty"`
	Verified     bool           `json:"verified,omitempty" default:"false" bson:"verified,omitempty"`
	IsTotp       bool           `json:"isTotp,omitempty" default:"false" bson:"isTotp,omitempty"`
	FileQRCode   string         `json:"fileQrCode,omitempty" bson:"fileQrCode,omitempty"`
	SecretQrCode string         `json:"secretQrCode,omitempty" bson:"secretQrCode,omitempty"`
	CreateAt     time.Time      `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt    time.Time      `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
