package compute

func Deletes(localFiles []string, bucketFiles []string) []string {
	res := make([]string, 0)
	locals := make(map[string]bool, 0)
	for _, localFile := range localFiles {
		locals[localFile] = true
	}

	for _, bucketFile := range bucketFiles {
		if _, found := locals[bucketFile]; !found {
			res = append(res, bucketFile)
		}
	}

	return res
}
