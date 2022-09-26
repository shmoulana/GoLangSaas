package dto

type SignInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
