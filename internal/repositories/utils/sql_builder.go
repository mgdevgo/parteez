package utils

import "strings"

type SQLBuilder struct {
	builder strings.Builder
	params  []any
}

func (b *SQLBuilder) WriteLine(s string, params ...any) {
	// paramsCount := 0
	// currentParamsCount := len(b.params)
	b.builder.WriteString(s)
	b.builder.WriteRune('\n')
	b.params = append(b.params, params...)
}

func (b *SQLBuilder) Params() []any {
	return b.params
}

func (b *SQLBuilder) String() string {
	return b.builder.String()
}
