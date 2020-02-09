package validators

type Validators interface {
	Validate(email string) ([]ValidatorResult, error)
}

type ReasonType string

const (
	ReasonTypeInvalidTLD      ReasonType = "INVALID_TLD"
	ReasonTypeUnableToConnect ReasonType = "UNABLE_TO_CONNECT"
)

type ValidatorResult struct {
	Name   string
	Valid  bool
	Reason ReasonType
}
