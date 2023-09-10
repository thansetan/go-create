package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
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
	)

	flag.Parse()

	if flag.NArg() > 0 {
		*projectName = flag.Arg(0)
	}

	splittedProjectName := strings.Split(*projectName, "/")
	dirName = splittedProjectName[len(splittedProjectName)-1]

	err = createDir(dirName)
	if err != nil {
		if strings.Contains(err.Error(), os.ErrExist.Error()) {
			fmt.Printf("Directory %s already exist!\n", dirName)
			os.Exit(1)
		}
		fmt.Println(err.Error())
		panic("failed to create directory")
	}

	err = os.Chdir(dirName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = initGoProject(*projectName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	goVersion, err := getGoVersion()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	createDirFromSlice(mainDirList)

	_, err = createFile(".env")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Chdir("cmd")

	err = createDir("app")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Chdir("app")

	mainDotGo, err := createFile("main.go")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = writeToFile(mainDotGo, mainDotGoContent)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Chdir("../../internal")

	createDirFromSlice(internalDirList)

	os.Chdir("app")
	createDirFromSlice(internalAppDirList)

	os.Chdir("../infrastructure")
	createDirFromSlice(infrastructureDirList)

	os.Chdir("../../")
	dockerFile, err := os.Create("Dockerfile")
	if err != nil {
		fmt.Println(err.Error())
	}

	dockerFileContent = strings.Replace(dockerFileContent, "GO_VERSION", goVersion, 1)
	writeToFile(dockerFile, dockerFileContent)
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
