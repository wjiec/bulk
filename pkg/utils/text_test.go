package utils

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestCamelCaseToHyphen(t *testing.T) {
	tests := []struct{ Camel, Expected string }{
		{"HelloWorld", "hello-world"},
		{"ILoveYou", "i-love-you"},
		{"JSON", "j-s-o-n"},
		{"Hello-World", "hello-world"},
	}

	for _, test := range tests {
		assert.Equal(t, CamelCaseToHyphen(test.Camel), test.Expected)
	}
}
