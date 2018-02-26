package builder

import (
	"fmt"
	"log"
	"strings"

	nn_collector "code.cloudfoundry.org/noisy-neighbor-nozzle/pkg/collector"
	nn_store "code.cloudfoundry.org/noisy-neighbor-nozzle/pkg/store"

	graphite "github.com/marpaia/graphite-golang"
)

type Builder interface {
	BuildPoints(timestamp int64) ([]interface{}, error)
}

// Fetcher provides a way of gathering the rates from the nozzles
type Fetcher interface {
	Rate(timestamp int64) (nn_store.Rate, error)
}

// GraphiteCollector handles fetch rates form multiple nozzles and summing their
// rates.
type GraphiteBuilder struct {
	fetcher       Fetcher
	store         nn_collector.AppInfoStore
	metricsPrefix string
}

// New initializes and returns a new GraphiteCollector.
func NewGraphiteBuilder(
	fetcher Fetcher,
	store nn_collector.AppInfoStore,
	metricsPrefix string,
) *GraphiteBuilder {

	gp := &GraphiteBuilder{
		fetcher:       fetcher,
		store:         store,
		metricsPrefix: metricsPrefix,
	}

	return gp
}

// BuildPoints satisfies the graphite Builder interface. It will
// request all the rates from all the known nozzles and sum their counts.
func (gp *GraphiteBuilder) BuildPoints(timestamp int64) ([]graphite.Metric, error) {
	rate, err := gp.fetcher.Rate(timestamp)
	if err != nil {
		return nil, err
	}

	var top counts

	for k, v := range rate.Counts {
		top = append(top, count{
			guidIndex: k,
			value:     v,
		})
	}

	var guids []string
	for _, c := range top {
		g := GUIDIndex(c.guidIndex).GUID()
		guids = append(guids, g)
	}
	// The underlying cached store does not return an error and instead simply
	// returns the cache when an error occurs.
	appInfo, err := gp.store.Lookup(guids)
	if err != nil {
		log.Printf("%s: failed to collect app metadata from API lookup", err)
	}

	var graphitePoints []graphite.Metric
	for _, c := range top {
		gi := GUIDIndex(c.guidIndex)

		orgSpaceAppName, ok := appInfo[nn_collector.AppGUID(gi.GUID())]

		if ok && checkOrgSpaceAppNameIsNotEmpty(orgSpaceAppName) {

			metricName := fmt.Sprintf("%s.%s.%s", gp.metricsPrefix, orgSpaceAppName, gi.Index())
			graphitePoints = append(graphitePoints, graphite.Metric{
				Name:      metricName,
				Value:     fmt.Sprintf("%d", c.value),
				Timestamp: rate.Timestamp,
			})

		} else {

			log.Printf("%s: failed to extract metric metadata from API lookup", c)
		}
	}

	return graphitePoints, nil
}

func checkOrgSpaceAppNameIsNotEmpty(orgSpaceAppName nn_collector.AppInfo) bool {

	if orgSpaceAppName.Name != "" && orgSpaceAppName.Space != "" && orgSpaceAppName.Org != "" {
		return true
	}

	return false
}

type count struct {
	guidIndex string
	value     uint64
}

type counts []count

func (c counts) Len() int           { return len(c) }
func (c counts) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c counts) Less(i, j int) bool { return c[i].value > c[j].value }

// GUIDIndex is a concatentation of GUID and instance index in the format
// some-guid/some-index, e.g., 7b8228a0-cf40-42d8-a7bb-b287a88198a3/0
type GUIDIndex string

// GUID returns the GUID of the GUIDIndex
func (g GUIDIndex) GUID() string {
	return strings.Split(string(g), "/")[0]
}

// Index returns the Index of the GUIDIndex
func (g GUIDIndex) Index() string {
	parts := strings.Split(string(g), "/")
	if len(parts) < 2 {
		return "0"
	}
	return parts[1]
}
