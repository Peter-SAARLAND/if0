package gitlabclient

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"if0/common"
	"if0/config"
)

func CreateProject(projectName, gitlabToken string) (string, error) {
	clientOptions := gitlab.WithBaseURL(getIf0RegistryUrl())
	git, err := gitlab.NewClient(gitlabToken, clientOptions)
	if err != nil {
		fmt.Println("Error: Creating gitlab client -", err)
		return "", err
	}
	projectOptions := &gitlab.CreateProjectOptions{
		Name:                  gitlab.String(projectName),
		Visibility:            gitlab.Visibility(gitlab.PrivateVisibility),
		MergeMethod:           gitlab.MergeMethod(gitlab.NoFastForwardMerge),
	}
	project, _, err := git.Projects.CreateProject(projectOptions)
	if err != nil {
		fmt.Println("Error: Creating GitLab project -", err)
		return "", err
	}
	fmt.Println("Project created successfully at ", project.HTTPURLToRepo)
	return project.SSHURLToRepo, nil
}

func getIf0RegistryUrl() string {
	config.ReadConfigFile(common.If0Default)
	return config.GetEnvVariable("IF0_REGISTRY_URL")
}