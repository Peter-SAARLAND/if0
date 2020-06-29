
# Development Report (29.06.2020)

## Requirements
1. [`if0 config --set` write back to if0.env](https://gitlab.com/peter.saarland/if0/-/issues/23)
2. [remove updates to environments via `-z` flag](https://gitlab.com/peter.saarland/if0/-/issues/24)
3. `if0 add` with prompts


## Features developed
    
1. `if0 config --set` write back to if0.env 
    
    The env variable that is set, will also be written to the `~/.if0/if0.env` config file.

2.  remove updates to environments via `-z` flag

    This is no longer necessary as we have env specific commands to handle config of zero environments.
    
3. `if0 add` with prompts
    
    To add zero environments with the following prompts:
    > `if0 add`
    >
    > **Name:** 
    >
    > **Use Cloud Provider? [Y/n]**: 
    >
    > **Cloud Provider:** [HCLOUD, digitalocean, aws]
    >
    > **Cloud Provider Authentication:**  
    >
    > **Custom Domain? [y/N]:** 
    >
    > **Enter Base Domain:** 
    
    The content of the prompt inputs will be written to `dash1.env` or `zero.env` as applicable.


## Bugfix/Learning
* Stream container logs

* [Make sure containers are removed](https://gitlab.com/peter.saarland/if0/-/issues/17)