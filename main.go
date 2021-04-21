package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	envMap := getEnvMap(".env")
	var processedYamlFile = getProcessedYaml("test.yaml", envMap)
	fmt.Print(processedYamlFile)
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
		var pair = strings.Split(envPair, "=")

		var key = pair[0]
		var value = pair[1]
		envMap[key] = value
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
