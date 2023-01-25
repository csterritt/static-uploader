package compute

import (
	"fmt"
	"testing"
)

func TestDeletesWithNoFiles(t *testing.T) {
	var localFiles []string
	var bucketFiles []string

	deleteFiles := Deletes(localFiles, bucketFiles)
	if len(deleteFiles) != 0 {
		t.Errorf("Expected zero files to delete, got %d\n", len(deleteFiles))
	}
}

func TestDeletesWithNoBucketFiles(t *testing.T) {
	localFiles := []string{"L/b", "L/c", "L/d"}
	var bucketFiles []string

	deleteFiles := Deletes(localFiles, bucketFiles)
	if len(deleteFiles) != 0 {
		t.Errorf("Expected zero files to delete, got %d\n", len(deleteFiles))
	}
}

func TestDeletesWithNoLocalFiles(t *testing.T) {
	var localFiles []string
	bucketFiles := []string{"B/b", "B/c", "B/d"}

	deleteFiles := Deletes(localFiles, bucketFiles)
	if len(deleteFiles) != 3 {
		t.Errorf("Expected three files to delete, got %d\n", len(deleteFiles))
	}
}

func TestDeletesWithMoreLocalThanBucket(t *testing.T) {
	localFiles := []string{"L/b", "L/d"}
	bucketFiles := []string{"B/b", "B/c", "B/d"}

	deleteFiles := Deletes(localFiles, bucketFiles)
	if len(deleteFiles) != 1 {
		t.Errorf("Expected one file to delete, got %d\n", len(deleteFiles))
	}
	if len(deleteFiles) == 1 && fmt.Sprintf("%#v", deleteFiles) != "[]string{\"B/c\"}" {
		t.Errorf("Expected file to delete %s, got %s\n", "[]string{\"B/c\"}", deleteFiles[0])
	}
}

func TestDeletesWithMoreBucketThanLocal(t *testing.T) {
	localFiles := []string{"L/b", "L/c", "L/d"}
	bucketFiles := []string{"B/b", "B/d"}

	deleteFiles := Deletes(localFiles, bucketFiles)
	if len(deleteFiles) != 0 {
		t.Errorf("Expected zero files to delete, got %d\n", len(deleteFiles))
	}
}
