[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release) [![pipeline status](https://gitlab.com/peter.saarland/if0/badges/master/pipeline.svg)](https://gitlab.com/peter.saarland/if0/-/commits/master)

Maintained by [Peter.SAARLAND | DevOps Consultants](https://www.peter.saarland) - Helping companies to develop software with startup speed and enterprise quality.

Additional Links:

- [ns0](https://gitlab.com/peter.saarland/ns0/) - The container-native DNS Proxy
- [zero](https://gitlab.com/peter.saarland/zero/) - The Application-Platform
- [dash1](https://gitlab.com/peter.saarland/dash1/) - Virtual Infrastructure for Zero

### **Setting up:**

To be able to use `if0` app, `cd` to `if0` directory, and run `go install if0`. This creates an executable in the GOPATH `bin` directory.

### Usages of **config** command:

1. `if0 config`

    * Prints the current running configuration on the console.
    * If there is no current running configuration available, `if0` creates one at `~/.ifo/if0.env`
2. `if0 config --set varName=varValue`
    
    * Updates the environment variable  _varName_ with _varValue_
    
    * Multiple variables can also be set using comma-separated key-value pairs. For example: `if0 config --set "varName1=varValue1,varName2=varValue2"`
3. `varName=varValue if0 config <args>`

    * Updates the variable _varName_ with value _varValue_ before running the `if0` command 
4. `if0 config --add=path/to/configFile.env` or `if0 config -z --add=path/to/configFile.env`

    * Takes a backup of the current running configuration file (`if0.env`), and replaces it with the configuration from `configFile.env`
    
    * `-z` or `--zero` adds or updates the config file in `~/.ifo/.environments` directory
5. `if0 config --merge --src=SRC_CONFIG.env [--dst=DST_CONFIG.env]`
    
    * `--merge` or `-m`

    * Takes a backup of `if0.env` and then merges it with configuration from `SRC_CONFIG.env` 
    
    * `--dst` flag is optional. By default, `if0.env` file is chosen.
    
    * This command can also be used to merge zero-cluster configuration files with the use of `--zero` or `-z`flag. 
    
    * `--dst` flag is again optional; in this case, `dst` file name is assumed to be the same as the `src` config file name. This requires the user to know the destination file name.
    
6. `if0 config --sync` (temporarily disabled, see `if0 sync`)
    
    * This command synchronizes configuration files with the git repository mentioned in the `if0.env` file under variable `REMOTE_STORAGE`.
    
    * If the user uses an SSH link as the `REMOTE_STORAGE`, then a passphrase protected `id_rsa` SSH key is required to be present at `~/.ssh` for authentication.
    
    * If the user uses an HTTPS link, they will be prompted to enter `username` and `password` during sync operation.
    
    * Additionally, the user can also choose to add/commit/push the local changes by entering 'y' when prompted, or 'n' if they do not want the local changes to be pushed to the repository.

### Environment commands:

1. `if0 add test-env [git@gitlab.com:test-env.git]`
    
    This command is used to add zero environments.

    There are three ways to add an environment:
    
    1. Using `GL_TOKEN`
    
        This requires the user to set the `GL_TOKEN` variable in the `~/.if0/if0.env` configuration file. The value for this variable is the GitLab personal access token.
    
        After setting the variable, running `if0 env add env-1` will create a private project titled `env-1` on Gitlab.
        
        The same environment is created locally with initial requirements (`zero.env`, `.gitlab-ci.yml`, `.ssh` directory with `id_rsa` and `id_rsa.pub` files) and synced with the private project `env-1`.
        
    2.  By running the command `if0 env add env-2 git@gitlab.com:peter.saarland/env-2.git`
    
        This command would clone the repository at `~/.if0/.environment/gitlab.com/peter.saarland/env-2` with the initial requirements, and sync these changes with the remote repository.
        
    3. Running the command `if0 env add env-3` with an empty `GL_TOKEN` or no remote repository url
        
        In this case, the environment is created locally at `~/.if0/.environments/env-3`
    
2. `if0 sync [env-name]`
    
    This command is used to synchronize a zero environment with its remote repository. 
    
    `env-name` is optional; when no `env-name` is provided, the current working directory is assumed to be the zero environment to be synced.

3. `if0 plan [env-name]`

    This command corresponds to `dash1 make plan`. It initializes the necessary Terraform provider modules for the Environment `env-name` and then creates a plan in ~/.if0/.environments/$NAME/dash1.plan`
    
    
4. `if0 infrastructure [env-name]` 
    
    This command corresponds to `dash1 make infrastructure`. It generates configuration necessary for zero.

5. `if0 platform [env-name]` 

    This command corresponds to `zero make platform`

6. `if0 destroy [env-name]`

    This command corresponds to `dash1 make destroy`

### Other commands:

1. `if0 status dep`

    This command is used to check for the status of dependencies such as docker and vagrant.

2. `if0 version`

    This command prints the `[if0 version](https://gitlab.com/peter.saarland/if0/-/blob/master/README.md#if0-version)`.
    
### **Developer Documentation**

1. ##### Making use of SSH Keys to login to a server via PUTTY  

    Enter your HOST Name or IP address, choose SSH as the 'Connection Type'
    
    ![](docs/images/ssh1.png)
    
    
    Choose 'Auth' under 'SSH' and provide private key 
    
    ![](docs/images/privkey.png)
    
    
2. ##### Installing go 
    
    * Download the latest binary from https://golang.org/dl/ and run it.
    * Ensure that the GOPATH (go workspace) is added to the PATH environment variable. This can be done automatically by choosing the option in the installer (in windows) or manually.  
    * The GOPATH should typically contain the following folders: `bin, pkg, src`
    
    OS-specific installation instructions can be found here: https://golang.org/doc/install
  
    **For windows**
    
    https://www.freecodecamp.org/news/setting-up-go-programming-language-on-windows-f02c8c14e2f/
    
    **For ubuntu**
    
    https://medium.com/better-programming/install-go-1-11-on-ubuntu-18-04-16-04-lts-8c098c503c5f
    
    **For MAC OS X**
    
    https://medium.com/golang-learn/quick-go-setup-guide-on-mac-os-x-956b327222b8

3. ##### Installing and Setting up Cobra

    * Run the following command to install cobra: `go get github.com/spf13/cobra/cobra`
    
    * Creating initial application code
    
        1. To create the initial application code, use the following command: `cobra init --pkg-name path/to/appName`
    ![](docs/images/cobrainit.png)
    
        2. To add new commands to your app, use the following command: `cobra add cmdName`. This will add the bare-bones code necessary for the command. 
    ![](docs/images/cmd.png)
    
        3. After including code for the necessary functionality, run the following commands to create a binary for the app. This binary is saved in `GOPATH/bin`.
        ![](docs/images/cmdrun.png)
        
    Reference: https://github.com/spf13/cobra/blob/master/cobra/README.md

4. #### `if0 version`
    Run the following commands to build and install `if0` app with commit SHA as the version number.

>     root@zero-gayathri-dev:~/if0# if0 version
>     if0 version:
>     root@zero-gayathri-dev:~/if0# export IF0VERSION=$(git rev-list -1 HEAD)
>     root@zero-gayathri-dev:~/if0# go build -ldflags "-X main.Version=$IF0VERSION"
>     root@zero-gayathri-dev:~/if0# go install -ldflags "-X main.Version=$IF0VERSION" if0
>     root@zero-gayathri-dev:~/if0# if0 version
>     if0 version: af6b3929a641e2ea9293f4c5956bbe9fc61f000f



    