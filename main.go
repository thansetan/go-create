package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	projectName           = flag.String("n", "my-project", "name of the project")
	perm                  = os.FileMode(0755)
	mainDirList           = []string{"cmd", "internal", "pkg"}
	internalDirList       = []string{"app", "domain", "infrastructure"}
	internalAppDirList    = []string{"delivery", "repository", "usecase"}
	infrastructureDirList = []string{"database", "http"}

	mainDotGoContent = `package main

func main(){

}`
	dockerFileContent = `FROM golang:GO_VERSION AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main-app ./cmd/app/

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main-app ./main-app

EXPOSE 8080
CMD "./main-app"`
)

func main() {
	var (
		err     error
		dirName string
		file    *os.File
	)

	flag.Parse()

	if flag.NArg() > 0 {
		*projectName = flag.Arg(0)
	}

	splittedProjectName := strings.Split(*projectName, "/")
	dirName = splittedProjectName[len(splittedProjectName)-1]
	cwd, _ := os.Getwd() // in case there's error, this will be used to delete the project dir

	err = createDir(dirName)
	if err != nil {
		if strings.Contains(err.Error(), os.ErrExist.Error()) {
			fmt.Printf("Directory %s already exist!\n", dirName)
			os.Exit(1)
		}
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("An error has occured: %s\nAborting and removing directory: %q for project: %q\n", r, dirName, *projectName)

			if file != nil {
				file.Close()
			}

			os.Chdir(cwd)
			err := os.RemoveAll(path.Join(cwd, dirName))
			if err != nil {
				fmt.Printf("Error when trying to remove directory: %s\n", err.Error())
			}
		}
	}()

	err = os.Chdir(dirName)
	if err != nil {
		panic(err)
	}

	err = initGoProject(*projectName)
	if err != nil {
		panic(err)
	}

	goVersion, err := getGoVersion()
	if err != nil {
		panic(err)
	}

	createDirFromSlice(mainDirList)

	file, err = createFile(".env")
	if err != nil {
		panic(err)
	}

	os.Chdir("cmd")

	err = createDir("app")
	if err != nil {
		fmt.Println("2")
		panic(err)
	}

	os.Chdir("app")

	file.Close() // .env needs to be closed manually
	file, err = createFile("main.go")
	if err != nil {
		panic(err)
	}

	err = writeToFile(file, mainDotGoContent)
	if err != nil {
		panic(err)
	}

	os.Chdir("../../internal")

	createDirFromSlice(internalDirList)

	os.Chdir("app")
	createDirFromSlice(internalAppDirList)

	os.Chdir("../infrastructure")
	createDirFromSlice(infrastructureDirList)

	os.Chdir("../../")

	file, err = os.Create("Dockerfile")
	if err != nil {
		panic(err)
	}

	dockerFileContent = strings.Replace(dockerFileContent, "GO_VERSION", goVersion, 1)
	writeToFile(file, dockerFileContent)

	fmt.Printf("Direcory: %q for project: %q created!\n", dirName, *projectName)
}

func createDir(dirName string) error {
	return os.Mkdir(dirName, perm)
}

func initGoProject(projectName string) error {
	cmd := exec.Command("go", "mod", "init", projectName)
	return cmd.Run()
}

func getGoVersion() (string, error) {
	goModFile, err := os.Open("go.mod")

	if err != nil {
		return "", err
	}

	defer goModFile.Close()
	fileScanner := bufio.NewReader(goModFile)
	// Skip 2 lines
	fileScanner.ReadLine() // this line is the module name
	fileScanner.ReadLine() // this line is an empty line
	line, _, _ := fileScanner.ReadLine()

	return strings.TrimLeft(string(line), "go "), nil
}

func createDirFromSlice(dirSlice []string) {
	for _, dir := range dirSlice {
		err := createDir(dir)
		if err != nil {
			panic("failed to create directory")
		}
	}
}

func createFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func writeToFile(file *os.File, content string) error {
	defer file.Close()
	_, err := file.Write([]byte(content))
	return err
}
