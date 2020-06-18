
# Development Report (18.06.2020)

## Requirements
1.  [if0 list](https://gitlab.com/peter.saarland/if0/-/issues/21)
2.  [if0 inspect](https://gitlab.com/peter.saarland/if0/-/issues/21)
3. [Default logo.png in zero environments](https://gitlab.com/peter.saarland/if0/-/issues/22)
4. Creating gitlab projects for groups

## Features developed
    
1. `if0 list` 
    
    This command lists all the zero environments available at `~/.if0/.environments`

2.  `if0 inspect [env-name]`

    This command displays the configuration available in all the *.env files of the environment `env-name`. If `env-name` is not provided, the current working directory is assumed to be the zero environment to be inspected.
    
3. Including default [logo.png](https://gitlab.com/peter.saarland/scratch/-/blob/master/logo.png) when zero environments are added/created.
    
4. Creating Gitlab projects for groups.

    A new environment variable, `IF0_REGISTRY_GROUP` has been included in `if0.env` to help with creation of Gitlab projects under specific groups. The value of the `IF0_REGISTRY_GROUP` variable is the name of the Gitlab group.

## Bugfix/Learning
* Dockerizing if0 application (WIP)