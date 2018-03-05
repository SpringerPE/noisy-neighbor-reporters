package builder_test

import (
	nn_collector "code.cloudfoundry.org/noisy-neighbor-nozzle/pkg/collector"
	nn_store "code.cloudfoundry.org/noisy-neighbor-nozzle/pkg/store"

	graphite_builder "github.com/SpringerPE/noisy-neighbor-reporters/pkg/builder/graphite"
	graphite "github.com/marpaia/graphite-golang"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GraphiteBuilder", func() {

	It("can build a set of graphite points from some appInfo metrics", func() {
		fetcher := &fakeFetcher{}
		store := &fakeStore{path: "happyPath"}

		b := graphite_builder.NewGraphiteBuilder(fetcher, store, "test")
		points, err := b.BuildPoints(1520259517)

		Expect(err).ToNot(HaveOccurred())
		Expect(points).To(HaveLen(2))

		expectPointA := graphite.Metric{
			Name:      "test.org1.space1.app1.0",
			Value:     "2",
			Timestamp: 1520259517,
		}

		expectPointB := graphite.Metric{
			Name:      "test.org2.space2.app2.0",
			Value:     "3",
			Timestamp: 1520259517,
		}

		Expect(points).To(ContainElement(expectPointA))
		Expect(points).To(ContainElement(expectPointB))
	})

	It("it excludes metrics with missing fields", func() {
		fetcher := &fakeFetcher{}
		store := &fakeStore{path: "missingInfo"}

		b := graphite_builder.NewGraphiteBuilder(fetcher, store, "test")
		points, err := b.BuildPoints(1520259517)

		Expect(err).ToNot(HaveOccurred())
		Expect(points).To(HaveLen(0))
	})

	It("it excludes metrics which are not being cached", func() {
		fetcher := &fakeFetcher{}
		store := &fakeStore{path: "missingCacheInfo"}

		b := graphite_builder.NewGraphiteBuilder(fetcher, store, "test")
		points, err := b.BuildPoints(1520259517)

		Expect(err).ToNot(HaveOccurred())
		Expect(points).To(HaveLen(0))
	})
})

type fakeFetcher struct {
}

func (f *fakeFetcher) Rate(timestamp int64) (nn_store.Rate, error) {

	rate := nn_store.Rate{
		timestamp,
		map[string]uint64{
			"a": 2,
			"b": 3,
		},
	}

	return rate, nil
}

type fakeStore struct {
	path string
}

func (s *fakeStore) Lookup(guids []string) (map[nn_collector.AppGUID]nn_collector.AppInfo, error) {

	switch path := s.path; path {
	case "happyPath":
		return map[nn_collector.AppGUID]nn_collector.AppInfo{
			"a": nn_collector.AppInfo{
				Name:  "app1",
				Space: "space1",
				Org:   "org1",
			},
			"b": nn_collector.AppInfo{
				Name:  "app2",
				Space: "space2",
				Org:   "org2",
			},
		}, nil
	case "missingInfo":
		return map[nn_collector.AppGUID]nn_collector.AppInfo{
			"a": nn_collector.AppInfo{
				Name:  "",
				Space: "space1",
				Org:   "org1",
			},
			"b": nn_collector.AppInfo{
				Name:  "app2",
				Space: "",
				Org:   "org2",
			},
		}, nil
	case "missingCacheInfo":
		return map[nn_collector.AppGUID]nn_collector.AppInfo{
			"c": nn_collector.AppInfo{
				Name:  "c",
				Space: "space3",
				Org:   "org3",
			},
		}, nil
	default:
		return nil, nil
	}
}
