package dto

type Role struct {
	ID          string   `json:"-" bson:"_id,omitempty"`
	Name        string   `json:"name,omitempty" example:"Admin" bson:"name,omitempty"`
	Permissions []string `json:"permissions,omitempty" bson:"permissions,omitempty"`
}

type RoleList struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	Name        string   `json:"name" bson:"name,omitempty"`
	Permissions []string `json:"permissions" bson:"permissions,omitempty"`
}
