package dto

type TokenDetail struct {
	AccessToken            string `json:"accessToken"`
	BearerToken            string `json:"bearerToken"`
	RefreshToken           string `json:"refreshToken"`
	AccessTokenExpireTime  int64  `json:"accessTokenExpireTime"`
	RefreshTokenExpireTime int64  `json:"refreshTokenExpireTime"`
}
