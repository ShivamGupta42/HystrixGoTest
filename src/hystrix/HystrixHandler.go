package hystrix

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gopkg.in/alexcesaro/statsd.v2"
	"net"
	"net/http"
)

func init() {
	NewHystrix()
}

/*
  type MetricResult struct {
	Attempts                float64
	Errors                  float64
	Successes               float64
	Failures                float64
	Rejects                 float64
	ShortCircuits           float64
	Timeouts                float64
	FallbackSuccesses       float64
	FallbackFailures        float64
	ContextCanceled         float64
	ContextDeadlineExceeded float64
	TotalDuration           time.Duration
	RunDuration             time.Duration
	ConcurrencyInUse        float64
}



*/

var (

	//Attempts = promauto.NewCounter(prometheus.CounterOpts{
	//Name: "hystrix_ops_concurrency",
	//Help: "The total number of hystrix requests",
	//})
	//
	//Errors = promauto.NewCounter(prometheus.CounterOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})
	//
	//Successes = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})
	//
	//Failures = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})
	//
	//Rejects = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})
	//
	//ShortCircuits = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})
	//
	//Timeouts = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})
	//
	//FallbackSuccesses = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})
	//
	//FallbackFailures = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})
	//
	//ContextCanceled = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})
	//
	//ContextDeadlineExceeded = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})
	//
	//TotalDuration = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})
	//
	//RunDuration = promauto.NewGauge(prometheus.GaugeOpts{
	//	Name: "hystrix_ops_concurrency",
	//	Help: "The total number of hystrix requests",
	//})

	ConcurrencyInUse = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "hystrix_ops_concurrency",
		Help: "The total number of hystrix concurrency",
	})
)

type PromHysCounter struct {
	curMetric metricCollector.MetricResult
}

func (p PromHysCounter) Update(m metricCollector.MetricResult) {
	ConcurrencyInUse.Set(m.ConcurrencyInUse)
}

func (p PromHysCounter) Reset() {

}

func NewHystrix() {
	hystrix.ConfigureCommand("updateCounter", hystrix.CommandConfig{
		Timeout:                1000,
		MaxConcurrentRequests:  10,
		RequestVolumeThreshold: 10,
		SleepWindow:            1000,
		ErrorPercentThreshold:  30,
	})

	//Hystrix dashboard
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("localhost", "9000"), hystrixStreamHandler)

	//Starting StatsD Client
	statsDConfig := StatsDConfig{
		Host:      "localhost",
		Port:      8125,
		Enabled:   true,
		Namespace: "hystrixTest",
		Tags:      nil,
	}

	//StatsD Collector
	NewHystrixStatsDCollector(statsDConfig)

	//Register Prom Metric Collector
	metricCollector.Registry.Register(initPromCollector)
}

func initPromCollector(name string) metricCollector.MetricCollector {
	return PromHysCounter{}
}

type StatsDConfig struct {
	Host      string
	Port      int
	Enabled   bool
	Namespace string
	Tags      []string
}

type StatsD struct {
	client *statsd.Client
}

func (cfg StatsDConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

func NewHystrixStatsDCollector(statsDConfig StatsDConfig) {
	if !statsDConfig.Enabled {
		return
	}

	address := fmt.Sprintf("%s:%d", statsDConfig.Host, statsDConfig.Port)
	c, err := plugins.InitializeStatsdCollector(&plugins.StatsdCollectorConfig{
		StatsdAddr: address,
		Prefix:     statsDConfig.Namespace + ".hystrix",
	})
	if err != nil {
		panic(errors.New("failed to initiate StatsD for Hystrix"))
	}

	metricCollector.Registry.Register(c.NewStatsdCollector)
}

func NewStatsD(cfg StatsDConfig) (*StatsD, error) {
	if !cfg.Enabled {
		return &StatsD{}, nil
	}

	client, err := statsd.New(statsd.Address(cfg.Address()), statsd.Prefix(cfg.Namespace))
	if err != nil {
		panic(errors.New("failed to initiate StatsD"))
	}
	return &StatsD{client: client}, nil
}
