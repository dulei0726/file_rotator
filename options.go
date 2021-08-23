package file_rotator

import (
	"time"
)

//Option ...
type Option func(rotator *FileRotator)

//SetRotationHandler ...
func SetRotationHandler(handler Handler) Option {
	return func(r *FileRotator) {
		r.handler = handler
	}
}

//SetFilenameSuffix ...
func SetFilenameSuffix(suffixFunc func(rotationTime time.Time) string) Option {
	return func(r *FileRotator) {
		r.filenameSuffix = suffixFunc
	}
}

func SetRotationJudgers(judges []Judger) Option {
	return func(r *FileRotator) {
		r.judgers = judges
	}
}

func AddRotationJudger(judger ...Judger) Option {
	return func(r *FileRotator) {
		r.judgers = append(r.judgers, judger...)
	}
}

//MaxRotatableFileNumber ...
func MaxRotatableFileNumber(num int) Option {
	var handler = &maxRotatableFileNumber{maxNum: num}
	return func(r *FileRotator) {
		r.handler = handler
	}
}
