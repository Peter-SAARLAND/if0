
# Development Report (07.05.2020)

## Requirements
Remote synchronization of zero-config files in `~/.if0/.environments` directory


## Features developed
1. `if0 environment add *repo_url*`

    This command clones the git repository present at *repo_url* within `~/.if0/.environments` directory to help synchronize the zero-cluster configuration files.
2. `if0 environment sync *repo_name*`

    This command synchronizes the contents of the local and remote copies of the git repository *repo_name*
    
## Learning/Others
* Adding unit test cases, improving code coverage.
* Disabling `if0 config --sync` feature so as to avoid discrepancies in synchronization of contents within `~/.if0/.environments`

## Bugfixes
None