package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"k8s.io/klog"
)

var (
	errorRate = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ci_operator_error_rate",
			Help: "number of errors, sorted by label/type",
		},
		[]string{"error"},
	)

	port string
)

func init() {
	flag.StringVar(&port, "port", "8080", "Port the server listens to")

	prometheus.MustRegister(errorRate)
}

// Payload uploaded from a ci operator
type Payload struct {
	Reason  string `json:"reason"`
	JobName string `json:"job_name"`
	State   string `json:"state"`
	Type    string `json:"type"`
	Cluster string `json:"cluster"`
}

func validatePayload(payload *Payload) error {
	if payload.Reason == "" {
		return fmt.Errorf("reason field in payload is empty")
	}
	if payload.JobName == "" {
		return fmt.Errorf("job_name field in payload is empty")
	}
	if payload.State == "" {
		return fmt.Errorf("state field in payload is empty")
	}
	if payload.Type == "" {
		return fmt.Errorf("type field in payload is empty")
	}
	if payload.Cluster == "" {
		return fmt.Errorf("cluster field in payload is empty")
	}
	return nil
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, err)
}

func handleCIOperatorResult() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handleError(w, fmt.Errorf("unable to ready request body: %v", err))
			return
		}

		payload := &Payload{}
		if err := json.Unmarshal(bytes, payload); err != nil {
			handleError(w, fmt.Errorf("unable to decode request body: %v", err))
			return
		}

		if err := validatePayload(payload); err != nil {
			handleError(w, err)
			return
		}

		labels := prometheus.Labels{"error": payload.Reason}
		errorRate.With(labels).Inc()

		w.WriteHeader(http.StatusOK)

		klog.Infof("Request with %#v payload processed", payload)
	})
}

func main() {
	flag.Parse()

	http.HandleFunc("/write", handleCIOperatorResult())
	http.Handle("/metrics", promhttp.Handler())

	klog.Infof("Server configured to list on port %v", port)

	klog.Fatal(http.ListenAndServe(":"+port, nil))
}
