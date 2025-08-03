package utils

import (
	"context"
	"log/slog"

	goValidator "github.com/go-playground/validator/v10"
	goNonStandardValidtors "github.com/go-playground/validator/v10/non-standard/validators"

	"github.com/Archisman-Mridha/x-clone/backend/pkg/assert"
)

func NewValidator(ctx context.Context) *goValidator.Validate {
	validator := goValidator.New(goValidator.WithRequiredStructEnabled())

	err := validator.RegisterValidation("notblank", goNonStandardValidtors.NotBlank)
	assert.AssertErrNil(ctx, err, "Failed registering notblank validator")

	return validator
}

type CustomFieldValidators = map[string]goValidator.Func

func RegisterCustomFieldValidators(
	validator *goValidator.Validate,
	customFieldValidators CustomFieldValidators,
) {
	ctx := context.Background()

	for id, customFieldValidator := range customFieldValidators {
		err := validator.RegisterValidation(id, customFieldValidator, false)
		assert.AssertErrNil(ctx, err,
			"Failed registering custom field validator",
			slog.String("id", id),
		)
	}
}
