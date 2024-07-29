package dto

import (
	"time"
)

type CreateVehicleRequestDTO struct {
	Name     string `json:"name,omitempty" bson:"make,omitempty"`
	Model    string `json:"model,omitempty" bson:"model,omitempty"`
	Year     int    `json:"year,omitempty" bson:"year,omitempty"`
	Price    int    `json:"price,omitempty" bson:"price,omitempty"`
	Location string `json:"location,omitempty" bson:"location,omitempty"`
}
type DBVehicleDTO struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name,omitempty" bson:"make,omitempty"`
	Model     string    `json:"model,omitempty" bson:"model,omitempty"`
	Year      int       `json:"year,omitempty" bson:"year,omitempty"`
	Price     int       `json:"price,omitempty" bson:"price,omitempty"`
	Location  string    `json:"location,omitempty" bson:"location,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type CreateVehicleRequest struct {
	ID        string    `json:"-" bson:"_id,omitempty"`
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	Model     string    `json:"model,omitempty" bson:"model,omitempty"`
	Year      int       `json:"year,omitempty" bson:"year,omitempty"`
	Price     int       `json:"price,omitempty" bson:"price,omitempty"`
	Location  string    `json:"location,omitempty" bson:"location,omitempty"`
	CreatedAt time.Time `json:"-" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"-" bson:"updated_at,omitempty"`
}

type DBVehicle struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	Model     string    `json:"model,omitempty" bson:"model,omitempty"`
	Year      int       `json:"year,omitempty" bson:"year,omitempty"`
	Price     int       `json:"price,omitempty" bson:"price,omitempty"`
	Location  string    `json:"location,omitempty" bson:"location,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type UpdateVehicleRequest struct {
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	Model     string    `json:"model,omitempty" bson:"model,omitempty"`
	Year      int       `json:"year,omitempty" bson:"year,omitempty"`
	Price     int       `json:"price,omitempty" bson:"price,omitempty"`
	Location  string    `json:"location,omitempty" bson:"location,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
