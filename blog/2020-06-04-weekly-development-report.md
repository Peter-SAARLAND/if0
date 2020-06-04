
# Development Report (04.06.2020)

## Requirements
1. [issue16](https://gitlab.com/peter.saarland/if0/-/issues/16) 

## Features developed
    
1. `if0 env add $ENV_NAME [$REPO_URL]`
    
    The functionality for `if0 env add` previously required the user to provide a repository URL to add the environment locally. This requirement has been changed in the current implementation.
    
    There are three possibilities to add an environment:
    
    1. Using `GL_TOKEN`
    
        This requires the user to set the `GL_TOKEN` variable in the `~/.if0/if0.env` configuration file. The value for this variable is the GitLab personal access token.
    
        After setting the variable, running `if0 env add env-1` will create a private project titled `env-1` on Gitlab.
        
        The same environment is created locally with initial requirements (`zero.env`, `.gitlab-ci.yml`, `.ssh` directory with `id_rsa` and `id_rsa.pub` files) and synced with the private project `env-1`.
        
    2.  By running the command `if0 env add env-2 git@gitlab.com:peter.saarland/env-2.git`
    
        This command would clone the repository at `~/.if0/.environment/gitlab.com/peter.saarland/env-2` with the initial requirements, and sync these changes with the remote repository.
        
    3. Running the command `if0 env add env-3` with no/an empty `GL_TOKEN` and no remote repository url
        
        In this case, the environment is created locally at `~/.if0/.environments/env-3`

2. Creating a `defaultIf0.env` in the repository to make it easier to add default env variables to be included in `~/.if0/if0.env` configuration file. 
        


## Bugfix/Learning
* Using GitLab API to create a private project.

    References: 
    * https://docs.gitlab.com/ee/api/projects.html#create-project
    * https://github.com/xanzy/go-gitlab