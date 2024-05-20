package fileds

import (
	"fmt"
	"strings"
)

type Fields struct {
	fields []*field
}

type field struct {
	Name        string
	Value       any
	Placeholder string
}

func (f *Fields) AddField(name string, value any, placeholder ...string) {
	field := &field{
		Name:        name,
		Value:       value,
		Placeholder: "$%d",
	}

	if len(placeholder) > 0 {
		field.Placeholder = placeholder[0]
	}

	f.fields = append(f.fields, field)
}

func (f *Fields) Build() (string, string, []any) {
	fields := make([]string, len(f.fields))
	values := make([]string, len(f.fields))
	args := make([]any, len(f.fields))
	for i, field := range f.fields {
		fields = append(fields, field.Name)
		values = append(values, fmt.Sprintf(field.Placeholder, i+1))
		args = append(args, field.Value)
	}

	return strings.Join(fields, " "), strings.Join(values, " "), args
}
