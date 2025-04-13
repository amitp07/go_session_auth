package dto

type MfaSession struct {
	Username string `json:"username"`
	Otp      string `json:"otp"`
}
