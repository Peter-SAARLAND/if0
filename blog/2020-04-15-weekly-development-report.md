
# Development Report (15.04.2020)

## Requirements
This week was all about adding more features to [the CLI prototype](https://gitlab.com/peter.saarland/if0/-/issues/2)

## Features developed

1. Rewriting `merge` logic. 

    * `merge` requires the user to mention `src` and/or `dst` configuration files to carry out the merge operation.
    * For `if0.env`, it is sufficient to include the `src` flag; `dst` flag is automatically set.
    * For `zero-cluster.env` files, it is sufficient to include only the `src` flag, but the filename of the src file should match with the  configuration file in the `~/.if0/.environments` directory with which the user intends to merge the configuration.
    * The user can also specify src and dst configuration files separately using `src` and `dst` flags respectively. Please ensure that the dst file mentioned in the flag is already present in the `~/.if0/.environments` directory.
    
2. Automatic Garbage Collection

    * This feature deletes older configuration files that are backed-up in the `~/.if0/.snapshots` directory. Configuration files that are older than a certain number of days (`GC_PERIOD`) are automatically deleted.
    * Requires the user to set the `GC_AUTO` environment variable to activate automatic garbage collection in the running configuration file `if0.env`
    * The user can also specify the garbage collection period by setting the `GC_PERIOD` variable. By default, the value is set to 30 days.

3. Renaming the command `addConfig` to `config`

4. Including a new flag `--add` to add or update configuration files.

## Learning

* Setting up a CI job to display test coverage results in the pipeline.

## Bugfixes
None

