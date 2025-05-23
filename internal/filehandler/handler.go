package handler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const templatesPath string = "internal/filehandler/templates/"

var rootFolder string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getRootFolder() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	rootFolder = filepath.Base(cwd)

}

func createFile(path string, outputName string) *os.File {
	fullFilePath := fmt.Sprintf("%s%s.go", path, outputName)

	outFile, err := os.Create(fullFilePath)
	check(err)

	return outFile
}

func createDirAndFile(path string, outputName string) *os.File {

	// Create the dir
	err := os.Mkdir(fmt.Sprintf("%s%s", path, outputName), os.ModePerm)
	check(err)

	// Create the file
	outFile, err := os.Create(fmt.Sprintf("%s%s/%s.go", path, outputName, outputName))
	check(err)

	return outFile
}

func generateFilesFromTemplates(createTyp string, selectTyp string, availParams []string, routNam string, needSubFolder bool) {
	var templateParams = map[string]string{}
	var outPutFile *os.File

	for _, param := range availParams {
		trimmedValue := strings.Trim(param, "{}")

		switch strings.ToLower(trimmedValue) {
		case "name":
			newRouteName := routNam

			// Check if the first letter of the original template was capitalized
			capitalize := false

			if unicode.IsUpper(rune(trimmedValue[0])) {
				capitalize = true
			}

			if capitalize {
				newRouteName = strings.ToUpper(newRouteName[:1]) + newRouteName[1:]
			}

			templateParams[param] = newRouteName
		case "rootfold":
			getRootFolder()
			templateParams[param] = rootFolder
		default:
			panic(fmt.Errorf("invalid command: %s", trimmedValue))
		}

	}

	rootPath, err := os.Getwd()
	check(err)

	var templatesFullPath string = filepath.Join(rootPath, templatesPath)
	
	
	if needSubFolder {
		outPutFile = createDirAndFile(fmt.Sprintf("config/%s/", createTyp), routNam)
	} else {
		outPutFile = createFile(fmt.Sprintf("config/%s/", createTyp), routNam)
	}

	entries, err := os.ReadDir(templatesFullPath)
	check(err)

	var templateFiles []string

	for _, entry := range entries {
		templateFiles = append(templateFiles, entry.Name())
	}

	writer := bufio.NewWriter(outPutFile)

	file, err := os.Open(templatesPath + selectTyp)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		for key, value := range templateParams {
			line = strings.ReplaceAll(line, key, value)

		}
		_, err := writer.WriteString(line + "\n")
		check(err)

		// Close the input file when done
		defer file.Close()

	}

	// Close the output file when done
	defer outPutFile.Close()
	writer.Flush() // Ensure all buffered writes are written to the file
}

func CreateRoute(routeName string) {
	// TODO This can be created into a func and also to be added in create handler
	if routeName == "" {
		fmt.Errorf("No name specified for the route")
		os.Exit(1)
	}

	var selectedTemplate string = "route"
	var createTypes string = selectedTemplate + "s"
	var templateAvailParams = []string{"{name}", "{Name}", "{rootfold}"}

	// Only thing left is to tell it which template to select
	generateFilesFromTemplates(createTypes, selectedTemplate, templateAvailParams, routeName, true)
}

func CreateHandler(handlerName string) {

	if handlerName == "" {
		fmt.Errorf("No name specified for the route")
		os.Exit(1)
	}

	var selectedTemplate string = "handler"
	var createTypes string = selectedTemplate + "s"
	var templateAvailParams = []string{"{name}", "{Name}"}

	generateFilesFromTemplates(createTypes, selectedTemplate, templateAvailParams, handlerName, false)

}
