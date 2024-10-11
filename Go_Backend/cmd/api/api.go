package main 

import (
  "log"
  "net/http"
  "time"
  
  "CoffeeMap/internal/store"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
)

type application struct {
  config  config
  store   store.Storage
  db      dbConfig
}

type config struct {
  addr     string
  db       dbConfig
  env      string
  apiURL   string
  mail     mailConfig
}

type mailConfig struct {
  exp time.Duration
}

type dbConfig struct {
  addr          string
  maxOpenConns  int
  maxIdleConns  int
  maxIdleTime   string
}

func (app *application) mount() http.Handler {
  r := chi.NewRouter()

  r.Use(middleware.RequestID)
  r.Use(middleware.RealIP)
  r.Use(middleware.Logger)
  r.Use(middleware.Recoverer)


  r.Route("/v1", func(r chi.Router) {
    r.Get("/health", app.healthCheckHandler)

  // comment
    r.Route("/comments", func(r chi.Router) {
      r.Post("/", app.createCommentHandler)
      r.Route("/{commentID}", func (r chi.Router) {
        r.Use(app.commentsContextMiddleware)

        r.Get("/", app.getCommentHandler) 
        r.Delete("/", app.deleteCommentHandler)
        r.Patch("/", app.updateCommentHandler)
      })
    })
  // users
   r.Route("/users", func(r chi.Router) {
     r.Route("/{userID}", func(r chi.Router) {
       r.Get("/", app.getUserHandler)
     }) 
   }) 
   r.Route("/shops", func (r chi.Router) {
     r.Route("/{shopID}", func (r chi.Router) {
       r.Use(app.shopsContextMiddleware)
       r.Get("/", app.getShopHandler) 
     })
  }) 
  // public routes
  r.Route("/authentication", func (r chi.Router) {
    r.Post("/user", app.registerUserHandler)
  })
})
  // auth

  return r
}

func (app *application) run(mux http.Handler) error {

  srv := &http.Server{
    Addr: app.config.addr,
    Handler: mux,
    WriteTimeout: time.Second * 30,
    ReadTimeout: time.Second * 10,
    IdleTimeout: time.Minute,
  }

  log.Printf("server running on port %s", app.config.addr)
  return srv.ListenAndServe();
}
