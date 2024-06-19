package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/huzaifac137/simple_bank/util"
)

// Custom age validator function
func customCurrencyValidator(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	if crncy := util.IsValidCurrency(currency); crncy {
		return true
	}

	return false
}
