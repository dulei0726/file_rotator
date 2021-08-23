package file_rotator

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

//FileRotator ...
type FileRotator struct {
	mu             sync.Mutex
	filepath       string
	currentFile    *os.File
	judgers        []Judger
	handler        Handler
	filenameSuffix func(rotationTime time.Time) string
}

//New ...
func New(dirPath string, filename string, opts ...Option) (*FileRotator, error) {
	err := checkOrCreateDir(dirPath)
	if err != nil {
		return nil, err
	}

	absFilepath, err := filepath.Abs(path.Join(dirPath, filename))
	if err != nil {
		return nil, err
	}

	newFile, err := openNewFile(absFilepath)
	if err != nil {
		return nil, err
	}

	var rotator = &FileRotator{
		mu:          sync.Mutex{},
		filepath:    absFilepath,
		currentFile: newFile,
		filenameSuffix: func(rotationTime time.Time) string {
			return strconv.FormatInt(rotationTime.UnixNano(), 10)
		},
	}

	for _, opt := range opts {
		opt(rotator)
	}
	if rotator.handler != nil {
		err = rotator.handler.Init(dirPath, filename)
		if err != nil {
			return nil, err
		}
	}
	return rotator, nil
}

func (r *FileRotator) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, judger := range r.judgers {
		if judger.ShouldRotate(r.currentFile) {
			if err := r.rotate(); err != nil {
				return 0, err
			}
			break
		}
	}

	return r.currentFile.Write(p)
}

func (r *FileRotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.currentFile.Close(); err != nil {
		return err
	}
	if r.handler != nil {
		if err := r.handler.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (r *FileRotator) Rotate() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rotate()
}

func (r *FileRotator) rotate() error {
	err := r.currentFile.Close()
	if err != nil {
		return err
	}

	rotationTime := time.Now().UTC()
	rotationFilepath := r.filepath + "." + r.filenameSuffix(rotationTime)

	err = os.Rename(r.filepath, rotationFilepath)
	if err != nil {
		return err
	}

	newFile, err := openNewFile(r.filepath)
	if err != nil {
		return err
	}

	r.currentFile = newFile

	if r.handler != nil {
		var fileInfo = FileInfo{
			Path:         rotationFilepath,
			RotationTime: rotationTime,
		}
		r.handler.Handle(Event{FileInfo: fileInfo})
	}

	return nil
}
