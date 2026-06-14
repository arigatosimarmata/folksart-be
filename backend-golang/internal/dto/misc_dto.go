package dto

type AccessRequestResponse struct {
	ID             string  `json:"id"`
	RequesterID    string  `json:"requester_id"`
	RequesterName  string  `json:"requester_name"`
	Resource       string  `json:"resource"`
	AccessLevel    string  `json:"access_level"`
	Justification  string  `json:"justification"`
	Status         string  `json:"status"`
	RequestedAt    string  `json:"requested_at"`
	ApprovedBy     *string `json:"approved_by,omitempty"`
	ApprovedAt     *string `json:"approved_at,omitempty"`
	ExpiresAt      *string `json:"expires_at,omitempty"`
}

type CreateAccessRequest struct {
	RequesterID   string `json:"requester_id" validate:"required"`
	RequesterName string `json:"requester_name" validate:"required"`
	Resource      string `json:"resource" validate:"required"`
	AccessLevel   string `json:"access_level" validate:"required"`
	Justification string `json:"justification"`
}

type ApproveAccessRequest struct {
	Operator string `json:"operator" validate:"required"`
	Note     string `json:"note"`
}

type RejectAccessRequest struct {
	Operator string `json:"operator" validate:"required"`
	Reason   string `json:"reason" validate:"required"`
}

type NotificationRuleResponse struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Trigger   string   `json:"trigger"`
	Severity  string   `json:"severity"`
	Channels  []string `json:"channels"`
	Active    bool     `json:"active"`
	CreatedAt string   `json:"created_at"`
}

type CreateNotificationRuleRequest struct {
	Name     string   `json:"name" validate:"required"`
	Trigger  string   `json:"trigger" validate:"required"`
	Severity string   `json:"severity" validate:"required"`
	Channels []string `json:"channels" validate:"required"`
	Active   bool     `json:"active"`
}

type NotificationResponse struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Severity  string `json:"severity"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"created_at"`
}
