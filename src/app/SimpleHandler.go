package app

import (
	"context"
	"encoding/json"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"sync/atomic"
)

var atomicCounter int64

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_ops_total",
		Help: "The total number of processed requests",
	})
)

func SimpleHandler(w http.ResponseWriter, r *http.Request) {
	hystrix.DoC(context.Background(), "updateCounter", func(ctx context.Context) error {
		updateCounter(w)
		return nil
	}, nil)
}

func updateCounter(w http.ResponseWriter) {
	v := atomic.AddInt64(&atomicCounter, 1)
	opsProcessed.Inc()
	if atomicCounter%1000 == 0 {
		hystrix.Flush()
	}

	m := map[string]int64{"counter": v}
	writeResponse(w, 200, m)
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
