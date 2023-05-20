package progress_bar

import (
	"io"
	"strings"
	"time"
)

const progressBarWidth = 20

func ProgressBar(currentSize, totalSize int64) string {
	progress := int(float64(currentSize) / float64(totalSize) * float64(progressBarWidth))
	return strings.Repeat("=", progress) + strings.Repeat("-", progressBarWidth-progress)
}

// LimitedReader limits the read speed by introducing delays between reads
type LimitedReader struct {
	r            io.Reader
	totalBytes   *int64
	isLimited    bool
	time         time.Time
	timeInterval int64
	speed        *float64
}

// NewLimitedReader creates a new LimitedReader with the specified read speed
func NewLimitedReader(r io.Reader, rateLimit int, totalBytes *int64, speed *float64) *LimitedReader {
	return &LimitedReader{
		r:            r,
		isLimited:    rateLimit > 0,
		totalBytes:   totalBytes,
		time:         time.Now(),
		timeInterval: 0,
		speed:        speed,
	}
}

// Close closes the underlying reader if it implements the io.Closer interface
func (lr *LimitedReader) Close() error {
	if closer, ok := lr.r.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Read reads data from the underlying reader with a limited speed
func (lr *LimitedReader) Read(p []byte) (n int, err error) {
	n, err = lr.r.Read(p)
	timeNow := time.Now()

	lr.timeInterval = int64(timeNow.Sub(lr.time))
	*lr.speed = float64(n) / (float64(lr.timeInterval) / float64(time.Second))
	lr.time = timeNow
	*lr.totalBytes += int64(n)

	if lr.isLimited {
		time.Sleep(time.Second)
	}

	return n, err
}
