package main

import (
	"fmt"
	"os"
)

var neededEnvVars = []string{"STORAGE_KEY_ID", "STORAGE_KEY_CONTENT", "STORAGE_BUCKET_NAME", "STORAGE_BUCKET_ENDPOINT"}

func isNotDir(dir string) bool {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("Cannot stat directory %s: %s\n", dir, err)
		return true
	}

	if !fileInfo.IsDir() {
		fmt.Printf("File %s is not a directory\n", dir)
		return true
	}

	return false
}

func verifyEnvironment() {
	fail := false
	for _, name := range neededEnvVars {
		val := os.Getenv(name)
		if val == "" {
			fail = true
			fmt.Printf("Error: environmental variable %s is not set\n", name)
		}
	}

	if fail {
		os.Exit(1)
	}

	if len(os.Args) != 2 || isNotDir(os.Args[1]) {
		fmt.Printf("Usage: static-upload distribution-directory\n")
		os.Exit(1)
	}
}

func findAllFiles(localDir string) ([]string, error) {
	dirs, err := os.ReadDir(localDir)
	if err != nil {
		fmt.Printf("Cannot read files in %s: %v\n", localDir, err)
		return nil, err
	}

	res := make([]string, 0)
	for _, file := range dirs {
		if file.IsDir() {
			subFiles, err := findAllFiles(localDir + "/" + file.Name())
			if err != nil {
				return nil, err
			}

			res = append(res, subFiles...)
		} else {
			res = append(res, localDir+"/"+file.Name())
		}
	}

	return res, nil
}

func copyLocalToBucket(localDir string, bucket string) {
	localFiles, err := findAllFiles(localDir)
	if err != nil {
		fmt.Printf("Cannot read localFiles in %s: %v\n", localDir, err)
		os.Exit(1)
	}

	dirLen := len(localDir) + 1
	for _, file := range localFiles {
		fmt.Printf("%s\n", file[dirLen:])
	}
}

func main() {
	verifyEnvironment()

	bucket := os.Getenv("STORAGE_BUCKET_NAME")
	localDir := os.Args[1]

	copyLocalToBucket(localDir, bucket)

	fmt.Printf("Done!\n")
}
