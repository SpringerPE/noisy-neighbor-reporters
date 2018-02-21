package reporter

import (
	graphite "github.com/marpaia/graphite-golang"

	"log"
	"time"
)

// Reporter stores configuration for reporting to Graphite.
type GraphiteReporter struct {
	pointBuilder   PointBuilder
	graphiteClient GraphiteClient
	interval       time.Duration
	metricsPrefix  string
}

// NewReporter initializes and returns a new Reporter.
func NewReporter(pointBuilder PointBuilder, graphiteClient GraphiteClient, opts ...ReporterOption) *GraphiteReporter {

	r := &GraphiteReporter{
		pointBuilder:   pointBuilder,
		graphiteClient: graphiteClient,
		interval:       time.Minute,
	}

	for _, o := range opts {
		o(r)
	}

	return r
}

// Run reports metrics from the configured PointBuilder to Graphite on a
// configured interval.
func (r *GraphiteReporter) Run() {

	ticker := time.NewTicker(r.interval)

	for timestamp := range ticker.C {
		func() {
			log.Printf("graphite reporter ticked at %s", timestamp)

			points, err := r.getAllPoints()
			if err != nil {
				log.Printf("failed to build points from points builder: %s", err)
				return
			}

			err = r.graphiteClient.Connect()
			defer func() {
				err = r.graphiteClient.Disconnect()
				if err != nil {
					log.Printf("Failed disconnecting from graphite: %s", err)
				}
			}()

			if err != nil {
				log.Printf("Failed connecting to graphite: %s", err)
			}

			err = r.graphiteClient.SendMetrics(points)
			if err != nil {
				log.Printf("failed to post to graphite: %s", err)
				return
			}

		}()
	}

}

func (r *GraphiteReporter) getAllPoints() (points []graphite.Metric, err error) {
	ts := time.Now().
		Add(-2 * r.interval).
		Truncate(r.interval).
		Unix()

	points, err = r.pointBuilder.BuildPoints(ts)
	if err != nil {
		return nil, err
	}

	return points, nil
}

// PointBuilder is the interface the GraphiteReporter will use to collect
// metrics to send to Graphite.
type PointBuilder interface {
	BuildPoints(int64) ([]graphite.Metric, error)
}

// ReporterOption is a func that is used to configure optional settings on a
// GraphiteReporter.
type ReporterOption func(*GraphiteReporter)

// WithInterval returns a ReporterOption for configuring the interval metrics
// will be reported to Graphite.
func WithInterval(d time.Duration) ReporterOption {
	return func(r *GraphiteReporter) {
		r.interval = d
	}
}

// GraphiteClient is the interface used for sending requests to Graphite.
type GraphiteClient interface {
	SendMetrics([]graphite.Metric) error
	Connect() error
	Disconnect() error
}
