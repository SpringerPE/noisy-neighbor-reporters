package reporter_test

import (
	"sync"
	"time"

	graphite "github.com/marpaia/graphite-golang"

	"github.com/SpringerPE/noisy-neighbor-reporters/pkg/reporter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GraphiteReporter", func() {
	It("sends data points to graphite on an interval", func() {
		pointBuilder := &spyPointBuilder{}
		graphiteClient := &spyGraphiteClient{}

		reporter := reporter.NewReporter(
			pointBuilder, graphiteClient,
			reporter.WithInterval(50*time.Millisecond),
		)
		go reporter.Run()

		Eventually(pointBuilder.buildCalled).Should(BeNumerically(">", 1))
		Expect(pointBuilder.buildPointsTimestamp()).To(BeNumerically("~",
			time.Now().Add(-2*(50*time.Millisecond)).Truncate(50*time.Millisecond).Unix(),
			1,
		))
		Eventually(graphiteClient._sendMetricsCount).Should(BeNumerically(">", 1))
	})
})

type spyPointBuilder struct {
	mu                    sync.Mutex
	_buildCalled          int
	_buildPointsTimestamp int64
}

func (s *spyPointBuilder) BuildPoints(timestamp int64) ([]graphite.Metric, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s._buildCalled++
	s._buildPointsTimestamp = timestamp

	return []graphite.Metric{
		{
			Name:      "application.ingress",
			Value:     "1234",
			Timestamp: 1257894000,
		},
		{
			Name:      "application.ingress",
			Value:     "1234",
			Timestamp: 1257894000,
		},
	}, nil
}

func (s *spyPointBuilder) buildCalled() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s._buildCalled
}

func (s *spyPointBuilder) buildPointsTimestamp() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s._buildPointsTimestamp
}

type spyReadCloser struct{}

func (s *spyReadCloser) Close() error {
	return nil
}

func (s *spyReadCloser) Read([]byte) (int, error) {
	return 0, nil
}

type spyGraphiteClient struct {
	mu                sync.Mutex
	_sendMetricsCount int
	_url              string
	_contentType      string
	_body             string
}

func (s *spyGraphiteClient) SendMetrics(points []graphite.Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s._sendMetricsCount++

	return nil
}

func (s *spyGraphiteClient) sendMetricsCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s._sendMetricsCount
}

func (s *spyGraphiteClient) Connect() error {
	return nil
}

func (s *spyGraphiteClient) Disconnect() error {
	return nil
}
