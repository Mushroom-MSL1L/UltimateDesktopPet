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
* Users can locally load their own image resources as animations by configuring the paths under `./assets/petImages`.
* However, if any required animation is missing, the system will automatically fall back to the default animations.
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
* Users can locally load their own image resources as animations by configuring the paths under `./assets/activities` and `./assets/items`, along with the corresponding attributes defined in `./assets/db/images.db`.
* However, if any required animation is missing, the system will automatically fall back to the default animations.
```
// in configuration 
imageDBDir:            ./assets/db/images.db                     // it means use the db ./assets/db/images.db
itemsImageFolder:      ./assets/items/myItemResource             // it means use the path of ./assets/items/myItemResource/*
activitiesImageFolder: ./assets/activities/myActivityResource    // it means use the path of ./assets/activities/myActivityResource/*


// path structure 
//// just like pet path structure
```
