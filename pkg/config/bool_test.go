package config

import (
	"os"
	"testing"

	git "github.com/libgit2/git2go/v30"
	"github.com/stretchr/testify/suite"
)

type BoolTestSuite struct {
	suite.Suite
	Config *git.Config

	f *os.File
}

func (suite *BoolTestSuite) SetupTest() {
	f, err := os.CreateTemp("", "config")
	suite.NoError(err)
	suite.f = f
	cfg, err := git.OpenOndisk(f.Name())
	suite.NoError(err)

	suite.Config = cfg
	suite.NoError(suite.Config.SetBool("key.bool", true))
}

func (suite *BoolTestSuite) TearDownTest() {
	suite.NoError(os.Remove(suite.f.Name()))
}

func (suite *BoolTestSuite) TestAnyBoolDirective() {
	abd := &AnyBoolDirective{"key.bool", suite.Config}
	suite.NoError(abd.Apply(false))

	abd = &AnyBoolDirective{"no.key", suite.Config}
	suite.EqualError(abd.Apply(false), "no.key is not a boolean")
	suite.NoError(abd.Apply(true))
	v, err := suite.Config.LookupBool("no.key")
	suite.NoError(err)
	suite.False(v)
}

func (suite *BoolTestSuite) TestBoolDirective() {
	bd := &BoolDirective{"key.bool", true, suite.Config}
	suite.NoError(bd.Apply(false))

	bd = &BoolDirective{"key.bool", false, suite.Config}
	suite.EqualError(bd.Apply(false), "key.bool is not a boolean equal to false")
	v, err := suite.Config.LookupBool("key.bool")
	suite.NoError(err)
	suite.True(v)

	suite.NoError(bd.Apply(true))
	v, err = suite.Config.LookupBool("key.bool")
	suite.NoError(err)
	suite.False(v)

	bd = &BoolDirective{"no.key", true, suite.Config}
	suite.EqualError(bd.Apply(false), "no.key is not a boolean equal to true")
	suite.NoError(bd.Apply(true))
	v, err = suite.Config.LookupBool("no.key")
	suite.NoError(err)
	suite.True(v)
}

func TestBoolTestSuite(t *testing.T) {
	suite.Run(t, new(BoolTestSuite))
}
