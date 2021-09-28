package handlers

import (
	"bytes"
	"context"
	"embed"
	"github.com/takeoff-projects/level-1/business/data/event"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	Templates embed.FS
	Log       *zap.Logger
	Store     event.Store
}

func (h Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	events, err := h.Store.GetEvents(ctx)
	if err != nil {
		h.Log.Error("failed to fetch events", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := HomePageData{
		PageTitle: "Home Page",
		Events:    events,
		Count:     len(events),
	}

	var tpl = template.Must(template.ParseFS(h.Templates, "templates/index.html", "templates/layout.html"))

	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		h.Log.Error("failed to execute template", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
	h.Log.Info("Home Page Served")
}

func (h Handler) AboutHandler(w http.ResponseWriter, r *http.Request) {
	data := AboutPageData{
		PageTitle: "About Go Website",
	}

	var tpl = template.Must(template.ParseFS(h.Templates, "templates/about.html", "templates/layout.html"))

	buf := &bytes.Buffer{}
	err := tpl.Execute(buf, data)
	if err != nil {
		h.Log.Error("failed to execute template", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf.WriteTo(w)
	h.Log.Info("About Page Served")
}

func (h Handler) AddHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method == http.MethodGet {
		data := AddPageData{
			PageTitle: "Add Event",
		}

		var tpl = template.Must(template.ParseFS(h.Templates, "templates/add.html", "templates/layout.html"))

		buf := &bytes.Buffer{}
		err := tpl.Execute(buf, data)
		if err != nil {
			h.Log.Error("failed to execute template", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		buf.WriteTo(w)

		h.Log.Info("Add Page Served")
	} else {
		// Add Event Here
		event := event.Event{
			Title:     r.FormValue("title"),
			Location:  r.FormValue("location"),
			EventDate: r.FormValue("when"),
		}
		if err := h.Store.AddEvent(ctx, event); err != nil {
			h.Log.Error("failed to create event", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Go back to home page
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h Handler) EditHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if r.Method == http.MethodGet {
		event, err := h.Store.GetEventByID(ctx, mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.Log.Error("failed to edit event", zap.Error(err))
			return
		}

		data := EditPageData{
			PageTitle: "Edit Event",
			Event:     event,
		}

		var tpl = template.Must(template.ParseFS(h.Templates, "templates/edit.html", "templates/layout.html"))

		buf := &bytes.Buffer{}
		if err := tpl.Execute(buf, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.Log.Error("failed to execute template", zap.Error(err))
			return
		}
		buf.WriteTo(w)

		h.Log.Info("Edit Page Served")
	} else {
		// Add Event Here
		event := event.Event{
			ID:        r.FormValue("id"),
			Title:     r.FormValue("title"),
			Location:  r.FormValue("location"),
			EventDate: r.FormValue("when"),
		}
		if err := h.Store.UpdateEvent(ctx, event); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.Log.Error("failed to update event", zap.Error(err))
			return
		}
		h.Log.Info("Event Updated")

		// Go back to home page
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (h Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if err := h.Store.DeleteEvent(ctx, mux.Vars(r)["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Log.Error("failed to update event", zap.Error(err))
		return
	}
	h.Log.Info("Event deleted")

	// Go back to home page
	http.Redirect(w, r, "/", http.StatusFound)
}

// HomePageData for Index template
type HomePageData struct {
	PageTitle string
	Events    []event.Event
	Count     int
}

// AboutPageData for About template
type AboutPageData struct {
	PageTitle string
}

// AddPageData for About template
type AddPageData struct {
	PageTitle string
}

// EditPageData for About template
type EditPageData struct {
	PageTitle string
	Event     event.Event
}
