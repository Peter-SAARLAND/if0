
# Development Report (29.05.2020)

## Requirements
1. Check for certain necessary .env files after env add
2. Change command `environment` to `env`


## Features developed
    
1. Check if the following .env files are present in the environment upon running the command `if0 env add envName`
    * `zero.env` - if not, add an empty `zero.env` file
    * `.ssh` directory with `id_rsa` and `id_rsa.pub` files
    * `.gitlab-ci.yml` with contents:
        >   include:   
             - remote: 'SHIPMATE_WORKFLOW_URL'
        
        where SHIPMATE_WORKFLOW_URL variable is defined in the `if0.env` file.
        
        If `if0.env` does not contain SHIPMATE_WORKFLOW_URL, it is added to the `if0.env` file with "https://gitlab.com/peter.saarland/shipmate/-/raw/master/shipmate.gitlab-ci.yml" as the default value.


## Bugfix/Learning
* Permission denied for `.ssh` folder in MacOS. 

    Solution - Set 0700 for `.ssh`, 0600 for `id_rsa`, and 0644 for `id_rsa.pub` (yet to be tested)
