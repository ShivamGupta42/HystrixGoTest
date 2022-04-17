package hystrixConfig

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/afex/hystrix-go/plugins"
	"github.com/prometheus/client_golang/prometheus"
	"strings"

	//"github.com/prometheus/client_golang/prometheus/promauto"
	"gopkg.in/alexcesaro/statsd.v2"
	"sync"
)

const (
	metricCircuitOpen       = "circuit_open"
	metricSuccesses         = "successes"
	metricAttempts          = "attempts"
	metricErrors            = "errors"
	metricFailures          = "failures"
	metricRejects           = "rejects"
	metricShortCircuits     = "short_circuits"
	metricTimeouts          = "timeouts"
	metricFallbackSuccesses = "fallback_successes"
	metricFallbackFailures  = "fallback_failures"
	metricTotalDuration     = "total_duration"
	metricRunDuration       = "run_duration"
	metricConcurrencyInUse  = "concurrency_in_use"
)

var (
	maxConcurrentRequests = 1.0
	gauges                = []string{metricCircuitOpen, metricTotalDuration, metricRunDuration}
	counters              = []string{metricSuccesses, metricAttempts, metricErrors, metricFailures,
		metricRejects, metricShortCircuits, metricTimeouts, metricFallbackSuccesses, metricFallbackFailures}
	histograms = []string{metricConcurrencyInUse}
)

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

//var (
//	Attempts = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_attempts",
//		Help: "The total number of requests",
//	})
//
//	Errors = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_errors",
//		Help: "The total number of errors",
//	})
//
//	Successes = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_successes",
//		Help: "The total number of successes",
//	})
//
//	Failures = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_failures",
//		Help: "The total number of failures",
//	})
//
//	Rejects = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_rejects",
//		Help: "The total number of rejects",
//	})
//
//	ShortCircuits = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_shortcircuits",
//		Help: "The total number of short-circuits",
//	})
//
//	Timeouts = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_timeouts",
//		Help: "The total number of timeouts",
//	})
//
//	FallbackSuccesses = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_fallback_success",
//		Help: "The total number of hystrixConfig requests",
//	})
//
//	FallbackFailures = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_fallback_failure",
//		Help: "The total number of hystrixConfig requests",
//	})
//
//	ContextCanceled = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_context_canceled",
//		Help: "The total number of hystrixConfig requests",
//	})
//
//	ContextDeadlineExceeded = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_context_deadline_exceeded",
//		Help: "The total number of hystrixConfig requests",
//	})
//
//	TotalDuration = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_total_duration",
//		Help: "The total number of hystrixConfig requests",
//	})
//
//	RunDuration = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_run_duration",
//		Help: "The total number of hystrixConfig requests",
//	})
//
//	MaxConcurrencyInUse = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_max_concurrency",
//		Help: "The total number of hystrixConfig concurrency",
//	})
//	MinConcurrencyInUse = promauto.NewGauge(prometheus.GaugeOpts{
//		Name: "hystrix_ops_min_concurrency",
//		Help: "The total number of hystrixConfig concurrency",
//	})
//)
//
//var maxConcurrency float64
//var minConcurrency float64
//var mu sync.Mutex

type PrometheusCollector struct {
	sync.RWMutex
	namespace  string
	subsystem  string
	gauges     map[string]prometheus.Gauge
	counters   map[string]prometheus.Counter
	summaries  map[string]prometheus.Summary
	histograms map[string]prometheus.Histogram
}

func (c *PrometheusCollector) Update(r metricCollector.MetricResult) {

	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	// check circuit open
	if r.Successes > 0 {
		gauge := c.gauges[metricCircuitOpen]
		gauge.Set(0)

		counter := c.counters[metricSuccesses]
		counter.Add(r.Successes)
	}
	if r.ShortCircuits > 0 {
		gauge := c.gauges[metricCircuitOpen]
		gauge.Set(1)

		counter := c.counters[metricShortCircuits]
		counter.Add(r.ShortCircuits)
	}
	// update  metric
	if r.Attempts > 0 {
		counter := c.counters[metricAttempts]
		counter.Add(r.Attempts)
	}
	if r.Errors > 0 {
		counter := c.counters[metricErrors]
		counter.Add(r.Errors)
	}
	if r.Failures > 0 {
		counter := c.counters[metricFailures]
		counter.Add(r.Failures)
	}
	if r.Rejects > 0 {
		counter := c.counters[metricRejects]
		counter.Add(r.Rejects)
	}
	if r.Timeouts > 0 {
		counter := c.counters[metricTimeouts]
		counter.Add(r.Timeouts)
	}
	if r.FallbackSuccesses > 0 {
		counter := c.counters[metricFallbackSuccesses]
		counter.Add(r.FallbackSuccesses)
	}
	if r.FallbackFailures > 0 {
		counter := c.counters[metricFallbackFailures]
		counter.Add(r.FallbackFailures)
	}

	gauge := c.gauges[metricTotalDuration]
	gauge.Set(r.TotalDuration.Seconds())

	gauge = c.gauges[metricRunDuration]
	gauge.Set(r.RunDuration.Seconds())

	histogram := c.histograms[metricConcurrencyInUse]
	histogram.Observe(maxConcurrentRequests * r.ConcurrencyInUse)
}

