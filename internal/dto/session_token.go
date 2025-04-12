package dto

type SessionToken struct {
	Username      string `json:"username"`
	IsOtpVerified bool   `json:"isOtpVerified"`
}
