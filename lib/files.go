package lib

import "os"

func FolderExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return false
}

func RemoveFileOrFolder(path string) bool {
	err := os.RemoveAll(path)
	if err != nil {
		return false
	} else {
		return true
	}
}
