package dto

type JwtPairDTO struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
