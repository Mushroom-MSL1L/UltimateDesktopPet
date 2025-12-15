# UltimateDesktopPet


* There're 2 subproject in to repo
    * desktop_pet : for user side desktop pet application 
    * sync_server : for developer maintain the synchronization server 
        * Not finished yet. 


## How to run desktop_pet
### For Windows user
* Install Golang 
    * We recommend you to install version 1.24
* Open your powershell 
* Execute `run_me.bat`
    * It's a script file under the root of this repository. 
    * There's a usage quide in it. 
* For the first run, you need to update all resources. 
    * Enter `.\run_me.bat update`
* Directly run 
    * Enter `.\run_me.bat dev`
* Create a shortcut
    * Enter `.\run_me.bat build`
    * It will generate file `UltimateDesktopPet.bat` under the root of this repository. 
    * `UltimateDesktopPet.bat` can be removed to elsewhere. 
    * In GUI, double click `UltimateDesktopPet.bat` to run. 


### For Mac / Linux user 
* Install Golang 
    * We recommend you to install version 1.24
* Open your terminal 
* Execute `run_me.sh`
    * It's a script file under the root of this repository. 
    * There's a usage quide in it. 
* For the first run, you need to update all resources. 
    * Enter `./run_me.sh update`
* Directly run 
    * Enter `./run_me.sh dev`
* Create a shortcut
    * Enter `./run_me.sh build`
    * It will generate file `UltimateDesktopPet.command` for mac, `UltimateDesktopPet.desktop` for linux under the root of this repository. 
    * `UltimateDesktopPet.command` and `UltimateDesktopPet.desktop` can be removed to elsewhere. 
        * MacOS FS has TCC, so not all path can execute `UltimateDesktopPet.command`. Still, you can run `./run_me.sh dev`. 
    * In GUI, double click `UltimateDesktopPet.command` or `UltimateDesktopPet.desktop` to run. 


## Modules 
* UI : [wails](https://github.com/wailsapp/wails)
* Framework : [gin](https://github.com/gin-gonic/gin)
* Database : [gorm](https://github.com/go-gorm/gorm)
* File layout : [go standard project layout](https://github.com/golang-standards/project-layout)
* API documentation : [swaggo/swag](https://github.com/swaggo/swag.git)