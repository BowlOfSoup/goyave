package validation

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"goyave.dev/goyave/v4/util/typeutil"
)

func TestTimezoneValidator(t *testing.T) {
	t.Run("Constructor", func(t *testing.T) {
		v := Timezone()
		assert.NotNil(t, v)
		assert.Equal(t, "timezone", v.Name())
		assert.True(t, v.IsType())
		assert.False(t, v.IsTypeDependent())
		assert.Empty(t, v.MessagePlaceholders(&ContextV5{}))
	})

	cases := []struct {
		value     any
		wantValue *time.Location
		want      bool
	}{
		{value: "UTC", want: true, wantValue: time.UTC},
		{value: "Europe/Paris", want: true, wantValue: typeutil.Must(time.LoadLocation("Europe/Paris"))},
		{value: "Europe/Paris", want: true, wantValue: typeutil.Must(time.LoadLocation("Europe/Paris"))}, // Second time to check cache
		{value: "America/St_Thomas", want: true, wantValue: typeutil.Must(time.LoadLocation("America/St_Thomas"))},
		{value: "GMT", want: true, wantValue: typeutil.Must(time.LoadLocation("GMT"))},
		{value: typeutil.Must(time.LoadLocation("Europe/Paris")), want: true, wantValue: typeutil.Must(time.LoadLocation("Europe/Paris"))},
		{value: "GMT+2", want: false},
		{value: "UTC+2", want: false},
		{value: "Local", want: false},
		{value: "string", want: false},
		{value: "", want: false},
		{value: 'a', want: false},
		{value: 2, want: false},
		{value: 2.5, want: false},
		{value: []string{"string"}, want: false},
		{value: map[string]any{"a": 1}, want: false},
		{value: true, want: false},
		{value: nil, want: false},
	}

	for _, c := range cases {
		c := c
		t.Run(fmt.Sprintf("Validate_%v_%t", c.value, c.want), func(t *testing.T) {
			v := Timezone()
			ctx := &ContextV5{
				Value: c.value,
			}
			ok := v.Validate(ctx)
			if assert.Equal(t, c.want, ok) && ok {
				assert.Equal(t, c.wantValue, ctx.Value)
			}
		})
	}
}