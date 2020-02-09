package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/josephn123/rekki/pkg/validators"
)

const (
	apiPathValidate     = "/email/validate"
	requestContentType  = "application/json"
	responseContentType = "application/json"
	contentTypeHeader   = "Content-Type"
)

type EmailAPI struct {
	handler http.Handler
	vd      validators.Validators
}

func NewEmailAPI(vd validators.Validators) http.Handler {
	api := &EmailAPI{
		vd: vd,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").Path(apiPathValidate).HandlerFunc(api.handlePost)

	return router
}

func (a *EmailAPI) handlePost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Header.Get(contentTypeHeader) != requestContentType {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	request := EmailAPIValidationRequest{}
	if err := json.Unmarshal(b, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := a.vd.Validate(request.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := EmailAPIValidationResult{
		Valid:      len(result) > 0,
		Validators: make(map[string]EmailAPIValidationValidatorResult),
	}

	for _, x := range result {
		response.Valid = response.Valid && x.Valid
		response.Validators[x.Name] = EmailAPIValidationValidatorResult{
			Valid:  x.Valid,
			Reason: string(x.Reason),
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
