package handlers

import (
  "github.com/go-chi/chi/v5"
  "github.com/jackc/pgx/v5/pgxpool"
  "net/http"
  "scrapper/cmd/api/handlers/postgrp"
  "scrapper/internal/repository"
  "scrapper/internal/service"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Cors() Middleware {
  m := func(handler http.HandlerFunc) http.HandlerFunc {
    h := func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Access-Control-Allow-Origin", "*")
      w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
      w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

      handler(w, r)
    }

    return h
  }

  return m
}

func API(db *pgxpool.Pool) http.Handler {
  mux := chi.NewMux()

  repo := repository.NewPostRepository(db)
  serv := service.NewPostService(repo)
  postHandler := postgrp.NewHandler(serv)

  cors := Cors()
  mux.Get("/post/query/{prompt}", cors(postHandler.Get))
  mux.Get("/post/{id}", cors(postHandler.GetOne))
  return mux
}