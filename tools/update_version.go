package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		log.Fatalf("failed to find git on path: %v", err)
	}

	cmd := exec.Command(gitPath, "rev-list", "master", "--count")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to run git: %v", err)
	}
	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		log.Fatalf("failed to parse git output: %v", err)
	}

	newVers := fmt.Sprintf("0.0.%d", count+1)

	verFile, err := os.Create("cmd/.version")
	if err != nil {
		log.Fatalf("failed to parse git output: %v", err)
	}

	if _, err := verFile.WriteString(newVers); err != nil {
		log.Fatalf("failed to parse git output: %v", err)
	}

	if err := verFile.Close(); err != nil {
		log.Fatalf("failed to parse git output: %v", err)
	}
}
