package file_rotator

import (
	"os"
	"time"
)

//FileInfo ...
type FileInfo struct {
	Path         string
	RotationTime time.Time
}

//Event ...
type Event struct {
	FileInfo FileInfo
}

//Handler ...
type Handler interface {
	//Init ...
	Init(dirPath string, filename string) error
	//Handle ...
	Handle(event Event)
	//Close ...
	Close() error
}

//RotationJudger ...
type Judger interface {
	//ShouldRotate ...
	ShouldRotate(file *os.File) bool
}
