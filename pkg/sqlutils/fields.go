package sqlutils

import (
	"fmt"
	"strings"
)

type FieldsBuilder struct {
	fields []string
	values []string
	args   []any
	count  int
}

func NewFieldsBuilder() *FieldsBuilder {
	return &FieldsBuilder{
		fields: make([]string, 0),
		values: make([]string, 0),
		args:   make([]any, 0),
		count:  0,
	}
}

func (b *FieldsBuilder) AddField(name string, value any, customPlaceholder ...string) {
	b.count++
	b.fields = append(b.fields, name)
	placeholder := "$%d"
	if len(customPlaceholder) > 0 {
		placeholder = customPlaceholder[0]
	}
	b.values = append(b.values, fmt.Sprintf(placeholder, b.count))
	b.args = append(b.args, value)
}

func (b *FieldsBuilder) Args() []any {
	return b.args
}

func (b *FieldsBuilder) Fields() string {
	return strings.Join(b.fields, ", ")
}

func (b *FieldsBuilder) Values() string {
	return strings.Join(b.values, ", ")
}
