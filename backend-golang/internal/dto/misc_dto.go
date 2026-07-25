package dto

type AccessRequestResponse struct {
	ID             string  `json:"id"`
	RequesterID    string  `json:"requesterId"`
	RequesterName  string  `json:"requesterName"`
	Resource       string  `json:"resource"`
	AccessLevel    string  `json:"accessLevel"`
	Justification  string  `json:"justification"`
	Status         string  `json:"status"`
	RequestedAt    string  `json:"requestedAt"`
	ApprovedBy     *string `json:"approvedBy,omitempty"`
	ApprovedAt     *string `json:"approvedAt,omitempty"`
	ExpiresAt      *string `json:"expiresAt,omitempty"`
}

type CreateAccessRequest struct {
	RequesterID   string `json:"requesterId" validate:"required"`
	RequesterName string `json:"requesterName" validate:"required"`
	Resource      string `json:"resource" validate:"required"`
	AccessLevel   string `json:"accessLevel" validate:"required"`
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
	CreatedAt string   `json:"createdAt"`
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
	CreatedAt string `json:"createdAt"`
}
