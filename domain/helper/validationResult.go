package helper

type ValidationResult struct {
	result map[string][]string
	valid  bool
}

func NewValidationResult() *ValidationResult {
	return &ValidationResult{result: map[string][]string{}, valid: true}
}

func (validationResult *ValidationResult) AddError(field string, err string) {
	validationResult.valid = false
	if _, ok := validationResult.result[field]; !ok {
		validationResult.result[field] = []string{err}
	} else {
		validationResult.result[field] = append(validationResult.result[field], err)
	}
}

func (validationResult *ValidationResult) IsValid() bool {
	return validationResult.valid
}

func (validationResult *ValidationResult) GetAllErrors() []string {
	var errors []string

	for _, v := range validationResult.result {
		errors = append(errors, v...)
	}

	return errors
}
