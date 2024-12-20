package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/backend/bff-cognito/pkg/errors"
	"github.com/backend/bff-cognito/pkg/helper"
	"github.com/go-playground/validator/v10"
)

func Validate(s any, errorMsg string) error {
	validate := validator.New()

	//err := validate.RegisterValidation("decimal", func(fl validator.FieldLevel) bool {
	//	_, err := decimal.NewFromString(fl.Field().String())
	//	return err == nil
	//})
	//if err != nil {
	//	return errors.NewInternalServerErrorWithError("error on add decimal validation", err)
	//}

	validate.RegisterValidation("cpf", func(fl validator.FieldLevel) bool {
		return helper.ValidateCPF(fl.Field().String())
	})
	validate.RegisterValidation("cnpj", func(fl validator.FieldLevel) bool {
		return helper.ValidateCNPJ(fl.Field().String())
	})
	validate.RegisterValidation("taxId", func(fl validator.FieldLevel) bool {
		return helper.ValidateTaxId(fl.Field().String())
	})
	validate.RegisterValidation("zipCode", func(fl validator.FieldLevel) bool {
		pattern := `^\d{5}-\d{3}$`
		re := regexp.MustCompile(pattern)

		return re.MatchString(fl.Field().String())
	})

	validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		return helper.IsValidEmail(fl.Field().String())
	})

	validate.RegisterValidation("pin", func(fl validator.FieldLevel) bool {
		return helper.IsValidPin(fl.Field().String())
	})

	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return helper.ValidatePassword(fl.Field().String())
	})

	validate.RegisterValidation("share", func(fl validator.FieldLevel) bool {
		return fl.Field().Float() >= 0 && fl.Field().Float() <= 1
	})
	vErr := errors.NewValidationError(errorMsg)
	if err := validate.Struct(s); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errStr := strings.Split(e.Error(), "'")
			fmt.Print(errStr)
			vErr.AddError(strings.ToLower(strings.Split(errStr[1], ".")[1]), errStr[5])
		}
	}

	if len(vErr.Errors) > 0 {
		return vErr
	}

	return nil
}
