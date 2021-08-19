package file_rotator

import (
	"os"
	"path/filepath"
	"strings"
)

func checkOrCreateDir(dirPath string) (err error) {
	absPathToDir, err := filepath.Abs(dirPath)
	if err != nil {
		return err
	}

	dir, err := os.OpenFile(absPathToDir, os.O_RDONLY, 0666)
	if err != nil {
		return os.MkdirAll(absPathToDir, os.ModePerm)
	}
	return dir.Close()
}

func openNewFile(absFilepath string) (f *os.File, err error) {
	return os.OpenFile(absFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}

//GetRotatableFiles ...
func GetRotatableFiles(dirPath string, filename string) (rotatableFiles []string, err error) {
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(absDirPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		walkFilename := info.Name()
		if strings.Contains(walkFilename, filename) && walkFilename != filename {
			rotatableFiles = append(rotatableFiles, path)
		}

		return nil
	})

	return
}
