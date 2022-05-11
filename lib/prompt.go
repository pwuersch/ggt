package lib

import (
	"bufio"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pwuersch/ggt/lib/api"
	"github.com/pwuersch/ggt/lib/globals"
)

const (
	GitUrlRegex = "((git|ssh|http(s)?)|(git@[\\w\\.]+))(:(\\/\\/)?)([\\w\\.@\\:/\\-~]+)(\\.git)?(\\/)?"
)

func ReadInput(prompt string, reader *bufio.Reader) string {
	fmt.Printf("%s: ", prompt)
	input, _ := reader.ReadString('\n')

	return strings.Replace(input, "\n", "", -1)
}

func PromptInput(label string, validate promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	return prompt.Run()
}

func PromptYesNo(prompt string) (bool, error) {
	_, answer, err := SelectItem(prompt, []string{"Yes", "No"})

	return answer == "Yes", err
}

func SelectItem(label string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	return prompt.Run()
}

func PromptSelectProfile(prompt string) (globals.Profile, int, error) {
	profiles, err := ReadProfiles()
	if err != nil {
		return globals.Profile{}, 0, err
	}

	if len(profiles) == 0 {
		return globals.Profile{}, 0, errors.New("no profile available, did you already create one?")
	}

	profileNames := globals.Map(profiles, func(profile globals.Profile) string {
		return profile.Name
	})

	index, _, err := SelectItem(prompt, profileNames)
	if err != nil {
		return globals.Profile{}, 0, err
	}

	return profiles[index], index, nil
}

func PromptNewProfile() (globals.Profile, error) {
	name, err := PromptInput("Profile Name", ValidateEmpty)
	if err != nil {
		return globals.Profile{}, err
	}

	email, err := PromptInput("Git Email", ValidateEmail)
	if err != nil {
		return globals.Profile{}, err
	}

	commitName, err := PromptInput("Git Commit Name", ValidateEmpty)
	if err != nil {
		return globals.Profile{}, err
	}

	profile := globals.Profile{
		Name:       name,
		Email:      email,
		CommitName: commitName,
	}

	createApiProfile, err := PromptYesNo("Do you want to add api credentials for additional features?")
	if err != nil {
		return profile, err
	}

	if createApiProfile {
		apiProfile, err := PromptNewApiProfile()
		if err != nil {
			return profile, err
		}
		profile.Api = apiProfile
	}

	return profile, nil
}

func PromptNewApiProfile() (globals.ApiProfile, error) {
	apiProfile := globals.ApiProfile{}

	_, adapterName, err := SelectItem("Select the adapter you want to use", api.GetAdapterNames())
	if err != nil {
		return apiProfile, err
	}
	apiProfile.Adapter = adapterName

	adapter, err := api.GetAdapterByName(adapterName)
	if err != nil {
		return apiProfile, err
	}

	if adapter.Info().RequiresHost {
		host, err := PromptInput("Instance Host (without proto)", ValidateEmpty)
		if err != nil {
			return apiProfile, err
		}
		apiProfile.Host = host
	}

	token, err := PromptInput("Access Token", ValidateEmpty)
	if err != nil {
		return apiProfile, err
	}
	apiProfile.Token = token

	return apiProfile, nil
}

func ValidateEmpty(input string) error {
	if strings.TrimSpace(input) == "" {
		return errors.New("input must not be blank")
	}
	return nil
}

func ValidateEmail(input string) error {
	if _, err := mail.ParseAddress(input); err != nil {
		return errors.New("input must be a valid email")
	}
	return nil
}

func ValidateGitUrl(input string) error {
	if matched, err := regexp.MatchString(GitUrlRegex, input); !matched || err != nil {
		return errors.New("input must be a valid git url [ssh or http(s)]")
	}
	return nil
}
