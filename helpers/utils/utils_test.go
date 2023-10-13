package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePhonenumberOpts(t *testing.T) {
	number1 := "08145070123"
	number2 := "+2348145070123"

	formatedNumber1, err1 := NormalizePhonenumber(NormalizePhonenumberOpts{Phonenumber: number1})
	formatedNumber2, err2 := NormalizePhonenumber(NormalizePhonenumberOpts{Phonenumber: number2})
	formatedNumber3, err3 := NormalizePhonenumber(NormalizePhonenumberOpts{Phonenumber: number1, CountryCode: "US"})

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
	assert.True(t, formatedNumber1 == formatedNumber2)
	assert.True(t, formatedNumber3 == "+108145070123")
}
