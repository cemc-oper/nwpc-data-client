package data_client

import "os"

func CheckLocalFile(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return true
	}

	return false
}
