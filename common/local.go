package common

import "os"

func CheckLocalFile(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		return false
	}

	return true
}
