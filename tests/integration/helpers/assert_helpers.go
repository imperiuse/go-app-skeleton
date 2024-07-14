package helpers

import (
	"errors"
)

var NotFoundRemoteEntry = errors.New("not found remote entry")

func genericAssert[E, I any](
	entries []E,
	id I,
	compareFunc func(E, I) bool,
) (E, error) {
	if len(entries) == 0 {
		return *new(E), NotFoundRemoteEntry
	}

	for _, entry := range entries {
		if compareFunc(entry, id) {
			return entry, nil
		}
	}

	return *new(E), NotFoundRemoteEntry
}
