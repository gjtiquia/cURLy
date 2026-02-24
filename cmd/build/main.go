package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const SOURCE_DIR = "./cmd/tui"
const BUILD_DIR = "./bin"
const BUILD_FILE = "cURLy"

func main() {
	fmt.Println("building... dir:", BUILD_DIR)

	err := os.RemoveAll("bin")
	if err != nil {
		log.Panicf("os.RemoveAll('bin'): %v\n", err)
	}

	listCmd := exec.Command("go", "tool", "dist", "list")
	bytes, err := listCmd.Output()
	if err != nil {
		log.Panicf("go tool dist list: %v\n", err)
	}

	listOutputStr := string(bytes)
	listOutput := strings.Split(listOutputStr, "\n")

	var wg sync.WaitGroup
	for _, output := range listOutput {
		segments := strings.Split(output, "/")
		if len(segments) < 2 {
			continue
		}

		targetOs := segments[0]
		targetArch := segments[1]

		wg.Add(1)
		go build(targetOs, targetArch, &wg)
	}

	wg.Wait()
	fmt.Println("build complete")
}

func build(targetOs, targetArch string, wg *sync.WaitGroup) {
	defer wg.Done()

	fileName := fmt.Sprintf("%s/%s_%s_%s", BUILD_DIR, BUILD_FILE, targetOs, targetArch)
	if strings.Contains(targetOs, "windows") {
		fileName += ".exe"
	}

	fmt.Printf("building %s... \n", fileName)

	buildCmd := exec.Command("go", "build", "-o", fileName, SOURCE_DIR)
	buildCmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", targetOs),
		fmt.Sprintf("GOARCH=%s", targetArch))

	if out, err := buildCmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build %s/%s failed: %v\n%s", targetOs, targetArch, err, out)
		return
	}

	fmt.Printf("build %s... succeeded!\n", fileName)
}
