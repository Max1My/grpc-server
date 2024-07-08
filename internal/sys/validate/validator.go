package validate

import "context"

type Condition func(ctx context.Context) error

func Validate(ctx context.Context, conds ...Condition) error {
	ve := NewValidationErrors()

	for _, c := range conds {
		err := c(ctx)
		if err != nil {
			if IsValidationError(err) {
				ve.addError(err.Error())
				continue
			}

			return err
		}
	}

	if ve.Messages == nil {
		return nil
	}

	return ve
}

func ValidateID(id int64) Condition {
	return func(ctx context.Context) error {
		if id <= 0 {
			return NewValidationErrors("id must be greater than 0")
		}

		return nil
	}
}

func OtherValidateID(id int64) Condition {
	return func(ctx context.Context) error {
		if id <= 100 {
			return NewValidationErrors("id must be greater than 100")
		}

		return nil
	}
}
