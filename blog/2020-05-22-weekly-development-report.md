
# Development Report (22.05.2020)

## Requirements
1. `if0 environment plan [$NAME]`
2. `if0 environment provision [$NAME]`
3. `if0 environment zero [$NAME]`

## Features developed
    
*$NAME* is optional. If no *$NAME* is given, the present working directory is assumed to be the *Environment*; if $NAME is given, the corresponding *Environment* in ~/.if0/.environments/ is chosen.

1. `if0 environment plan [$NAME]`
    
    This command initializes the necessary Terraform provider modules for the environment `[$NAME]`.

2. `if0 environment provision [$NAME]`

    This command is used to trigger the `make provision` command from `zero` to provision the environment `[$NAME]`.
    
3. `if0 environment zero [$NAME]`

    This command is used to trigger the `make zero` command from `dash1`, which creates the zero infrastructure for the environment `[$NAME]`.

    
    
## Learning/Others
* Using [Docker API](https://docs.docker.com/engine/api/sdk/) to run docker commands via golang.

