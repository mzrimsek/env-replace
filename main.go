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
		checkError(err)

		var envFilePath = filepath.Join(pwd, ".env")
		var yamlFilePath = filepath.Join(pwd, os.Args[1])

		envMap := getEnvMap(envFilePath)
		var processedYamlFile = getProcessedYaml(yamlFilePath, envMap)
		fmt.Println(processedYamlFile)
	} else {
		panic("Invalid parameters")
	}
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func getEnvMap(envFilePath string) map[string]string {
	envFile, err := os.Open(envFilePath)
	checkError(err)

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

func getProcessedYaml(yamlFilePath string, envMap map[string]string) string {
	yamlFile, err := os.Open(yamlFilePath)
	checkError(err)

	scanner := bufio.NewScanner(yamlFile)
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

	yamlFile.Close()
	return strings.Join(lines, "\n")
}
