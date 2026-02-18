package main

import (
	"fmt"
	"os/exec"
	"strings"
)

const SOURCE_DIR = "./cmd/tui"
const BUILD_DIR = "./bin"
const BUILD_FILE = "cURLy"

func main() {

	fmt.Println("building... dir:", BUILD_DIR)

	listCmd := exec.Command("go", "tool", "dist", "list")
	bytes, err := listCmd.Output()
	if err != nil {
		panic(err)
	}

	listOutputStr := string(bytes)
	listOutput := strings.Split(listOutputStr, "\n")

	for _, output := range listOutput {
		segments := strings.Split(output, "/")
		if len(segments) < 2 {
			continue
		}

		os := segments[0]
		arch := segments[1]

		fileName := fmt.Sprintf("%s/%s_%s_%s", BUILD_DIR, BUILD_FILE, os, arch)
		if strings.Contains(os, "windows") {
			fileName += ".exe"
		}

		fmt.Printf("building %s... \n", fileName)

		buildCmd := exec.Command("go", "build", "-o", fileName, SOURCE_DIR)
		buildCmd.Env = append(buildCmd.Env, fmt.Sprintf("GOOS=%s", os), fmt.Sprintf("GOARCH=%s", arch))

		err := buildCmd.Run()
		if err != nil {
			panic(err)
		}
	}

}
