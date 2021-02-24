package utils

import (
	"strings"
)

// Check if a given string vector has at least one string
func Any(v []string, k string) (bool, int) {
	for i, str := range v {
		if str == k {
			return true, i
		}
	}

	return false, -1
}

func ExtractFilename(contentDispositon []string, url string) string {

	if contentDispositon != nil {
		mark := "filename="
		if len(contentDispositon) == 1 && strings.Contains(contentDispositon[0], mark) {
			filenameIndex := strings.Index(contentDispositon[0], mark)
			if filenameIndex > 0 {
				filename := contentDispositon[0][filenameIndex+len(mark):]
				if filename[0] == '"' && filename[len(filename)-1] == '"' {
					filename = filename[1 : len(filename)-1]
				}

				return filename
			}
		}
	}

	splited := strings.Split(url, "/")
	return splited[len(splited)-1]
}
