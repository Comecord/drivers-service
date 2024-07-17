package models

type Auth struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

type AuthTOTP struct {
	Issuer      string `json:"issuer,omitempty" bson:"issuer,omitempty"`
	AccountName string `json:"accountName,omitempty" bson:"accountName,omitempty"`
}

type AuthSecret struct {
	SecretKey string `json:"secretKey,omitempty" bson:"secretKey,omitempty"`
}
