package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	ExpiresIn   int           `json:"expires_in"`
	User        *UserResponse `json:"user,omitempty"`
}

type SessionResponse struct {
	ID           string `json:"id"`
	Device       string `json:"device"`
	IPAddress    string `json:"ip_address"`
	Location     string `json:"location"`
	CreatedAt    string `json:"created_at"`
	LastActiveAt string `json:"last_active_at"`
	IsCurrent    bool   `json:"is_current"`
}
