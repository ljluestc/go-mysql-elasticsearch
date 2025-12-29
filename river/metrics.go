package river

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	esInsertNum = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mysql2es_inserted_num",
			Help: "The number of docs inserted to elasticsearch",
		}, []string{"index"},
	)
	esUpdateNum = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mysql2es_updated_num",
			Help: "The number of docs updated to elasticsearch",
		}, []string{"index"},
	)
	esDeleteNum = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mysql2es_deleted_num",
			Help: "The number of docs deleted from elasticsearch",
		}, []string{"index"},
	)
	canalSyncState = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "mysql2es_canal_state",
			Help: "The canal slave running state: 0=stopped, 1=ok",
		},
	)
	canalDelay = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "mysql2es_canal_delay",
			Help: "The canal slave lag",
		},
	)
)

func (r *River) collectMetrics() {
	for range time.Tick(10 * time.Second) {
		canalDelay.Set(float64(r.canal.GetDelay()))
	}
}

func (r *River) RunHTTP(addr string, path string) {
	http.Handle(path, promhttp.Handler())
	http.HandleFunc("/reindex", r.handleReindex)
	http.ListenAndServe(addr, nil)
}

func (r *River) handleReindex(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var body map[string]interface{}
	if err := decoder.Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := r.es.Reindex(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
