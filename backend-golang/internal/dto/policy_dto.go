package dto

type PolicyResponse struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Condition   interface{} `json:"condition"`
	Action      string      `json:"action"`
	Priority    int         `json:"priority"`
	Active      bool        `json:"active"`
	CreatedAt   string      `json:"created_at"`
}

type CreatePolicyRequest struct {
	Name        string      `json:"name" validate:"required"`
	Description string      `json:"description"`
	Condition   interface{} `json:"condition"`
	Action      string      `json:"action" validate:"required"`
	Priority    int         `json:"priority"`
	Active      bool        `json:"active"`
}

type PolicyEvaluationResponse struct {
	UserID        string          `json:"user_id"`
	Resource      string          `json:"resource"`
	Action        string          `json:"action"`
	Decision      string          `json:"decision"`
	MatchedPolicy *PolicyResponse `json:"matched_policy,omitempty"`
}

type EvaluatePolicyRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	Resource string `json:"resource" validate:"required"`
	Action   string `json:"action" validate:"required"`
}
