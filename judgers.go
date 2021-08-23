package file_rotator

import "time"

func MaxFileSize(size int64) Judger {
	return maxFileSizeImpl{size: size}
}

func Every(d time.Duration) Judger {
	return newEvery(d)
}
