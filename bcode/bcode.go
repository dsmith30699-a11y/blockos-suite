package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	projectsDir = "/home/kaung/Desktop/Block OS/projects"
	vmUser      = "root"
	vmHost      = "localhost"
	vmPort      = "2222"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	command := os.Args[1]
	switch command {
	case "list":
		listProjects()
	case "build":
		if len(os.Args) < 3 { usage(); return }
		buildProject(os.Args[2])
	case "deploy":
		if len(os.Args) < 3 { usage(); return }
		deployProject(os.Args[2])
	case "run":
		if len(os.Args) < 3 { usage(); return }
		runProject(os.Args[2], os.Args[3:]...)
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage: bcode {list|build|deploy|run} [project_name]")
	fmt.Println("  list            - Show all custom projects")
	fmt.Println("  build <name>    - Compile Go project for Block OS")
	fmt.Println("  deploy <name>   - Build and SCP binary to VM")
	fmt.Println("  run <name>      - Build, Deploy, and SSH Execute")
}

func listProjects() {
	entries, _ := os.ReadDir(projectsDir)
	fmt.Printf("%-15s %-20s\n", "PROJECT", "PATH")
	fmt.Println("---------------------------------------------------")
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("%-15s %s\n", entry.Name(), filepath.Join(projectsDir, entry.Name()))
		}
	}
}

func buildProject(name string) {
	fmt.Printf("Building [%s] for Block OS...\n", name)
	projPath := filepath.Join(projectsDir, name)
	srcFile := filepath.Join(projPath, name+".go")
	
	if _, err := os.Stat(srcFile); os.IsNotExist(err) {
		fmt.Printf("Error: %s.go not found in %s\n", name, projPath)
		return
	}

	cmd := exec.Command("go", "build", "-o", name, srcFile)
	cmd.Dir = projPath
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Build failed:\n%s\n", string(out))
		return
	}
	fmt.Println("Build successful!")
}

func deployProject(name string) {
	buildProject(name)
	fmt.Printf("Deploying [%s] to Block OS (port %s)...\n", name, vmPort)
	
	binaryPath := filepath.Join(projectsDir, name, name)
	remotePath := fmt.Sprintf("%s@%s:/usr/bin/%s", vmUser, vmHost, name)
	
	cmd := exec.Command("scp", "-P", vmPort, binaryPath, remotePath)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Deploy failed:\n%s\n", string(out))
		return
	}
	fmt.Println("Deploy successful!")
}

func runProject(name string, args ...string) {
	deployProject(name)
	fmt.Printf("Running [%s] on Block OS...\n", name)
	
	sshArgs := []string{"-p", vmPort, "-t", fmt.Sprintf("%s@%s", vmUser, vmHost), name}
	sshArgs = append(sshArgs, args...)
	
	cmd := exec.Command("ssh", sshArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
