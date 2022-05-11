package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/pwuersch/ggt/lib/globals"
)

type BitbucketProject struct {
	Name     string `json:"name"`
	SshUrl   string `json:"ssh_url"`
	CloneUrl string `json:"clone_url"`
}

type BitbucketAdapter struct{}

func (BitbucketAdapter) Info() AdapterInfo {
	return AdapterInfo{
		Name:           "bitbucket",
		RequiresHost:   true,
		ProvidesSearch: true,
		ProvidesScope:  true,
	}
}

func (bitbucket BitbucketAdapter) GetProjects(profile globals.Profile, search string, scope string) ([]Project, error) {
	token := profile.Api.Token

	client := &http.Client{}
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s?q=name~\"%s\"", scope, search)
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

	githubProjects := []BitbucketProject{}
	if err = json.Unmarshal(data, &githubProjects); err != nil {
		return []Project{}, err
	}

	return globals.Map(githubProjects, func(bitbucketProjects BitbucketProject) Project {
		return bitbucket.mapProject(bitbucketProjects)
	}), nil
}

func (BitbucketAdapter) CreateProject(profile globals.Profile) error {
	return errors.New("not implemented")
}

func (BitbucketAdapter) mapProject(source BitbucketProject) Project {
	return Project{
		DisplayName: source.Name,
		SshUrl:      source.SshUrl,
		HttpUrl:     source.CloneUrl,
	}
}
