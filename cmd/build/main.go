package main

import (
	"fmt"
	"os"
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
		fmt.Fprintf(os.Stderr, "go tool dist list: %v\n", err)
		os.Exit(1)
	}

	listOutputStr := string(bytes)
	listOutput := strings.Split(listOutputStr, "\n")

	for _, output := range listOutput {
		segments := strings.Split(output, "/")
		if len(segments) < 2 {
			continue
		}

		targetOs := segments[0]
		targetArch := segments[1]

		fileName := fmt.Sprintf("%s/%s_%s_%s", BUILD_DIR, BUILD_FILE, targetOs, targetArch)
		if strings.Contains(targetOs, "windows") {
			fileName += ".exe"
		}

		fmt.Printf("building %s... \n", fileName)

		buildCmd := exec.Command("go", "build", "-o", fileName, SOURCE_DIR)
		buildCmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", targetOs), fmt.Sprintf("GOARCH=%s", targetArch))

		if out, err := buildCmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "build %s/%s failed: %v\n%s", targetOs, targetArch, err, out)
			os.Exit(1)
		}
	}
}
