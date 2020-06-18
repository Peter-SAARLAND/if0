package gitlabclient

import (
	"errors"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"if0/common"
	"if0/config"
)

func CreateProject(projectName, gitlabToken string) (string, string, error) {
	clientOptions := gitlab.WithBaseURL(getIf0RegistryUrl())
	git, err := gitlab.NewClient(gitlabToken, clientOptions)
	if err != nil {
		fmt.Println("Error: Creating gitlab client -", err)
		return "", "", err
	}

	// if group ID is 0 (group not found/invalid)
	groupId, err := getIf0GroupId(git)
	if groupId == 0 || err != nil {
		return "", "", err
	}

	projectOptions := &gitlab.CreateProjectOptions{
		Name:        gitlab.String(projectName),
		Visibility:  gitlab.Visibility(gitlab.PrivateVisibility),
		MergeMethod: gitlab.MergeMethod(gitlab.NoFastForwardMerge),
		NamespaceID: gitlab.Int(groupId),
	}

	project, _, err := git.Projects.CreateProject(projectOptions)
	if err != nil {
		fmt.Println("Error: Creating GitLab project -", err)
		return "", "", err
	}
	fmt.Println("Project created successfully at ", project.HTTPURLToRepo)
	return project.SSHURLToRepo, project.HTTPURLToRepo, nil
}

func getIf0RegistryUrl() string {
	config.ReadConfigFile(common.If0Default)
	return config.GetEnvVariable("IF0_REGISTRY_URL")
}

func getIf0GroupId(client *gitlab.Client) (int, error) {
	config.ReadConfigFile(common.If0Default)
	groupName := config.GetEnvVariable("IF0_REGISTRY_GROUP")

	var namespaceId int
	namespace, _, err := client.Namespaces.SearchNamespace(groupName)
	if err != nil {
		return 0, err
	}
	for _, n := range namespace {
		if n.Name == groupName {
			namespaceId = n.ID
		}
	}
	if namespaceId == 0 {
		return 0, errors.New("group not found/invalid group")
	}
	return namespaceId, nil
}
