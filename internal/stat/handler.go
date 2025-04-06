package stat

import (
	"firstServer/configs"
	"firstServer/pkg/middleware"
	"firstServer/pkg/res"
	"net/http"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandler struct {
	StatRepo *StatRepository
}

type StatHandlerDeps struct {
	StatRepo *StatRepository
	Config   *configs.Config
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepo: deps.StatRepo,
	}

	router.Handle("GET /stat", middleware.IsAuthed(handler.GetStat(), deps.Config))
}

func (sh *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const layout = "2006-01-02"
		from, err := time.Parse(layout, r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		to, err := time.Parse(layout, r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		by := r.URL.Query().Get("by")

		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid by param", http.StatusBadRequest)
			return
		}

		stats := sh.StatRepo.GetStats(by, from, to)
		res.Json(w, stats, http.StatusOK)
	}
}
