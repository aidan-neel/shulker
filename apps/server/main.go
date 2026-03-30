package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/aidan-neel/shulker/apps/proto/gen/go/auth/authconnect"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/blob/blobconnect"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/health/healthconnect"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/user/userconnect"

	"github.com/aidan-neel/shulker/apps/server/middleware"
	db "github.com/aidan-neel/shulker/apps/server/postgres/gen"
	"github.com/aidan-neel/shulker/apps/server/routes"
	"github.com/aidan-neel/shulker/apps/server/storage"
)

func main() {
	godotenv.Load()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	queries := db.New(stdlib.OpenDBFromPool(pool))
	fs := storage.NewLocalFilesystem("./blobs")
	enc := storage.NewAESEncryption()

	mux := http.NewServeMux()

	healthPath, healthHandler := healthconnect.NewHealthServiceHandler(routes.NewHealthHandler())
	mux.Handle(healthPath, middleware.Auth(healthHandler))

	usersPath, usersHandler := userconnect.NewUserServiceHandler(routes.NewUsersHandler(queries))
	mux.Handle(usersPath, middleware.Auth(usersHandler))

	blobsPath, blobsHandler := blobconnect.NewBlobServiceHandler(routes.NewBlobsHandler(queries, fs, enc))
	mux.Handle(blobsPath, middleware.Auth(blobsHandler))

	authPath, authHandler := authconnect.NewAuthServiceHandler(routes.NewAuthHandler(queries))
	mux.Handle(authPath, middleware.InjectResponseWriter(middleware.Auth(authHandler)))

	mux.Handle("/docs/specs/", http.StripPrefix("/docs/specs/", http.FileServer(http.Dir("../proto/gen/openapi"))))
	mux.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		spec := r.URL.Query().Get("spec")
		specs := map[string]string{
			"auth":   "/docs/specs/auth/auth.swagger.json",
			"blob":   "/docs/specs/blob/blob.swagger.json",
			"user":   "/docs/specs/user/user.swagger.json",
			"health": "/docs/specs/health/health.swagger.json",
		}
		url, ok := specs[spec]
		if !ok {
			url = specs["auth"]
			spec = "auth"
		}

		services := []string{"auth", "blob", "user", "health"}
		var navParts []string
		for _, s := range services {
			active := ""
			if s == spec {
				active = "background:#2d2d2d;color:#fff;"
			}
			navParts = append(navParts, `<a href="?spec=`+s+`" style="text-decoration:none;padding:6px 14px;border-radius:6px;color:#aaa;font-size:13px;`+active+`">`+s+`</a>`)
		}
		nav := strings.Join(navParts, "")

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
  <title>Shulker API</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width,initial-scale=1"/>
</head>
<body style="margin:0">
<div style="position:fixed;top:0;left:0;right:0;z-index:9999;display:flex;align-items:center;gap:4px;padding:8px 16px;background:#1a1a1a;border-bottom:1px solid #333;font-family:sans-serif">
  <span style="color:#666;font-size:13px;margin-right:8px">Shulker API</span>
  ` + nav + `
</div>
<div style="margin-top:49px">
  <script id="api-reference" data-url="` + url + `"></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</div>
</body>
</html>`))
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type", "Connect-Protocol-Version", "Connect-Timeout-Ms", "Grpc-Timeout", "Authorization"},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Grpc-Status", "Grpc-Message"},
	})

	log.Println("Running on localhost:8080")
	if err := http.ListenAndServe(":8080", c.Handler(h2c.NewHandler(mux, &http2.Server{}))); err != nil {
		log.Fatal(err)
	}
}
