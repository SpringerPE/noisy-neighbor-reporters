package app

import (
	"log"
	"net/http"
	//	"strconv"
	"time"

	"code.cloudfoundry.org/noisy-neighbor-nozzle/pkg/auth"
	"code.cloudfoundry.org/noisy-neighbor-nozzle/pkg/collector"

	nn_collector "code.cloudfoundry.org/noisy-neighbor-nozzle/pkg/collector"

	"github.com/SpringerPE/noisy-neighbor-reporters/pkg/builder"
	graphite_builder "github.com/SpringerPE/noisy-neighbor-reporters/pkg/builder/graphite"

	"github.com/SpringerPE/noisy-neighbor-reporters/pkg/reporter"

	graphite "github.com/marpaia/graphite-golang"
)

// Reporter is the constructor for the datadog reporter application.
type Reporter struct {
	reporter *reporter.GraphiteReporter
}

// NewReporter configures and returns a new Reporter
func NewReporter(cfg Config) *reporter.GraphiteReporter {

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: cfg.TLSConfig,
		},
	}

	a := auth.NewAuthenticator(cfg.ClientID, cfg.ClientSecret, cfg.UAAAddr,
		auth.WithHTTPClient(client),
	)

	httpStore := builder.NewCFLightApiAppInfoStore(cfg.CAPIAddr, client)
	cache := collector.NewCachedAppInfoStore(
		httpStore,
		collector.WithCacheTTL(cfg.AppInfoCacheTTL),
	)

	log.Printf("initializing collector with accumulator: %v", cfg.AccumulatorAddr)
	c := nn_collector.New([]string{cfg.AccumulatorAddr}, a, "", cache,
		collector.WithReportLimit(cfg.ReportLimit),
		collector.WithHTTPClient(client),
	)

	b := graphite_builder.NewGraphiteBuilder(c, cache, cfg.GraphitePrefix)

	//graphitePort, err := strconv.Atoi(cfg.GraphitePort)
	/*	if err != nil {
			log.Fatalf("Please make sure that graphite port is a numeric value, %s", err)
		}
	*/
	graphiteClient, err := graphite.NewGraphite(cfg.GraphiteHost, cfg.GraphitePort)
	if err != nil {
		log.Fatalf("Error while connecting to graphite %s:%s: %s", cfg.GraphiteHost, cfg.GraphitePort, err)
	}

	log.Printf("initializing graphite reporter")

	r := reporter.NewReporter(b, graphiteClient,
		reporter.WithInterval(cfg.ReportInterval),
	)

	return r
}

// Run starts the graphite reporter. This is a blocking method call.
func (r *Reporter) Run() {
	r.reporter.Run()
}
