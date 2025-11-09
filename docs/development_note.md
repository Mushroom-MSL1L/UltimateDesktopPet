# Development Notes

## How to build
* Write correct Go packages and imported module. 
    * from outside e.g. `github.com/gin-gonic/gin`
    * from inside e.g. `"UltimateDesktopPet/internal/synchronization"`
* Initialize Go module if not yet done
    * `go mod init UltimateDesktopPet`
* Download dependencies
    * `go mod tidy`
* Build / Run the project
    * `go build main.go`
    * `go run main.go`
* If you want to **stop** anything about internet process,
    * Press `Ctrl + C` to gracefully stop the process.
    * Don't `Ctrl + Z` or `Ctrl + D`, it may cause unexpected behavior. (usually block your port resources)

## How to document API with Swagger
* Install swag CLI tool (There website to download (swaggo)[https://github.com/swaggo/swag.git])
* Write comments (follow the [documentation](https://github.com/swaggo/swag.git))
    * Please ensure the url of swag comments in `main.go` is equal to the one of your system. 
* Generate docs 
    * `swag fmt ; swag init`
* Run your application. (For the first time, you need to initialize modules first.)
* Browse to `http://localhost:8080/swagger/index.html` to see the API documentation UI.


## Workflow 
```
Wails startup
    |
    | Load configuration
    | Initialize system
        | database
        | network
    | Start services
        | go: pet 
            | init 
            | loop   <--------------|
        | go: synchronization       |
            | init                  |
            | loop   <--------------|
        | go: chat box              |
            | init                  |
            | loop   <--------------|
Wails binding                       |
    | Start GUI                     |
        | Bottoms and event         |
            | go: event -------------
```

## image management 
### Pet image management 
* Configuration need to specify what path under `./assets/petImages` need to load as animation. 
* If there's anything that expected to exist but not in user pet image path, use the default ones. 
```
// in configuration 
petImageFolder: "myPetResource" // it means use the path of ./assets/petImages/myPetResource/*

// path structure 
.
└── assets
    ├── db
    │   └── udp.db
    └── petImages
        ├── default
        │   └── left_move 
        |   │   ├── 1.png
        |   │   ├── 2.png
        |   │   ├── 3.png
        |   │   └── 4.png
        │   └── drag 
        |       ├── 1.png
        |       ├── 2.png
        |       ├── 3.png
        |       └── 4.png
        └── myPetResource
            └── left_move 
            │   ├── 1.png
            │   ├── 2.png
            │   ├── 3.png
            │   └── 4.png
            └── drag 
                ├── 1.png
                ├── 2.png
                ├── 3.png
                └── 4.png
```

### Items / Activity image management 
* Configuration need to specify what path under `./assets/food` need to load as animation. 
* If there's anything that expected to exist but not in user pet image path, use the default ones. 
* Attributes and exact path of one food image should be specified in the database. 
```
// in configuration 
itemsDBDir:            myItemDB.db          // it means use the path of ./assets/db/myItemDB.db 
itemsImageFolder:      myItemResource       // it means use the path of ./assets/items/myFoodResource/*
activitiesDBDir:       myActivityDB.db      // it means use the path of ./assets/db/myActivityDB.db
activitiesImageFolder: myActivityResource   // it means use the path of ./assets/activities/myActivityResource/*


// path structure 
//// just like pet path structure
```
