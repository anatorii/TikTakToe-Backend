package auth

type AuthResponse struct {
	Success bool   `json:"success"`
	UUID    string `json:"uuid"`
	Login   string `json:"login,omitempty"`
	Message string `json:"message,omitempty"`
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JwtRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JwtResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshJwtRequest struct {
	RefreshToken string `json:"refresh_token"`
}
