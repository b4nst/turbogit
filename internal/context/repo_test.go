package context

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	dir, err := ioutil.TempDir("", "turbogit-test")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir)

	if err := os.Chdir(dir); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestCurrRepo(t *testing.T) {
	d, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	expect, err := git.PlainInit(d, false)
	if err != nil {
		t.Fatal(err)
	}

	actual, err := currRepo()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expect, actual)
}
