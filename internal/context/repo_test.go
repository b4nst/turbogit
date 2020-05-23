package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrRepo(t *testing.T) {
	expect, teardown, err := setUp()
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	actual, err := currRepo()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, actual)
}
