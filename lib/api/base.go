package api

import (
	"fmt"

	"github.com/pwuersch/ggt/lib/globals"
)

type ApiAdapter interface {
	Info() AdapterInfo
	GetProjects(globals.Profile, string, string) ([]Project, error)
	CreateProject(globals.Profile) error
}

type AdapterInfo struct {
	Name           string
	RequiresHost   bool
	ProvidesSearch bool
	ProvidesScope  bool
}

type Project struct {
	HttpUrl     string
	SshUrl      string
	DisplayName string
}

var AvailableAdapters = []ApiAdapter{
	GitlabAdapter{},
	GitHubAdapter{},
	BitbucketAdapter{},
}

func GetAdapterByName(name string) (ApiAdapter, error) {
	for _, adapter := range AvailableAdapters {
		if adapter.Info().Name == name {
			return adapter, nil
		}
	}

	return GitlabAdapter{}, fmt.Errorf("no provider with name %s found", name)
}

func GetAdapterNames() []string {
	return globals.Map(AvailableAdapters, func(adapter ApiAdapter) string {
		return adapter.Info().Name
	})
}
