package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pwuersch/ggt/lib/globals"
)

type GitlabProject struct {
	Name          string `json:"name"`
	SshUrlToRepo  string `json:"ssh_url_to_repo"`
	HttpUrlToRepo string `json:"http_url_to_repo"`
}

type GitlabAdapter struct{}

func (GitlabAdapter) Info() AdapterInfo {
	return AdapterInfo{
		Name:           "gitlab",
		RequiresHost:   true,
		ProvidesSearch: true,
		ProvidesScope:  false,
	}
}

func (gitlab GitlabAdapter) GetProjects(profile globals.Profile, search string, _ string) ([]Project, error) {
	host := profile.Api.Host
	token := profile.Api.Token

	client := &http.Client{}
	url := fmt.Sprintf("https://%s/api/v4/projects?search=%s&simple=false&membership=true&per_page=100&order_by=last_activity_at", host, search)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Project{}, err
	}
	req.Header.Add("PRIVATE-TOKEN", token)

	globals.Debug(fmt.Sprintf("Getting %s", url))
	res, err := client.Do(req)
	if err != nil {
		return []Project{}, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Project{}, err
	}

	gitlabProjects := []GitlabProject{}
	if err = json.Unmarshal(data, &gitlabProjects); err != nil {
		return []Project{}, err
	}

	return globals.Map(gitlabProjects, func(gitlabProject GitlabProject) Project {
		return gitlab.mapProject(gitlabProject)
	}), nil
}

func (GitlabAdapter) CreateProject(profile globals.Profile) error {
	return errors.New("not implemented")
}

func (GitlabAdapter) mapProject(source GitlabProject) Project {
	return Project{
		HttpUrl:     source.HttpUrlToRepo,
		SshUrl:      source.SshUrlToRepo,
		DisplayName: source.Name,
	}
}
