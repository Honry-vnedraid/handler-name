package main

import (
	"encoding/json"
	"fmt"
	"handler-service/data"
	"net/http"
	"strconv"
)

type Monitor struct {
	handler *Handler
}

func (monitor *Monitor) listenAndServe(addr string) error {
	fmt.Printf("http://%s", addr)
	return http.ListenAndServe(addr, nil)
}

func (monitor *Monitor) initHandling() {
	http.Handle("/add/news", monitor.AddNews())
	http.Handle("/news", monitor.GetNews())
	http.Handle("/summary", monitor.Summary())
}

func (monitor *Monitor) AddNews() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

		var news data.News
		err := json.NewDecoder(r.Body).Decode(&news)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		go monitor.handler.addNews(&news)

		w.Write([]byte("ok"))
	})
}

func (monitor *Monitor) GetNews() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		offset, err := strconv.Atoi(r.FormValue("offset"))
		if err != nil {
			http.Error(w, "offset variable should be a number", http.StatusUnprocessableEntity)
			return
		}

		limit, err := strconv.Atoi(r.FormValue("limit"))
		if err != nil {
			http.Error(w, "limit variable should be a number", http.StatusUnprocessableEntity)
			return
		}

		news, err := monitor.handler.getNews(offset, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(news)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	})
}

func (monitor *Monitor) Summary() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
	})
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
