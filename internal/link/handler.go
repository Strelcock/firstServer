package link

import (
	"firstServer/configs"
	"firstServer/internal/di"
	"firstServer/pkg/event"
	"firstServer/pkg/middleware"
	"firstServer/pkg/req"
	"firstServer/pkg/res"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandler struct {
	LinkRepo *LinkRepository
	StatRepo di.IStatRepository
	EventBus *event.EventBus
}

type LinkHandlerDeps struct {
	LinkRepo *LinkRepository
	Config   *configs.Config
	EventBus *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := LinkHandler{
		LinkRepo: deps.LinkRepo,
		EventBus: deps.EventBus,
	}
	router.Handle("POST /link", middleware.IsAuthed(handler.Create(), deps.Config))
	router.HandleFunc("GET /{hash}", handler.GoTo())
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update(), deps.Config))
	router.Handle("DELETE /link/{id}", middleware.IsAuthed(handler.Delete(), deps.Config))
	router.Handle("GET /link", middleware.IsAuthed(handler.GetAll(), deps.Config))

}

func (h *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LinkCreateRequset](&w, r)
		if err != nil {
			return
		}

		link := NewLink(body.Url)
		for {
			if res, _ := h.LinkRepo.GetByHash(link.Hash); res != nil {
				link.RegenerateHash()
			} else {
				break
			}
		}
		createdLink, err := h.LinkRepo.Create(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createdLink, 201)

	}
}

func (h *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		link, err := h.LinkRepo.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		go h.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited,
			Data: link.ID,
		})
		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}

func (h *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxEmail := r.Context().Value(middleware.ContextEmailKey)
		fmt.Println(ctxEmail)

		body, err := req.HandleBody[LinkUpdateRequest](&w, r)
		if err != nil {
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		link, err := h.LinkRepo.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, link, 201)

	}
}

func (h *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		_, err = h.LinkRepo.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = h.LinkRepo.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, nil, 200)
	}
}

func (h *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}
		links := h.LinkRepo.GetAll(limit, offset)
		count := h.LinkRepo.Count()
		res.Json(w, GetAllLinksResponse{
			Links: links,
			Count: count,
		}, 200)

	}
}
