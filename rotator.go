package file_rotator

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

//RotationHandler ...
type RotationHandler interface {
	//Init ...
	Init(dirPath string, filename string) error
	//Handle ...
	Handle(rotatableFilepath string)
	//Close ...
	Close() error
}

type fileRotator struct {
	currentFileSize int64
	currentFile     *os.File
	mu              *sync.Mutex
	filepath        string
	maxFileSize     int64
	rotationHandler RotationHandler
	filenameSuffix  func() string
}

//New ...
func New(dirPath string, filename string, maxFileSize int64, opts ...Option) (io.WriteCloser, error) {
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

	st, err := newFile.Stat()
	if err != nil {
		return nil, err
	}

	var rotator = &fileRotator{
		currentFile:     newFile,
		currentFileSize: st.Size(),
		mu:              new(sync.Mutex),
		filepath:        absFilepath,
		maxFileSize:     maxFileSize,
		filenameSuffix: func() string {
			return strconv.FormatInt(time.Now().UnixNano(), 10)
		},
	}

	for _, opt := range opts {
		opt(rotator)
	}
	if rotator.rotationHandler != nil {
		err = rotator.rotationHandler.Init(dirPath, filename)
		if err != nil {
			return nil, err
		}
	}
	return rotator, nil
}

func (r *fileRotator) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.currentFileSize+int64(len(p)) >= r.maxFileSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = r.currentFile.Write(p)
	r.currentFileSize += int64(n)
	return
}

func (r fileRotator) Close() error {
	if err := r.currentFile.Close(); err != nil {
		return err
	}
	if r.rotationHandler != nil {
		if err := r.rotationHandler.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (r *fileRotator) rotate() error {
	err := r.currentFile.Close()
	if err != nil {
		return err
	}

	currentFilepath, err := filepath.Abs(r.currentFile.Name())
	if err != nil {
		return err
	}

	rotatableFilepath := r.filepath + "." + r.filenameSuffix()

	err = os.Rename(currentFilepath, rotatableFilepath)
	if err != nil {
		return err
	}

	newFile, err := openNewFile(r.filepath)
	if err != nil {
		return err
	}

	st, err := newFile.Stat()
	if err != nil {
		return err
	}

	r.currentFile = newFile
	r.currentFileSize = st.Size()

	if r.rotationHandler != nil {
		r.rotationHandler.Handle(rotatableFilepath)
	}

	return nil
}
