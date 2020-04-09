
# Development Report (09.04.2020)

## Requirements
The features described in the **Features developed** section follow the requirements detailed in the following git issues/features:

1. [issue1](https://gitlab.com/peter.saarland/if0/-/issues/1)
2. [issue2](https://gitlab.com/peter.saarland/if0/-/issues/2)

## Features developed

1. Printing current running configuration for `if0`. If no configuration is available, create a default if0.env file at `~/.if0` configuration with key "IF0_VERSION" and appropriate value.
2. Setting environment variables using `--set` flag.
Examples: 

    `if0 addConfig --set var=value`
    
    `if0 addConfig --set "var1=value1,var2=value2`
3. Setting environment variables directly before running the `if0` command
Example: `var=value if0 addConfig`

4. Replacing the current running configuration.
Example: `if0 addConfig path/to/config.env`

5. Merging new configuration with the current running configuration.
Example: `if0 addConfig --merge path/to/config.env`

## Learning
* We have used [cobra](https://github.com/spf13/cobra) to design the if0 CLI tool. Setting this up is pretty simple, as explained in their README document.
* Adding new commands to if0 is simple. 
Example: `cobra add sampleCmd` adds the bare-bones code necessary to work with the command `sampleCmd`
* Setting environment variables directly from cmd terminal. 
**Expectation**: set an environment variable to a desired value before executing commands in if0 tool. 
    * `testvar=testval; if0 addConfig`
    * `testvar=testval if0 addConfig`

    Note the missing semicolon in the second command. This makes a difference to how these commands are executed.
    
    Why?  [Reference](https://unix.stackexchange.com/questions/36745/when-to-use-a-semi-colon-between-environment-variables-and-a-command/36829#36829?newreg=b41d7ccacbb843d0b9fa11556b515668)

    `VAR=value; somecommand` is equivalent to
    
     `VAR=value`  
     `somecommand`  
     These are unrelated commands executed one after the other. 
     The first command assigns a value to the shell variable VAR. 
     Unless VAR is already an environment variable, it is not exported to the environment, it remains internal to the shell. 
     A statement export VAR would export VAR to the environment.
    
     `VAR=value somecommand` is a different syntax. 
     The assignment `VAR=value` is to the environment, but this assignment is only made in the execution environment of `somecommand`, not for the subsequent execution of the shell. 


## Bugfixes
* To include a default if0.env file when the command is run with no config file initially
* Check for already existing dir during dir creation
* Include go.mod and go.sum files to be able to run the application outside of gopath
* Include O_RDWR mode to if0.env file (just create works on windows, but not on mac os)