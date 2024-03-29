package integrations

import (
	"testing"

	"github.com/b4nst/turbogit/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvidersFrom(t *testing.T) {
	r := test.TestRepo(t)
	defer test.CleanupRepo(t, r)
	c, err := r.Config()
	require.NoError(t, err)
	require.NoError(t, c.SetBool("jira.enabled", true))
	require.NoError(t, c.SetString("jira.username", "alice@ecorp.com"))
	require.NoError(t, c.SetString("jira.token", "supersecret"))
	require.NoError(t, c.SetString("jira.domain", "foo.bar"))
	require.NoError(t, c.SetString("jira.filter", "query"))
	require.NoError(t, c.SetString("gitlab.token", "supersecret"))
	r.Remotes.Create("origin", "git@gitlab.com:namespace/project.git")

	p, err := Issuers(r)
	assert.NoError(t, err)
	assert.Len(t, p, 2)
	assert.IsType(t, JiraProvider{}, p[0])
	assert.IsType(t, GitLabProvider{}, p[1])
}
