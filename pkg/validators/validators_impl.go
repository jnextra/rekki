package validators

import (
	"fmt"
	"net"
	"regexp"
	"sync"
)

type validatorsImpl struct{}

var emailRegex = regexp.MustCompile("(\\S+)@(\\S+)$")

const (
	regexpValidatorName = "regexp"
	domainValidatorName = "domain"
	smtpValidatorName   = "smtp"
	mxValidatorName     = "mx"
	numValidators       = 4
	smtpPortNumber      = 25
)

type validatorResults struct {
	Regexp validatorStageResult
	Domain validatorStageResult
	MX     validatorStageResult
	SMTP   validatorStageResult
}

type validatorStageResult struct {
	Valid  bool
	Reason ReasonType
}

func NewValidators() Validators {
	return &validatorsImpl{}
}

func (v *validatorsImpl) Validate(email string) ([]ValidatorResult, error) {
	output, err := validate(email)
	if err != nil {
		return []ValidatorResult{}, err
	}

	results := make([]ValidatorResult, numValidators)
	results[0] = mapValidatorResult(regexpValidatorName, output.Regexp)
	results[1] = mapValidatorResult(domainValidatorName, output.Domain)
	results[2] = mapValidatorResult(mxValidatorName, output.MX)
	results[3] = mapValidatorResult(smtpValidatorName, output.SMTP)

	return results, nil
}

func mapValidatorResult(name string, stage validatorStageResult) ValidatorResult {
	return ValidatorResult{
		Name:   name,
		Valid:  stage.Valid,
		Reason: stage.Reason,
	}
}

func validate(email string) (validatorResults, error) {
	result := validatorResults{}
	emailParts := emailRegex.FindStringSubmatch(email)
	if len(emailParts) < 3 {
		return result, nil
	}

	result.Regexp.Valid = true

	records, err := net.LookupMX(emailParts[2])
	if err != nil {
		result.Domain.Reason = ReasonTypeInvalidTLD
		return result, nil
	}

	result.Domain.Valid = true

	if len(records) == 0 {
		return result, nil
	}

	result.MX.Valid = true

	wg := sync.WaitGroup{}
	wg.Add(len(records))
	smtpResultsChan := make(chan bool, len(records))

	for _, x := range records {
		go func(host string) {
			smtpResultsChan <- checkSMTP(host)
			wg.Done()
		}(x.Host)
	}

	wg.Wait()

	for range records {
		result.SMTP.Valid = result.SMTP.Valid || <-smtpResultsChan
	}

	if !result.SMTP.Valid {
		result.SMTP.Reason = ReasonTypeUnableToConnect
	}

	return result, nil
}

func checkSMTP(host string) bool {
	_, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, smtpPortNumber))

	return err == nil
}
