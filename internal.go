package file_rotator

import (
	"os"
	"sort"
	"time"
)

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

func (h *maxRotatableFileNumber) Handle(event Event) {
	h.rotatableFiles = append(h.rotatableFiles, event.FileInfo.Path)

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

type maxFileSizeImpl struct {
	size int64
}

func (j maxFileSizeImpl) ShouldRotate(file *os.File) bool {
	st, err := file.Stat()
	if err != nil {
		return false
	}
	return st.Size() >= j.size
}

type everyImpl struct {
	d time.Duration
	t time.Time
}

func newEvery(d time.Duration) *everyImpl {
	return &everyImpl{
		d: d,
		t: time.Now().UTC().Truncate(d),
	}
}

func (j *everyImpl) ShouldRotate(*os.File) bool {
	nowTruncate := time.Now().UTC().Truncate(j.d)
	if nowTruncate != j.t {
		j.t = nowTruncate
		return true
	}
	return false
}
