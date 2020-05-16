
# Development Report (16.05.2020)

## Requirements
1. [Provide function to check for local Dependencies](https://gitlab.com/peter.saarland/if0/-/issues/5) 
2. [Add 'if0 environment load'](https://gitlab.com/peter.saarland/if0/-/issues/12)
3. [Add version command](https://gitlab.com/peter.saarland/if0/-/issues/10)

## Features developed
1. `if0 environment load [$NAME]`

    $NAME is optional. If no $NAME is given, the present working directory is assumed to be the *Environment* to be loaded; if $NAME is given, the corresponding *Environment* in ~/.if0/.environments/ is loaded.
    
    This command loads the environment variables present at the environment **[$NAME]** within `~/.if0/.environments` directory. By load, we mean that it fetches all the .env files from the *Environment*, and exports them to the environment of the shell in which `if0` is running.

2. `if0 status dep`

    This command checks if required dependencies such as Docker and Vagrant are installed.
    
## Learning/Others
* To set `if0` version during `go build`

    This was done as part of the following feature: [Add version command](https://gitlab.com/peter.saarland/if0/-/issues/10)
    
    Steps to build and install `if0` app with commit SHA as the version number:

>     root@zero-gayathri-dev:~/if0# if0 version
>     if0 version:
>     root@zero-gayathri-dev:~/if0# export IF0VERSION=$(git rev-list -1 HEAD)
>     root@zero-gayathri-dev:~/if0# go build -ldflags "-X main.Version=$IF0VERSION"
>     root@zero-gayathri-dev:~/if0# go install -ldflags "-X main.Version=$IF0VERSION" if0
>     root@zero-gayathri-dev:~/if0# if0 version
>     if0 version: af6b3929a641e2ea9293f4c5956bbe9fc61f000f


## Bugfixes
* Cloning empty repository when `if0 environment add repo_name` is called with a repository that is empty.