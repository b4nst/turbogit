package git

import (
	"regexp"
	"testing"

	"github.com/b4nst/turbogit/pkg/test"
	git "github.com/libgit2/git2go/v30"
	"github.com/stretchr/testify/suite"
)

type RemoteTestSuite struct {
	suite.Suite
	Repo *git.Repository
}

func (suite *RemoteTestSuite) SetupSuite() {
	suite.Repo = test.TestRepo(suite.T())
	_, err := suite.Repo.Remotes.Create("origin", "git@alice.com:namespace/project.git")
	suite.NoError(err)
	_, err = suite.Repo.Remotes.Create("fork", "git@bob.com:namespace/project.git")
	suite.NoError(err)
}

func (suite *RemoteTestSuite) TearDownSuite() {
	test.CleanupRepo(suite.T(), suite.Repo)
}

func (suite *RemoteTestSuite) TestParseRemote() {
	// Direct
	u, err := ParseRemote(suite.Repo, "origin", false)
	suite.NoError(err)
	suite.Equal(u.String(), "ssh://git@alice.com/namespace/project.git")

	// No fallback
	u, err = ParseRemote(suite.Repo, "rename", false)
	suite.EqualError(err, "remote 'rename' does not exist")

	// Fallback
	u, err = ParseRemote(suite.Repo, "rename", true)
	suite.NoError(err)
	suite.Equal(u.String(), "ssh://git@bob.com/namespace/project.git")
}

func (suite *RemoteTestSuite) TestRemotes() {
	remotes, err := Remotes(suite.Repo)
	suite.NoError(err)
	suite.Len(remotes, 2)
	urls := make([]string, len(remotes))
	for i, remote := range remotes {
		urls[i] = remote.Url()
	}
	suite.ElementsMatch(urls, []string{"git@alice.com:namespace/project.git", "git@bob.com:namespace/project.git"})
}

func (suite *RemoteTestSuite) TestAnyRemoteMatch() {
	match, err := AnyRemoteMatch(suite.Repo, regexp.MustCompile("namespace/project"))
	suite.NoError(err)
	suite.True(match)
	match, err = AnyRemoteMatch(suite.Repo, regexp.MustCompile("nospace/project"))
	suite.NoError(err)
	suite.False(match)
	match, err = AnyRemoteMatch(suite.Repo, nil)
	suite.NoError(err)
	suite.False(match)
}

func TestRemoteTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteTestSuite))
}
