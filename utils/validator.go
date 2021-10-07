package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

// ErrorResponse is derived from the error result returned by go-validator
type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func Validate(s interface{}) (errors []*ErrorResponse) {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors = make([]*ErrorResponse, 0, len(validationErrors))
		for _, err := range validationErrors {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return
}

// CheckNoLoopRisk checks if users abuse the shortener
// by shorting the host (e.g. localhost:3000) itself,
// leading to an infinite loop, so don't accept the url for shortening.
// It returns false if there is this sort of risk.
func CheckNoLoopRisk(url, host string, port int) bool {
	// basically this functions removes all the commonly found
	// prefixes from URL such as [http, https, www]
	// then checks of the remaining string is our host itself
	hostWithPort := fmt.Sprintf(host+":%d", port)
	if url == host || url == hostWithPort {
		return false
	}
	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]
	return newURL != host && strings.Split(newURL, ":")[0] != host && newURL != hostWithPort
}
