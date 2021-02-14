package integrations

import (
	"io/ioutil"
	"testing"

	git "github.com/libgit2/git2go/v30"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJiraProvider(t *testing.T) {
	cf, err := ioutil.TempFile("", "gitconfig")
	require.NoError(t, err)
	c, err := git.OpenOndisk(cf.Name())
	require.NoError(t, err)

	p, err := jiraProvider(c)
	assert.NoError(t, err)
	assert.Nil(t, p)

	// No username
	err = c.SetBool("jira.enable", true)
	require.NoError(t, err)
	p, err = jiraProvider(c)
	assert.EqualError(t, err, "config value 'jira.username' was not found")
	assert.Nil(t, p)

	// No token
	err = c.SetString("jira.username", "alice@ecorp.com")
	require.NoError(t, err)
	p, err = jiraProvider(c)
	assert.EqualError(t, err, "config value 'jira.token' was not found")
	assert.Nil(t, p)

	// No domain
	err = c.SetString("jira.token", "supersecret")
	require.NoError(t, err)
	p, err = jiraProvider(c)
	assert.EqualError(t, err, "config value 'jira.domain' was not found")
	assert.Nil(t, p)

	// No filter
	err = c.SetString("jira.domain", "foo.bar")
	require.NoError(t, err)
	p, err = jiraProvider(c)
	assert.EqualError(t, err, "config value 'jira.filter' was not found")
	assert.Nil(t, p)

	// All's ok
	err = c.SetString("jira.filter", "query filter")
	require.NoError(t, err)
	p, err = jiraProvider(c)
	assert.NoError(t, err)
	assert.IsType(t, &JiraProvider{}, p)
}
