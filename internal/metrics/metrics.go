package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const namespace = "report_service"

const (
	level       = "level"
	message     = "message"
	status      = "status"
	destination = "destination"
	errName     = "error"
)

const (
	Success = "success"
	Failed  = "failed"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "reports_service_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})

	logs = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "logs",
			Name:      "count",
			Help:      "Log count",
		},
		[]string{level, message},
	)

	processedReports = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "reports",
			Name:      "processed_reports",
			Help:      "Processed reports logs count",
		},
		[]string{status, destination},
	)

	kafkaProcessedMsgs = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "kafka",
		Name:      "kafka_processed_msgs",
		Help:      "Processed kafka msg count",
	},
		[]string{status, errName},
	)
)

func LogsInc(lvl string, msg string) {
	logs.WithLabelValues(lvl, msg).Inc()
}

func ProcessedReportsInc(status string, destination string) {
	processedReports.WithLabelValues(status, destination).Inc()
}

func KafkaProcessedMsgsInc(status string, errName string) {
	kafkaProcessedMsgs.WithLabelValues(status, errName).Inc()
}
