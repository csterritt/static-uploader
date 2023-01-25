package compute

import "strings"

func Deletes(localFiles []string, bucketFiles []string) []string {
	res := make([]string, 0)
	locals := make(map[string]bool, 0)
	for _, localFile := range localFiles {
		file := localFile
		index := strings.Index(file, "/")
		if index > -1 {
			file = file[index+1:]
		}
		locals[file] = true
	}

	for _, bucketFile := range bucketFiles {
		bucket := bucketFile
		index := strings.Index(bucket, "/")
		if index > -1 {
			bucket = bucket[index+1:]
		}
		if _, found := locals[bucket]; !found {
			res = append(res, bucketFile)
		}
	}

	return res
}
