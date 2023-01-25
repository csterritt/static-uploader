package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

var neededVars = []string{"STORAGE_KEY_ID", "STORAGE_KEY_CONTENT", "STORAGE_BUCKET_NAME", "STORAGE_BUCKET_ENDPOINT"}

func isNotDir(dir string) bool {
	// This returns an *os.FileInfo type
	fileInfo, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("Cannot stat directory %s: %s\n", dir, err)
		return true
	}

	// IsDir is short for fileInfo.Mode().IsDir()
	if !fileInfo.IsDir() {
		fmt.Printf("File %s is not a directory\n", dir)
		return true
	}

	return false
}

func verifyEnvironment() {
	fail := false
	for _, name := range neededVars {
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

func main() {
	verifyEnvironment()

	bucket := aws.String(os.Getenv("STORAGE_BUCKET_NAME"))
	//key := aws.String("fileName.txt")

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(os.Getenv("STORAGE_KEY_ID"), os.Getenv("STORAGE_KEY_CONTENT"), ""),
		Endpoint:         aws.String("https://" + os.Getenv("STORAGE_BUCKET_ENDPOINT")),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	}
	//newSession := session.New(s3Config)
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		fmt.Printf("Failed to create session %s\n", err.Error())
	}

	s3Client := s3.New(newSession)

	// List files
	fmt.Printf("bucket: %s\n", *bucket)
	data, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: bucket})
	if err != nil {
		fmt.Printf("ListObjectsV2 got error %v\n", err)
	} else {
		// fmt.Printf("ListObjectsV2 returned data %#v\n", data)
		for _, content := range data.Contents {
			fmt.Printf("File: %s\n", *content.Key)
		}
	}

	// Delete an existing object
	//_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
	//	Bucket: bucket,
	//	Key:    key,
	//})
	//if err != nil {
	//	fmt.Printf("Failed to delete object %s/%s, %s\n", *bucket, *key, err.Error())
	//}

	fmt.Printf("Done!\n")
}
