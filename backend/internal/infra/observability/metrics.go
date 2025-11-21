package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	BackupLastRunTimestamp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "savesync_backup_last_run_timestamp",
			Help: "Timestamp of the last backup run",
		},
		[]string{"source_id", "source_name"},
	)

	BackupStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "savesync_backup_status",
			Help: "Status of the last backup (1=success, 0=failed)",
		},
		[]string{"source_id", "source_name"},
	)

	BackupDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "savesync_backup_duration_seconds",
			Help:    "Duration of backup operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"source_id", "source_name"},
	)

	BytesTransferredTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "savesync_bytes_transferred_total",
			Help: "Total number of bytes transferred",
		},
		[]string{"source_id", "source_name"},
	)

	ErrorCountTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "savesync_error_count_total",
			Help: "Total number of errors",
		},
		[]string{"operation"},
	)

	SnapshotCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "savesync_snapshot_count",
			Help: "Number of snapshots per source",
		},
		[]string{"source_id", "source_name"},
	)

	ChunkDeduplicationRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "savesync_chunk_deduplication_rate",
			Help: "Percentage of chunks that were deduplicated",
		},
		[]string{"source_id", "source_name"},
	)
)
