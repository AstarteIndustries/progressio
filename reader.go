package progressio

import (
	"io"
	"time"
)

type ProgressReader struct {
	reader             io.Reader
	totalSizeBytes     int64
	transferredBytes   int64
	progressCallback   func(ProgressUpdate)
	startTime          time.Time
	intervalStartTime  time.Time
	intervalStartBytes int64
}

func NewReader(reader io.Reader, totalSizeBytes int64, progressCallback func(ProgressUpdate)) *ProgressReader {
	return &ProgressReader{
		reader:             reader,
		totalSizeBytes:     totalSizeBytes,
		progressCallback:   progressCallback,
		startTime:          time.Now(),
		intervalStartTime:  time.Now(),
		intervalStartBytes: 0,
	}
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if n <= 0 {
		return n, err
	}
	pr.transferredBytes += int64(n)

	now := time.Now()
	intervalDuration := now.Sub(pr.intervalStartTime)
	intervalBytesTransferred := pr.transferredBytes - pr.intervalStartBytes
	intervalRateBytesPerSecond := float64(pr.transferredBytes-pr.intervalStartBytes) / intervalDuration.Seconds()
	pr.progressCallback(ProgressUpdate{
		StartTime:                pr.startTime,
		ElapsedTime:              now.Sub(pr.startTime),
		TransferredBytes:         pr.transferredBytes,
		TotalBytes:               pr.totalSizeBytes,
		IntervalStartTime:        pr.intervalStartTime,
		IntervalDuration:         intervalDuration,
		IntervalTransferredBytes: intervalBytesTransferred,
		IntervalBytesPerSecond:   intervalRateBytesPerSecond,
	})

	pr.intervalStartTime = time.Now()
	pr.intervalStartBytes = pr.transferredBytes
	return n, err
}
