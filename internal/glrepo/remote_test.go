package glrepo

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/profclems/glab/internal/git"
	"github.com/stretchr/testify/assert"
)

func eq(t *testing.T, got interface{}, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestFindByName(t *testing.T) {
	list := Remotes{
		&Remote{Remote: &git.Remote{Name: "mona"}, Repo: New("monalisa", "myfork")},
		&Remote{Remote: &git.Remote{Name: "origin"}, Repo: New("monalisa", "octo-cat")},
		&Remote{Remote: &git.Remote{Name: "upstream"}, Repo: New("hubot", "tools")},
	}

	r, err := list.FindByName("upstream", "origin")
	eq(t, err, nil)
	eq(t, r.Name, "upstream")

	r, err = list.FindByName("nonexist", "*")
	eq(t, err, nil)
	eq(t, r.Name, "mona")

	_, err = list.FindByName("nonexist")
	eq(t, err, errors.New(`no GitLab remotes found`))
}

func TestTranslateRemotes(t *testing.T) {
	publicURL, _ := url.Parse("https://gitlab.com/monalisa/hello")
	originURL, _ := url.Parse("http://example.com/repo")

	gitRemotes := git.RemoteSet{
		&git.Remote{
			Name:     "origin",
			FetchURL: originURL,
		},
		&git.Remote{
			Name:     "public",
			FetchURL: publicURL,
		},
	}

	identityURL := func(u *url.URL) *url.URL {
		return u
	}
	result := TranslateRemotes(gitRemotes, identityURL)

	if len(result) != 1 {
		t.Errorf("got %d results", len(result))
	}
	if result[0].Name != "public" {
		t.Errorf("got %q", result[0].Name)
	}
	if result[0].RepoName() != "hello" {
		t.Errorf("got %q", result[0].RepoName())
	}
}

func Test_remoteNameSortingScore(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output int
	}{
		{
			name:   "upstream",
			input:  "upstream",
			output: 3,
		},
		{
			name:   "gitlab",
			input:  "gitlab",
			output: 2,
		},
		{
			name:   "origin",
			input:  "origin",
			output: 1,
		},
		{
			name:   "else",
			input:  "anyOtherName",
			output: 0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := remoteNameSortScore(tC.input)
			assert.Equal(t, tC.output, got)
		})
	}
}
