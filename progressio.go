package progressio

import "time"

type ProgressUpdate struct {
	StartTime                time.Time
	ElapsedTime              time.Duration
	TransferredBytes         int64
	TotalBytes               int64
	IntervalStartTime        time.Time
	IntervalDuration         time.Duration
	IntervalTransferredBytes int64
	IntervalBytesPerSecond   float64
}
