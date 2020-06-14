
# Development Report (12.06.2020)

## Requirements
1. Remove `env` command
2. Change commands under `env` to be standalone commands
3. Htpasswd hash in the `zero.env` file when an environment is added
4. Create GitLab projects based on `IF0_REGISTRY_URL` ([issue20](https://gitlab.com/peter.saarland/if0/-/issues/20))

## Features developed
    
1. `if0 env commands` to standalone commands
    
    Commands which were previously under `env` are now standalone commands:
    * `if0 add test-env [git@gitlab.com:test-env.git]`
    * `if0 sync [env-name]`
    * `if0 plan [env-name]`
    * `if0 infrastructure [env-name]` (was previously `if0 env zero`)
    * `if0 platform [env-name]` (was previously `if0 provision`)
    * `if0 destroy [env-name]`

2.  Htpasswd hash in the `zero.env` file

    Creation of htpasswd hash has been implemented two ways:
    
    * Using htpasswd command on OSs where the command is available.
    * Using docker run on OSs (eg: windows) where htpasswd command is not available.
    
3. Using `IF0_REGISTRY_URL` in the creation of GitLab Projects. 
    
    By default, the value of `IF0_REGISTRY_URL` is `gitlab.com`. This value can be configured to any other self-hosted instance of GitLab.

## Bugfix/Learning
* `if0 add`: To not include env init files when the remote repository already contains env init files.
* Handling ssh failure for self-hosted GitLab registries. 
    
    `ssh: handshake failed: knownhosts: key is unknown`
    This was fixed by using `InsecureIgnoreHostKey` in SSH auth.
* Remove port number from the env directory path