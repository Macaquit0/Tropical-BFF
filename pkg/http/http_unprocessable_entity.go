package sharedhttp

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func FormatMessage(tag string) string {
	switch tag {
	case "min":
		return "the field requires an minimum of 1 item inside array"
	case "uuid":
		return "invalid UUID, use format: 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'"
	case "required":
		return "field is required, and cannot be left blank or empty"
	case "validateBalanceType":
		return "invalid balanceType. use: 'positive', 'negative', or 'notValidated'"
	default:
		return ""
	}
}

type ValidationError struct {
	Message string              `json:"message"`
	Details []ValidationDetails `json:"details"`
}

type ValidationDetails struct {
	Field       string `json:"field"`
	Value       string `json:"value"`
	Constraint  string `json:"constraint"`
	Description string `json:"description"`
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("%s: %v", v.Message, v.Details)
}

func LowerFirstChar(s string) string {
	if len(s) <= 2 {
		return strings.ToLower(s)
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	return string(r)
}

type HttpUnprocessableEntity struct {
	ErrorCode string              `json:"error_code"`
	Message   string              `json:"message"`
	TracerID  string              `json:"tracer_id"`
	Errors    []ValidationDetails `json:"errors"`
}

func NewHttpUnprocessableEntity(tracerId, errorCode string) *HttpUnprocessableEntity {
	return &HttpUnprocessableEntity{
		Message:   "validation error",
		TracerID:  tracerId,
		ErrorCode: errorCode,
		Errors:    make([]ValidationDetails, 0),
	}
}

func (u *HttpUnprocessableEntity) AddError(vl ValidationDetails) {
	u.Errors = append(u.Errors, vl)
}

func (u *HttpUnprocessableEntity) AddErrorFromField(f validator.FieldError) {
	var vl ValidationDetails
	vl.Field = LowerFirstChar(f.Field())
	vl.Constraint = f.ActualTag()
	vl.Value = fmt.Sprintf("%v", f.Value())
	vl.Description = FormatMessage(f.ActualTag())
	u.Errors = append(u.Errors, vl)
}

func (u *HttpUnprocessableEntity) Direct(field, contraint, value, description string) {
	var vl ValidationDetails
	vl.Field = field
	vl.Constraint = contraint
	vl.Value = value
	vl.Description = description
	u.Errors = append(u.Errors, vl)
}
