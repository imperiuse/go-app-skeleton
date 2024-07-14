package apihelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToEnumValueOrDefault(t *testing.T) {
	type (
		Enum1 string
	)

	const (
		yes Enum1 = "yes"
		no  Enum1 = "no"
	)

	assert.Equal(t, yes, ConvertToEnumValueOrDefault("", []Enum1{yes, no}, yes))
	assert.Equal(t, no, ConvertToEnumValueOrDefault("", []Enum1{yes, no}, no))

	assert.Equal(t, no, ConvertToEnumValueOrDefault("ewvqf", []Enum1{yes, no}, no))
	assert.Equal(t, no, ConvertToEnumValueOrDefault("yes1", []Enum1{yes, no}, no))
	assert.Equal(t, no, ConvertToEnumValueOrDefault("no1", []Enum1{yes, no}, no))

	assert.Equal(t, yes, ConvertToEnumValueOrDefault("yes", []Enum1{yes, no}, yes))
	assert.Equal(t, yes, ConvertToEnumValueOrDefault("yes", []Enum1{yes, no}, no))
	assert.Equal(t, no, ConvertToEnumValueOrDefault("no", []Enum1{yes, no}, yes))
	assert.Equal(t, yes, ConvertToEnumValueOrDefault("yes", []Enum1{yes, no}, no))
}

func TestConvertToEnumValue(t *testing.T) {
	type (
		Enum1 string
	)

	const (
		yes Enum1 = "yes"
		no  Enum1 = "no"
	)

	v, err := ConvertToEnumValue("", []Enum1{yes, no})
	assert.Error(t, err)
	assert.Equal(t, Enum1(""), v)

	v, err = ConvertToEnumValue("ewvqf", []Enum1{yes, no})
	assert.Error(t, err)
	assert.Equal(t, Enum1(""), v)

	v, err = ConvertToEnumValue("yes", []Enum1{yes, no})
	assert.NoError(t, err)
	assert.Equal(t, yes, v)

	v, err = ConvertToEnumValue("no", []Enum1{yes, no})
	assert.NoError(t, err)
	assert.Equal(t, no, v)
}

func Test_getIntFromStr(t *testing.T) {
	assert.Equal(t, 10, GetIntFromStr("", 10, 0, 10))
	assert.Equal(t, 10, GetIntFromStr("", 10, 0, 0))

	assert.Equal(t, 10, GetIntFromStr("10", 1, 0, 20))
	assert.Equal(t, 10, GetIntFromStr("10", 1, 0, 10))

	assert.Equal(t, 1, GetIntFromStr("10", 1, 0, 5))

	assert.Equal(t, -10, GetIntFromStr("-10", 1, -20, 5))

	assert.Equal(t, 1, GetIntFromStr("adfs", 1, 0, 50))
	assert.Equal(t, 1, GetIntFromStr("10.0", 1, 0, 50))
	assert.Equal(t, 1, GetIntFromStr("10.21fr", 1, 0, 50))
	assert.Equal(t, 1, GetIntFromStr("1wf!@F13rvr", 1, 0, 50))
}

func TestConvertToEnumsSliceWithDefault(t *testing.T) {
	type (
		Enum1 string
	)

	const (
		yes Enum1 = "yes"
		no  Enum1 = "no"
	)

	assert.Equal(t, []Enum1{yes}, ConvertToEnumsSliceWithDefault([]string{}, []Enum1{yes, no}, yes))
	assert.Equal(t, []Enum1{no}, ConvertToEnumsSliceWithDefault([]string{}, []Enum1{yes, no}, no))

	assert.Equal(t, []Enum1{no}, ConvertToEnumsSliceWithDefault([]string{"1", "3", "", "yes1"}, []Enum1{yes, no}, no))
	assert.Equal(t, []Enum1{yes}, ConvertToEnumsSliceWithDefault([]string{"1", "3", "", "yes1"}, []Enum1{yes, no}, yes))

	assert.Equal(t, []Enum1{yes}, ConvertToEnumsSliceWithDefault([]string{"yes"}, []Enum1{yes, no}, yes))
	assert.Equal(t, []Enum1{yes}, ConvertToEnumsSliceWithDefault([]string{"yes"}, []Enum1{yes, no}, no))
	assert.Equal(t, []Enum1{no}, ConvertToEnumsSliceWithDefault([]string{"no"}, []Enum1{yes, no}, yes))
	assert.Equal(t, []Enum1{yes}, ConvertToEnumsSliceWithDefault([]string{"yes"}, []Enum1{yes, no}, no))

	assert.Equal(t, []Enum1{yes, no}, ConvertToEnumsSliceWithDefault([]string{"yes", "no"}, []Enum1{yes, no}, yes))
	assert.Equal(t, []Enum1{yes, no}, ConvertToEnumsSliceWithDefault([]string{"yes", "123", "no"}, []Enum1{yes, no}, yes))

	assert.Equal(t, []Enum1{yes}, ConvertToEnumsSliceWithDefault([]string{"yes", "123", "no"}, []Enum1{yes}, yes))
	assert.Equal(t, []Enum1{no}, ConvertToEnumsSliceWithDefault([]string{"yes", "123", "no"}, []Enum1{no}, yes))

	assert.Equal(t, []Enum1{yes}, ConvertToEnumsSliceWithDefault([]string{}, []Enum1{no}, yes))
}

func TestConvertToEnumsSlice(t *testing.T) {
	type (
		Enum1 string
	)

	const (
		yes Enum1 = "yes"
		no  Enum1 = "no"
	)

	assert.Equal(t, []Enum1{}, ConvertToEnumsSlice([]string{}, []Enum1{yes, no}))
	assert.Equal(t, []Enum1{}, ConvertToEnumsSlice([]string{"1", "3", "", "yes1"}, []Enum1{yes, no}))

	assert.Equal(t, []Enum1{yes}, ConvertToEnumsSlice([]string{"yes"}, []Enum1{yes, no}))
	assert.Equal(t, []Enum1{no}, ConvertToEnumsSlice([]string{"no"}, []Enum1{yes, no}))

	assert.Equal(t, []Enum1{yes, no}, ConvertToEnumsSlice([]string{"yes", "no"}, []Enum1{yes, no}))
	assert.Equal(t, []Enum1{yes, no}, ConvertToEnumsSlice([]string{"yes", "123", "no"}, []Enum1{yes, no}))

	assert.Equal(t, []Enum1{yes}, ConvertToEnumsSlice([]string{"yes", "123", "no"}, []Enum1{yes}))
	assert.Equal(t, []Enum1{no}, ConvertToEnumsSlice([]string{"yes", "123", "no"}, []Enum1{no}))

	assert.Equal(t, []Enum1{}, ConvertToEnumsSlice([]string{}, []Enum1{no}))
}
