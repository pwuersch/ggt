package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/pwuersch/ggt/lib/globals"
)

type GithubProject struct {
	Name     string `json:"name"`
	SshUrl   string `json:"ssh_url"`
	CloneUrl string `json:"clone_url"`
}

type GitHubAdapter struct{}

func (GitHubAdapter) GetName() string {
	return "github"
}

func (GitHubAdapter) Info() AdapterInfo {
	return AdapterInfo{
		Name:           "github",
		RequiresHost:   false,
		ProvidesSearch: false,
		ProvidesScope:  false,
	}
}

func (github GitHubAdapter) GetProjects(profile globals.Profile, _ string, _ string) ([]Project, error) {
	token := profile.Api.Token

	client := &http.Client{}
	url := "https://api.github.com/user/repos?per_page=100&sort=updated"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Project{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	globals.Debug(fmt.Sprintf("Getting %s", url))
	res, err := client.Do(req)
	if err != nil {
		return []Project{}, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []Project{}, err
	}

	githubProjects := []GithubProject{}
	if err = json.Unmarshal(data, &githubProjects); err != nil {
		return []Project{}, err
	}

	return globals.Map(githubProjects, func(githubProject GithubProject) Project {
		return github.mapProject(githubProject)
	}), nil
}

func (GitHubAdapter) CreateProject(profile globals.Profile) error {
	return errors.New("not implemented")
}

func (GitHubAdapter) mapProject(source GithubProject) Project {
	return Project{
		DisplayName: source.Name,
		SshUrl:      source.SshUrl,
		HttpUrl:     source.CloneUrl,
	}
}