func (c *PrometheusCollector) Reset() {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
}

func NewHystrix() {
	hystrix.ConfigureCommand("updateCounter", hystrix.CommandConfig{
		Timeout:                1000,
		MaxConcurrentRequests:  int(maxConcurrentRequests),
		RequestVolumeThreshold: 10,
		SleepWindow:            1000,
		ErrorPercentThreshold:  30,
	})

	//Hystrix dashboard
	//hystrixStreamHandler := hystrix.NewStreamHandler()
	//hystrixStreamHandler.Start()
	//go http.ListenAndServe(net.JoinHostPort("localhost", "9000"), hystrixStreamHandler)

	////Starting StatsD Client
	//statsDConfig := StatsDConfig{
	//	Host:      "localhost",
	//	Port:      8125,
	//	Enabled:   true,
	//	Namespace: "hystrixTest",
	//	Tags:      nil,
	//}

	//StatsD Collector
	//NewHystrixStatsDCollector(statsDConfig)

	//Register Prom Metric Collector
	metricCollector.Registry.Register(NewPrometheusCollector("default", map[string]string{"app": "shivamTestApp"}))
}

//func initPromCollector(name string) metricCollector.MetricCollector {
//	return PromHysCounter{}
//}

// NewPrometheusCollector returns wrapper function returning an implemented struct from MetricCollector.
func NewPrometheusCollector(namespace string, labels map[string]string) func(string) metricCollector.MetricCollector {
	return func(name string) metricCollector.MetricCollector {
		name = strings.Replace(name, "/", "_", -1)
		name = strings.Replace(name, ":", "_", -1)
		name = strings.Replace(name, ".", "_", -1)
		name = strings.Replace(name, "-", "_", -1)
		collector := &PrometheusCollector{
			namespace:  namespace,
			subsystem:  name,
			gauges:     map[string]prometheus.Gauge{},
			counters:   map[string]prometheus.Counter{},
			histograms: map[string]prometheus.Histogram{},
			summaries:  map[string]prometheus.Summary{},
		}

		// make gauges
		for _, metric := range gauges {
			opts := prometheus.GaugeOpts{
				Namespace: collector.namespace,
				Subsystem: collector.subsystem,
				Name:      metric,
				Help:      fmt.Sprintf("[gauge] namespace : %s, metric : %s", collector.namespace, metric),
			}
			if labels != nil {
				opts.ConstLabels = labels
			}
			gauge := prometheus.NewGauge(opts)
			collector.gauges[metric] = gauge
			prometheus.Register(gauge)
		}
		// make counters
		for _, metric := range counters {
			opts := prometheus.CounterOpts{
				Namespace: collector.namespace,
				Subsystem: collector.subsystem,
				Name:      metric,
				Help:      fmt.Sprintf("[counter] namespace : %s, metric : %s", collector.namespace, metric),
			}
			if labels != nil {
				opts.ConstLabels = labels
			}
			counter := prometheus.NewCounter(opts)
			collector.counters[metric] = counter
			prometheus.Register(counter)
		}

		for _, metric := range histograms {
			opts := prometheus.HistogramOpts{
				Namespace:   collector.namespace,
				Subsystem:   collector.subsystem,
				Name:        metric,
				Help:        "",
				ConstLabels: nil,
				Buckets:     generateHistBuckets(),
			}

			if labels != nil {
				opts.ConstLabels = labels
			}
			histogram := prometheus.NewHistogram(opts)
			collector.histograms[metric] = histogram
			prometheus.Register(histogram)

		}
		return collector
	}
}

func generateHistBuckets() []float64 {
	bucketSize := maxConcurrentRequests / 10
	var buckets []float64
	for i := 1; i <= 10; i++ {
		buckets = append(buckets, float64(i)*bucketSize)
	}
	return buckets
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
		Prefix:     statsDConfig.Namespace + ".hystrixConfig",
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
