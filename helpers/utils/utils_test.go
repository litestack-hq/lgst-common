package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizePhonenumberOpts(t *testing.T) {
	defaultExpectedOutput := "+2348145070123"
	number1 := "08145070123"
	number2 := "+2348145070123"
	number3 := "0814 507 0123"

	formatedNumber1, err1 := NormalizePhonenumber(NormalizePhonenumberOpts{Phonenumber: number1})
	formatedNumber2, err2 := NormalizePhonenumber(NormalizePhonenumberOpts{Phonenumber: number2})
	formatedNumber3, err3 := NormalizePhonenumber(NormalizePhonenumberOpts{Phonenumber: number3})
	formatedNumber4, err4 := NormalizePhonenumber(NormalizePhonenumberOpts{Phonenumber: number1, CountryCode: "US"})

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)
	assert.Nil(t, err4)
	assert.Equal(t, defaultExpectedOutput, formatedNumber1)
	assert.Equal(t, defaultExpectedOutput, formatedNumber2)
	assert.Equal(t, defaultExpectedOutput, formatedNumber3)
	assert.Equal(t, "+108145070123", formatedNumber4)
}
