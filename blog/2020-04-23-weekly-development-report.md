
# Development Report (23.04.2020)

## Requirements
During this week, 
* We developed features for requirements described [here](https://gitlab.com/peter.saarland/if0/-/issues/3)

* And analyzed godoc for [issue 8](https://gitlab.com/peter.saarland/if0/-/issues/8) 

## Features developed

**Including `sync` logic.** 

   * Adding a new command `sync` to `if0` to manage synchronization of configuration files with remote git repositories.
   * `if0 sync` does the following:
       * `git init` to create a git repository on the local machine, if one is not already present at `~/.if0`
       * `git remote add $REMOTE_STORAGE` to add a remote origin for the repository mentioned in `REMOTE_STORAGE`. 
            * Include REMOTE_STORAGE variable in the `if0.env` file with either an HTTPS or an SSH git repository link as the value. 
            * For SSH link, it is required to also include `SSH_KEY_PATH` variable in the `if0.env` config file, and its value should be filepath of the private ssh key.
       * `git pull $REMOTE_STORAGE` to pull in changes from the repository.
       * `git status` to check for local changes
       * `git add`, `git commit`, `git push` commands to push local changes to the remote repository. This is optional. The user is prompted to enter a 'y' to push changes, or an 'n' to abort.
    
## Learning
1. **Setting up `if0` on Ubuntu** 
    * Requires installation of Go, and Cobra (refer [README.md](https://gitlab.com/peter.saarland/if0/-/blob/master/README.md#installing-go))
2. **Setting up ssh-keys in Windows, and Ubuntu to sync repositories in Gitlab**
	* On Ubuntu, generate the keys using the `ssh-keygen` command. The newly created keys can be found at `~/root/.ssh`
	
	* On windows - we were getting "no ssh key error" 
	
    	This is because the ssh keys were not in OpenSSH format. When SSH keys are generated using PuttyGen, it does not generate it in OpenSSH format by default. This can be resolved by either of the options below:
    	* We have to explicitly export it in OpenSSH format.(Reference: https://stackoverflow.com/a/5514768/5772695)
    	 
	    * Create ssh keys on the cmd terminal or git bash using the `ssh-keygen` command
3. **Analyzing godoc for automated documentation**
    * **Pros**: godoc requires us to properly comment all of our functions, variables etc. and generates documentation based on our comments. This doesn't require us to do anything extra, so, it is definitely a pro.
    * **Cons**: godoc generates documentation for all the packages that are being used. Not just for the packages that we have implemented (which is desirable). Documentation for other packages isn't really needed as they are anyway available otherwise too. Also, makes it tedious to search for our package.

## Bugfixes
None

