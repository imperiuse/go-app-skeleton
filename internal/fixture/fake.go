package fixture

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/jaswdr/faker"
)

var fake = faker.New()

func GetPointer[T any](val T) *T {
	var v = val

	return &v
}

func MarshalWithIgnoreError(some any) []byte {
	b, _ := jsoniter.Marshal(some)
	return b
}

func UnmarshallWithIgnoreError[T any](data []byte) T {
	var t T
	_ = jsoniter.Unmarshal(data, &t)

	return t
}
