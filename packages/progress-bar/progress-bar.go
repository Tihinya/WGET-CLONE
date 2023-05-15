package progress_bar

import (
	"io"
	"strings"
	"time"
)

const progressBarWidth = 40

func ProgressBar(currentSize, totalSize int64) string {
	progress := int(float64(currentSize) / float64(totalSize) * float64(progressBarWidth))
	return strings.Repeat("=", progress) + strings.Repeat("-", progressBarWidth-progress)
}

// LimitedReader limits the read speed by introducing delays between reads
type LimitedReader struct {
	r          io.Reader
	totalBytes *int64
	isLimited  bool
}

// NewLimitedReader creates a new LimitedReader with the specified read speed
func NewLimitedReader(r io.Reader, rateLimit int, totalBytes *int64) *LimitedReader {
	return &LimitedReader{
		r:          r,
		isLimited:  rateLimit > 0,
		totalBytes: totalBytes,
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

	*lr.totalBytes += int64(n)

	if lr.isLimited {
		time.Sleep(time.Second)
	}

	return n, err
}
