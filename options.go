package file_rotator

import (
	"os"
	"sort"
)

//Option ...
type Option func(rotator *fileRotator)

//SetRotationHandler ...
func SetRotationHandler(handler RotationHandler) Option {
	return func(r *fileRotator) {
		r.rotationHandler = handler
	}
}

//SetFilenameSuffix ...
func SetFilenameSuffix(suffixFunc func() string) Option {
	return func(r *fileRotator) {
		r.filenameSuffix = suffixFunc
	}
}

type maxRotatableFileNumber struct {
	rotatableFiles []string
	maxNum         int
}

func (h *maxRotatableFileNumber) Init(dirPath string, filename string) error {
	files, err := GetRotatableFiles(dirPath, filename)
	if err != nil {
		return err
	}
	sort.Strings(files)
	h.rotatableFiles = files
	return nil
}

func (h *maxRotatableFileNumber) Handle(rotatableFilepath string) {
	h.rotatableFiles = append(h.rotatableFiles, rotatableFilepath)

	for len(h.rotatableFiles) > h.maxNum {
		delFilepath := h.rotatableFiles[0]
		err := os.Remove(delFilepath)
		if err != nil {
			return
		}
		h.rotatableFiles = h.rotatableFiles[1:]
	}
}

func (h *maxRotatableFileNumber) Close() error {
	return nil
}

//MaxRotatableFileNumber ...
func MaxRotatableFileNumber(num int) Option {
	var handler = &maxRotatableFileNumber{maxNum: num}
	return func(r *fileRotator) {
		r.rotationHandler = handler
	}
}
