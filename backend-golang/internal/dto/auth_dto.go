package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterCustomerRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone"`
}

type TokenResponse struct {
	AccessToken string        `json:"accessToken"`
	TokenType   string        `json:"tokenType"`
	ExpiresIn   int           `json:"expiresIn"`
	User        *UserResponse `json:"user,omitempty"`
}

type SessionResponse struct {
	ID           string `json:"id"`
	Device       string `json:"device"`
	IPAddress    string `json:"ipAddress"`
	Location     string `json:"location"`
	CreatedAt    string `json:"createdAt"`
	LastActiveAt string `json:"lastActiveAt"`
	IsCurrent    bool   `json:"isCurrent"`
}
