package api

type EmailAPIValidationRequest struct {
	Email string `json:"email"`
}

type EmailAPIValidationResult struct {
	Valid      bool                                         `json:"valid"`
	Validators map[string]EmailAPIValidationValidatorResult `json:"validators"`
}

type EmailAPIValidationValidatorResult struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason,omitempty"`
}
