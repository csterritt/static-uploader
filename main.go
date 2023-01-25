package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jawher/mow.cli"
)

var neededEnvVars = []string{"STORAGE_KEY_ID", "STORAGE_KEY_CONTENT", "STORAGE_BUCKET_ENDPOINT"}

func isNotDir(dir string) bool {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("Cannot stat directory %s: %s\n", dir, err)
		return true
	}

	if !fileInfo.IsDir() {
		return true
	}

	return false
}

func verifyEnvironment() error {
	for _, name := range neededEnvVars {
		val := os.Getenv(name)
		if val == "" {
			return errors.New(fmt.Sprintf("Error: environmental variable %s is not set\n", name))
		}
	}

	return nil
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

func getS3Client() (*s3.S3, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(os.Getenv("STORAGE_KEY_ID"), os.Getenv("STORAGE_KEY_CONTENT"), ""),
		Endpoint:         aws.String("https://" + os.Getenv("STORAGE_BUCKET_ENDPOINT")),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession, err := session.NewSession(s3Config)
	if err != nil {
		fmt.Printf("Failed to create session %s\n", err.Error())
		return nil, err
	}

	return s3.New(newSession), nil
}

func findBucketFiles(s3Client *s3.S3, bucket string) ([]string, error) {
	res := make([]string, 0)

	fmt.Printf("bucket: %s\n", bucket)
	data, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		fmt.Printf("ListObjectsV2 got error %v\n", err)
		return nil, err
	} else {
		for _, content := range data.Contents {
			res = append(res, *content.Key)
		}
	}

	return res, nil
}

func copyLocalToBucket(s3Client *s3.S3, localDir string, bucket string, testOnly bool) {
	localFiles, err := findAllFiles(localDir)
	if err != nil {
		fmt.Printf("Cannot read local files in %s: %v\n", localDir, err)
		os.Exit(1)
	}

	dirLen := len(localDir) + 1
	for _, file := range localFiles {
		if testOnly {
			fmt.Printf("L: %s\n", file[dirLen:])
		}
	}

	if testOnly {
		fmt.Printf("--------\n")
	}

	bucketFiles, err := findBucketFiles(s3Client, bucket)
	if err != nil {
		fmt.Printf("Cannot read bucket files in %s: %v\n", bucket, err)
		os.Exit(1)
	}

	bucketLen := len(bucket) + 1
	for _, file := range bucketFiles {
		if testOnly {
			fmt.Printf("B: %s\n", file[bucketLen:])
		}
	}
}

func main() {
	// create an app
	app := cli.App("static-uploader", "Copy files from a local directory to an S3-type bucket")
	app.Spec = "[-t] LOCAL_DIRECTORY BUCKET"

	var (
		// declare the -r flag as a boolean flag
		testOnly = app.BoolOpt("t test", false, "Don't copy the files, just show what would be done")
		// declare the SRC argument as a multi-string argument
		localDir = app.StringArg("LOCAL_DIRECTORY", "", "Local directory to copy from")
		// declare the DST argument as a single string (string slice) arguments
		bucket = app.StringArg("BUCKET", "", "Destination bucket to copy files to")
	)

	// Specify the action to execute when the app is invoked correctly
	app.Action = func() {
		err := verifyEnvironment()
		if err != nil {
			fmt.Printf("Got error %v\n", err)
			os.Exit(1)
		}

		if isNotDir(*localDir) {
			fmt.Printf("Error: %s is not a directory.\n", *localDir)
			os.Exit(1)
		}

		s3Client, err := getS3Client()
		if err != nil {
			fmt.Printf("Cannot create S3 client: %v\n", err)
			os.Exit(1)
		}

		copyLocalToBucket(s3Client, *localDir, *bucket, *testOnly)

		fmt.Printf("Done!\n")
	}

	// Invoke the app passing in os.Args
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Got error %v\n", err)
		os.Exit(1)
	}
}
