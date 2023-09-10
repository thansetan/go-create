# Go-Create

A helper to create a Go project directory (I'm too lazy to create it one by one every time I start a new project, lol).

## How to Use:

1. Install the binary using:

```
go install github.com/thansetan/go-create@latest
```

2. To create a new project directory, you can use either of the following commands:
    - `go-create <project-name>`
    - `go-create -n <project-name>`

## The `go-create` command will create the following directories and files:


    <PROJECT-DIR-NAME>
    ├── cmd/
    │   └── app/
    │       └── main.go
    ├── internal/
    │   ├── app/
    │   │   ├── delivery/
    │   │   │   └── ...
    │   │   ├── repository/
    │   │   │   └── ...
    │   │   └── usecase/
    │   │       └── ...
    │   ├── domain/
    │   │   └── ...
    │   └── infrastructure/
    │       ├── database/
    │       │   └── ...
    │       └── http/
    │           └── ...
    ├── pkg/
    │   └── ...
    ├── .env
    ├── Dockerfile
    └── go.mod


## Notes:

1. If no project name is specified, it will create a project named "my-project".
2. If the `<project-name>` is specified as `github.com/username/project-name` (e.g., `go-create github.com/username/project-name`), the project directory name will be `project-name`, but the module name will still be `github.com/username/project-name`.