package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/pz3-http/internal/storage"
)

func TestCreateTaskTitleLength(t *testing.T) {
	tests := []struct {
		name       string
		title      string
		wantStatus int
	}{
		{name: "one character", title: "a", wantStatus: http.StatusCreated},
		{name: "144 characters", title: strings.Repeat("a", 144), wantStatus: http.StatusCreated},
		{name: "145 characters", title: strings.Repeat("a", 145), wantStatus: http.StatusBadRequest},
		{name: "144 unicode characters", title: strings.Repeat("я", 144), wantStatus: http.StatusCreated},
		{name: "145 unicode characters", title: strings.Repeat("я", 145), wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(createTaskRequest{Title: tt.title})
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			NewHandlers(storage.NewMemoryStore()).CreateTask(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestPatchTaskDone(t *testing.T) {
	store := storage.NewMemoryStore()
	task := store.Create("test")
	h := NewHandlers(store)

	body, _ := json.Marshal(patchTaskRequest{Done: true})
	req := httptest.NewRequest(http.MethodPatch, "/tasks/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.PatchTask(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got storage.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if !got.Done {
		t.Fatalf("done = false, want true")
	}
	if got.ID != task.ID {
		t.Fatalf("id = %d, want %d", got.ID, task.ID)
	}
}

func TestPatchTaskNotFound(t *testing.T) {
	h := NewHandlers(storage.NewMemoryStore())

	body, _ := json.Marshal(patchTaskRequest{Done: true})
	req := httptest.NewRequest(http.MethodPatch, "/tasks/999", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.PatchTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestDeleteTask(t *testing.T) {
	store := storage.NewMemoryStore()
	store.Create("test")
	h := NewHandlers(store)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	rec := httptest.NewRecorder()

	h.DeleteTask(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}

	if _, err := store.Get(1); err == nil {
		t.Fatal("task still exists after delete")
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	h := NewHandlers(storage.NewMemoryStore())

	req := httptest.NewRequest(http.MethodDelete, "/tasks/999", nil)
	rec := httptest.NewRecorder()

	h.DeleteTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestGetTask(t *testing.T) {
	store := storage.NewMemoryStore()
	task := store.Create("test")
	h := NewHandlers(store)

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	rec := httptest.NewRecorder()

	h.GetTask(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got storage.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Title != task.Title {
		t.Fatalf("title = %q, want %q", got.Title, task.Title)
	}
}
