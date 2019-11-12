package id_test

import (
	"github.com/nasermirzaei89/realworld-go/libs/id"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("String Length", func(t *testing.T) {
		for i := range []int{1, 2, 3, 5, 8} {
			res := id.New(i)
			if len(res) != i {
				t.Errorf("should return string with length '%d', but got length '%d'", i, len(res))
			}
		}
	})
}
