package app

import (
	"crypto/tls"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	uaaAddr              = kingpin.Flag("uaa-addr", "UAA address").Envar("UAA_ADDR").Required().String()
	capiAddr             = kingpin.Flag("capi-addr", "Api endpoint address.").Envar("CAPI_ADDR").Required().String()
	accumulatorAddr      = kingpin.Flag("accumulator-addr", "Api endpoint address.").Envar("ACCUMULATOR_ADDR").Required().String()
	syslogServer         = kingpin.Flag("syslog-server", "Syslog server.").Envar("SYSLOG_ENDPOINT").String()
	clientID             = kingpin.Flag("client-id", "Client ID.").Envar("CLIENT_ID").Required().String()
	clientSecret         = kingpin.Flag("client-secret", "Client secret.").Envar("CLIENT_SECRET").Required().String()
	metricsHost          = kingpin.Flag("metrics-host", "Metrics Host.").Envar("METRICS_HOST").Required().String()
	metricsPort          = kingpin.Flag("metrics-port", "Metrics Port.").Envar("METRICS_PORT").Required().Int()
	graphitePrefix       = kingpin.Flag("graphite-prefix", "Graphite metrics prefix").Envar("GRAPHITE_PREFIX").Required().String()
	skipCertVerify       = kingpin.Flag("skip-cert-verify", "Please don't").Default("false").Envar("SKIP_CERT_VERIFY").Bool()
	reportInterval       = kingpin.Flag("report-interval", "Report interval").Default("1m").Envar("REPORT_INTERVAL").Duration()
	reportLimit          = kingpin.Flag("report-limit", "Report limit").Default("50").Envar("REPORT_LIMIT").Int()
	appInfoCacheDuration = kingpin.Flag("cache-duration", "APP INFO CACHE DURATION").Default("150s").Envar("APP_INFO_CACHE_TTL").Duration()
)

// Config stores configuration data for the accumulator.
type Config struct {
	UAAAddr         string
	CAPIAddr        string
	AccumulatorAddr string
	ClientID        string
	ClientSecret    string
	GraphiteHost    string
	GraphitePort    int
	GraphitePrefix  string
	SkipCertVerify  bool
	ReportInterval  time.Duration
	ReportLimit     int

	AppInfoCacheTTL time.Duration

	TLSConfig *tls.Config
}

// LoadConfig loads the configuration settings from the current environment.
func LoadConfig() Config {

	kingpin.Parse()

	cfg := Config{
		UAAAddr:         *uaaAddr,
		CAPIAddr:        *capiAddr,
		AccumulatorAddr: *accumulatorAddr,
		ClientID:        *clientID,
		ClientSecret:    *clientSecret,
		GraphiteHost:    *metricsHost,
		GraphitePort:    *metricsPort,
		GraphitePrefix:  *graphitePrefix,
		SkipCertVerify:  *skipCertVerify,
		ReportInterval:  *reportInterval,
		ReportLimit:     *reportLimit,

		AppInfoCacheTTL: *appInfoCacheDuration,
	}

	cfg.TLSConfig = &tls.Config{InsecureSkipVerify: cfg.SkipCertVerify}

	return cfg
}
