package integrations

import (
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGitLabProvider(t *testing.T) {
	t.Skip("TODO")
	// r := test.TestRepo(t)
	// defer test.CleanupRepo(t, r)
	// test.InitRepoConf(t, r)
	// r.Remotes.Create("blank", "git@blank.com:project.git")

	// // nil, nil when not in a GitLab repo
	// provider, err := NewGitLabProvider(r)
	// assert.NoError(t, err)
	// assert.Nil(t, provider)

	// r.Remotes.Create("origin", "git@gitlab.com:namespace/project.git")
	// // no token
	// provider, err = NewGitLabProvider(r)
	// assert.EqualError(t, err, "config value 'gitlab.token' was not found")

	// c, err := r.Config()
	// require.NoError(t, err)
	// require.NoError(t, c.SetString("gitlab.token", "supersecret"))
	// // default values
	// provider, err = NewGitLabProvider(r)
	// assert.NoError(t, err)
	// assert.IsType(t, &GitLabProvider{}, provider)
	// assert.Equal(t, "namespace/project", provider.project)
}

func TestGitLabSearch(t *testing.T) {
	t.Skip("TODO")
	// r := test.TestRepo(t)
	// defer test.CleanupRepo(t, r)
	// test.InitRepoConf(t, r)
	// r.Remotes.Create("origin", "git@gitlab.com:namespace/project.git")
	// c, err := r.Config()
	// require.NoError(t, err)
	// require.NoError(t, c.SetString("gitlab.token", "supersecret"))

	// ts := gitlabMockServer(t, "myproject")
	// defer ts.Close()

	// client, err := gitlab.NewClient("supersecret", gitlab.WithBaseURL(ts.URL), gitlab.WithHTTPClient(ts.Client()))
	// require.NoError(t, err)
	// provider := GitLabProvider{
	// 	project: "myproject",
	// 	client:  client,
	// }
	// ids, err := provider.Search()
	// assert.NoError(t, err)
	// assert.Len(t, ids, 1)
	// assert.Equal(t, IssueDescription{
	// 	ID:          "1",
	// 	Name:        "Ut commodi ullam eos dolores perferendis nihil sunt.",
	// 	Description: "Omnis vero earum sunt corporis dolor et placeat.",
	// 	Provider:    GITLAB_PROVIDER,
	// }, ids[0])
}

func gitlabMockServer(t *testing.T, project string) *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v4/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc(path.Join("/api/v4/projects/", project, "issues"), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "assigned_to_me", r.URL.Query().Get("scope"))

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(gitlabIssue))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Unexpected path '%s'", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	})

	return httptest.NewServer(mux)
}

const gitlabIssue = `
[
   {
      "project_id" : 4,
      "milestone" : {
         "due_date" : null,
         "project_id" : 4,
         "state" : "closed",
         "description" : "Rerum est voluptatem provident consequuntur molestias similique ipsum dolor.",
         "iid" : 3,
         "id" : 11,
         "title" : "v3.0",
         "created_at" : "2016-01-04T15:31:39.788Z",
         "updated_at" : "2016-01-04T15:31:39.788Z"
      },
      "author" : {
         "state" : "active",
         "web_url" : "https://gitlab.example.com/root",
         "avatar_url" : null,
         "username" : "root",
         "id" : 1,
         "name" : "Administrator"
      },
      "description" : "Omnis vero earum sunt corporis dolor et placeat.",
      "state" : "closed",
      "iid" : 1,
      "assignees" : [{
         "avatar_url" : null,
         "web_url" : "https://gitlab.example.com/lennie",
         "state" : "active",
         "username" : "lennie",
         "id" : 9,
         "name" : "Dr. Luella Kovacek"
      }],
      "assignee" : {
         "avatar_url" : null,
         "web_url" : "https://gitlab.example.com/lennie",
         "state" : "active",
         "username" : "lennie",
         "id" : 9,
         "name" : "Dr. Luella Kovacek"
      },
      "labels" : ["foo", "bar"],
      "upvotes": 4,
      "downvotes": 0,
      "merge_requests_count": 0,
      "id" : 41,
      "title" : "Ut commodi ullam eos dolores perferendis nihil sunt.",
      "updated_at" : "2016-01-04T15:31:46.176Z",
      "created_at" : "2016-01-04T15:31:46.176Z",
      "closed_at" : "2016-01-05T15:31:46.176Z",
      "closed_by" : {
         "state" : "active",
         "web_url" : "https://gitlab.example.com/root",
         "avatar_url" : null,
         "username" : "root",
         "id" : 1,
         "name" : "Administrator"
      },
      "user_notes_count": 1,
      "due_date": "2016-07-22",
      "web_url": "http://gitlab.example.com/my-group/my-project/issues/1",
      "references": {
        "short": "#1",
        "relative": "#1",
        "full": "my-group/my-project#1"
      },
      "time_stats": {
         "time_estimate": 0,
         "total_time_spent": 0,
         "human_time_estimate": null,
         "human_total_time_spent": null
      },
      "has_tasks": true,
      "task_status": "10 of 15 tasks completed",
      "confidential": false,
      "discussion_locked": false,
      "_links":{
         "self":"http://gitlab.example.com/api/v4/projects/4/issues/41",
         "notes":"http://gitlab.example.com/api/v4/projects/4/issues/41/notes",
         "award_emoji":"http://gitlab.example.com/api/v4/projects/4/issues/41/award_emoji",
         "project":"http://gitlab.example.com/api/v4/projects/4"
      },
      "task_completion_status":{
         "count":0,
         "completed_count":0
      }
   }
]
`
