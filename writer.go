package progressio

import (
	"io"
	"time"
)

type ProgressWriter struct {
	writer             io.Writer
	totalSizeBytes     int64
	transferredBytes   int64
	progressCallback   func(ProgressUpdate)
	startTime          time.Time
	intervalStartTime  time.Time
	intervalStartBytes int64
}

func NewWriter(writer io.Writer, totalSizeBytes int64, progressCallback func(ProgressUpdate)) *ProgressWriter {
	return &ProgressWriter{
		writer:             writer,
		totalSizeBytes:     totalSizeBytes,
		progressCallback:   progressCallback,
		startTime:          time.Now(),
		intervalStartTime:  time.Now(),
		intervalStartBytes: 0,
	}
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	pw.transferredBytes += int64(n)

	now := time.Now()
	intervalDuration := now.Sub(pw.intervalStartTime)
	intervalBytesTransferred := pw.transferredBytes - pw.intervalStartBytes
	intervalRateBytesPerSecond := float64(pw.transferredBytes-pw.intervalStartBytes) / intervalDuration.Seconds()
	pw.progressCallback(ProgressUpdate{
		StartTime:                pw.startTime,
		ElapsedTime:              now.Sub(pw.startTime),
		TransferredBytes:         pw.transferredBytes,
		TotalBytes:               pw.totalSizeBytes,
		IntervalStartTime:        pw.intervalStartTime,
		IntervalDuration:         intervalDuration,
		IntervalTransferredBytes: intervalBytesTransferred,
		IntervalBytesPerSecond:   intervalRateBytesPerSecond,
	})

	pw.intervalStartTime = time.Now()
	pw.intervalStartBytes = pw.transferredBytes
	return n, err
}
