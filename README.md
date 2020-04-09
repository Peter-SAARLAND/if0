# ToDo

- Implement https://github.com/semantic-release/semantic-release

**Setting up:**

To be able to use `if0` app, `cd` to `if0` directory, and run `go install if0`. This creates an executable in the GOPATH `bin` directory.

Usages of **addConfig** command:

1. `if0 addConfig`

    Prints the current running configuration on the console.
    If there is no current running configuration available, `if0` creates on at `~/.ifo/if0.env`
2. `if0 addConfig --set varName=varValue`
    
    Updates the environment variable  _varName_ with _varValue_
    
    Multiple variables can also be set using comma-separated key-value pairs. For example: `if0 addConfig --set "varName1=varValue1,varName2=varValue2"`
3. `varName=varValue if0 addConfig <args>`

    Updates the variable _varName_ with value _varValue_ before running the `if0` command 
4. `if0 addConfig path/to/configFile.env`

    Takes a backup of the current running configuration file (`if0.env`), and replaces it with the configuration from `configFile.env`