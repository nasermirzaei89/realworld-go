package slug_test

import (
	"github.com/nasermirzaei89/realworld-go/pkg/slug"
	"testing"
)

func TestMake(t *testing.T) {
	tt := map[string]string{
		"Hello World":       "hello-world",
		"__Hello--World___": "hello-world",
	}

	for tc, expected := range tt {
		res := slug.Make(tc)
		if res != expected {
			t.Errorf("expected '%s', but got '%s'", expected, res)
		}
	}
}
