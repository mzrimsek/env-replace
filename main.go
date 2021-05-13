package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) == 2 {
		pwd, err := os.Getwd()
		checkError(err, "Failed to get pwd location")

		var envFilePath = filepath.Join(pwd, ".env")
		var fileToProcessPath = filepath.Join(pwd, os.Args[1])

		envMap := getEnvMap(envFilePath)
		var processedFile = getProcessedFile(fileToProcessPath, envMap)
		fmt.Println(processedFile)
	} else {
		printError("Invalid parameters")
	}
}

func printError(msg string) {
	fmt.Printf("Error: %s\n", msg)
}

func checkError(e error, msg string) {
	if e != nil {
		printError(msg)
	}
}

func getEnvMap(envFilePath string) map[string]string {
	envFile, err := os.Open(envFilePath)
	checkError(err, "No .env file found in pwd or unable to open")

	envMap := make(map[string]string)

	scanner := bufio.NewScanner(envFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		var envPair = scanner.Text()
		if envPair != "" {
			var pair = strings.Split(envPair, "=")

			var key = pair[0]
			var value = pair[1]
			envMap[key] = value
		}
	}

	envFile.Close()
	return envMap
}

func getProcessedFile(fileToProcessPath string, envMap map[string]string) string {
	fileToProcess, err := os.Open(fileToProcessPath)
	checkError(err, "Target file not found in pwd")

	scanner := bufio.NewScanner(fileToProcess)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		var line = scanner.Text()

		for key := range envMap {
			var testString = fmt.Sprintf("${%s}", key)
			if strings.Contains(scanner.Text(), testString) {
				line = strings.ReplaceAll(line, testString, envMap[key])
			}
		}

		lines = append(lines, line)
	}

	fileToProcess.Close()
	return strings.Join(lines, "\n")
}
