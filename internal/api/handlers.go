package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"example.com/pz3-http/internal/storage"
)

type Handlers struct {
	Store *storage.MemoryStore
}

func NewHandlers(store *storage.MemoryStore) *Handlers {
	return &Handlers{Store: store}
}

// GET /tasks
func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.Store.List()

	// Поддержка простых фильтров через query: ?q=text
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q != "" {
		filtered := tasks[:0]
		for _, t := range tasks {
			if strings.Contains(strings.ToLower(t.Title), strings.ToLower(q)) {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}

	JSON(w, http.StatusOK, tasks)
}

type createTaskRequest struct {
	Title string `json:"title"`
}

// POST /tasks
func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "" && !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		BadRequest(w, "Content-Type must be application/json")
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json: "+err.Error())
		return
	}
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		BadRequest(w, "title is required")
		return
	}
	if utf8.RuneCountInString(req.Title) > 144 {
		BadRequest(w, "title must be between 1 and 144 characters")
		return
	}

	t := h.Store.Create(req.Title)
	JSON(w, http.StatusCreated, t)
}

func parseTaskID(r *http.Request) (int64, bool, string) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 {
		return 0, false, "invalid path"
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, false, "invalid id"
	}
	return id, true, ""
}

// GET /tasks/{id}
func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id, ok, msg := parseTaskID(r)
	if !ok {
		if msg == "invalid id" {
			BadRequest(w, msg)
		} else {
			NotFound(w, msg)
		}
		return
	}

	t, err := h.Store.Get(id)
	if err != nil {
		NotFound(w, "task not found")
		return
	}
	JSON(w, http.StatusOK, t)
}

type patchTaskRequest struct {
	Done bool `json:"done"`
}

// PATCH /tasks/{id}
func (h *Handlers) PatchTask(w http.ResponseWriter, r *http.Request) {
	id, ok, msg := parseTaskID(r)
	if !ok {
		if msg == "invalid id" {
			BadRequest(w, msg)
		} else {
			NotFound(w, msg)
		}
		return
	}

	var req patchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		BadRequest(w, "invalid json: "+err.Error())
		return
	}

	t, err := h.Store.SetDone(id, req.Done)
	if err != nil {
		NotFound(w, "task not found")
		return
	}
	JSON(w, http.StatusOK, t)
}

// DELETE /tasks/{id}
func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok, msg := parseTaskID(r)
	if !ok {
		if msg == "invalid id" {
			BadRequest(w, msg)
		} else {
			NotFound(w, msg)
		}
		return
	}

	if err := h.Store.Delete(id); err != nil {
		NotFound(w, "task not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
