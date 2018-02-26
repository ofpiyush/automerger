package automerger

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	ErrUnknownEvent  = errors.New("This is not the event we are looking for")
	ErrUnknownBranch = errors.New("This is not the branch we are looking for")
	ErrBadSignature  = errors.New("This is not the signature we are looking for")
)

func fatalErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
func httpErrors(w http.ResponseWriter, code int, errs ...error) bool {
	var errorResp = &ErrorResp{
		Errors: handleErrs(errs...),
	}
	if len(errorResp.Errors) > 0 {
		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(errorResp); err != nil {
			log.Println(err)
		}
		return true
	}
	return false
}

func handleErrs(errs ...error) []error {
	resp := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			// Log to the console anyway
			log.Println(err)
			resp = append(resp, err)
		}
	}

	return resp
}
