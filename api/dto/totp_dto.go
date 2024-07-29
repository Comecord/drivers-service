package dto

type TotpRequest struct {
	Issuer      string `json:"issuer,omitempty" example:"comecord.com" bson:"issuer,omitempty"`
	AccountName string `json:"accountName,omitempty" example:"user@comecord.com" bson:"accountName,omitempty"`
}
type TotpCodeVerify struct {
	Code  string `json:"code,omitempty"`
	Email string `json:"-"`
}

type TotpResponse struct {
	SecretKey string `json:"secretKey,omitempty" bson:"secretKey,omitempty"`
	TotpURL   string `json:"totpURL,omitempty" bson:"totpURL,omitempty"`
	QrCode    string `json:"qrCode,omitempty" bson:"qrCode,omitempty"`
}

type TOTP struct {
	Issuer      string `json:"issuer,omitempty" bson:"issuer,omitempty"`
	AccountName string `json:"accountName,omitempty" bson:"accountName,omitempty"`
	SecretKey   string `json:"secretKey,omitempty" bson:"secretKey,omitempty"`
	FileName    string `json:"fileName,omitempty" bson:"fileName,omitempty"`
	TotpURL     string `json:"totpURL,omitempty" bson:"totpURL,omitempty"`
}
