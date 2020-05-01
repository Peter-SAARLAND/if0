
# Development Report (30.04.2020)

## Requirements
The following issues are WIP:
1. [issue5](https://gitlab.com/peter.saarland/if0/-/issues/5) 
2. [issue11](https://gitlab.com/peter.saarland/if0/-/issues/11)

## Features developed
A working code to check if local dependencies such as docker, vagrant have already been installed.
    
## Learning
* Setting up Docker on Windows 10
    *   Pre-requisite: Windows 10 Pro, Enterprise, or Education version
    * Download Docker Desktop for Windows from https://hub.docker.com/editions/community/docker-ce-desktop-windows/
    * To be able to use it, enable Hyper-V in Control Panel  
    Reference: https://success.docker.com/article/manually-enable-docker-for-windows-prerequisites
* Setting up Vagrant on Windows 10
    * Download and install the latest version of VirtualBox. It is important to install the latest version as `vagrant up` doesn't work otherwise.  
    Reference: https://community.oracle.com/thread/3639877
    * If the Windows SmartScreen doesn't allow you to run the Vagrant application, set it to 'Warn' instead of 'Off'
    * Download and install the latest version from https://www.vagrantup.com/downloads.html
    
    Useful references to setup Vagrant:  
    * https://www.swtestacademy.com/quick-start-vagrant-windows-10/
    * https://www.taniarascia.com/what-are-vagrant-and-virtualbox-and-how-do-i-use-them/
    
* Dependency injection in golang

    When writing unit test cases in golang, to mock third party dependencies, dependency injection is the recommended way to go.
    
    For further references:   
    * https://stackoverflow.com/a/19168875/5772695
    * https://medium.com/agrea-technogies/mocking-dependencies-in-go-bb9739fef008


## Bugfixes
* To prompt the user for passphrase only if the SSH key demands it.
Solved as a part of: https://gitlab.com/peter.saarland/if0/-/issues/11

* To handle git pull deleting local changes - WIP.