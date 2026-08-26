package validation

import "strings"

type fieldError struct {
	field  string
	reason string
}

type Errors []fieldError

type Builder struct {
	errs Errors
}

func (e fieldError) Error() string {
	return e.field + ": " + e.reason
}

func (e Errors) Error() string {
	reasons := make([]string, 0, len(e))
	for _, fe := range e {
		reasons = append(reasons, fe.Error())
	}

	return "validation error: " + strings.Join(reasons, "; ")
}

func (e Errors) Fields() map[string]string {
	fields := make(map[string]string, len(e))
	for _, fe := range e {
		fields[fe.field] = fe.reason
	}

	return fields
}

func (b *Builder) Add(field, reason string) {
	b.errs = append(b.errs, fieldError{field: field, reason: reason})
}

func (b *Builder) Err() error {
	if len(b.errs) == 0 {
		return nil
	}

	return b.errs
}
