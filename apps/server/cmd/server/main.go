package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/aidan-neel/shulker/apps/proto/gen/go/auth/authconnect"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/health/healthconnect"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/user/userconnect"
	auth "github.com/aidan-neel/shulker/apps/server/src/auth"
	health "github.com/aidan-neel/shulker/apps/server/src/health"
	"github.com/aidan-neel/shulker/apps/server/src/middleware"
	user "github.com/aidan-neel/shulker/apps/server/src/user"

	db "github.com/aidan-neel/shulker/apps/server/postgres/gen"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	sqlDB := stdlib.OpenDBFromPool(pool)
	queries := db.New(sqlDB)

	userRepo := user.NewUserRepo(queries)
	userService := user.NewService(userRepo)
	authService := auth.NewService(userRepo)

	mux := http.NewServeMux()

	healthPath, healthHandler := healthconnect.NewHealthServiceHandler(health.New())
	mux.Handle(healthPath, middleware.AuthMiddleware(healthHandler))

	userPath, userHandler := userconnect.NewUserServiceHandler(user.New(userService))
	mux.Handle(userPath, middleware.AuthMiddleware(userHandler))

	authPath, authHandler := authconnect.NewAuthServiceHandler(auth.New(authService))
	mux.Handle(authPath, authHandler)

	mux.Handle("/docs/swagger/", http.StripPrefix("/docs/swagger/", http.FileServer(http.Dir("../proto/gen/openapi"))))
	mux.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <title>Shulker API</title>
    <meta charset="utf-8"/>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
<script>
window.onload = function() {
    SwaggerUIBundle({
		urls: [
			{ url: "/docs/swagger/health/health.swagger.json", name: "Health" },
			{ url: "/docs/swagger/user/user.swagger.json", name: "User" },
			{ url: "/docs/swagger/auth/auth.swagger.json", name: "Auth" },
		],
        "urls.primaryName": "Health",
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
        layout: "StandaloneLayout"
    })
}
</script>
</body>
</html>`))
	})

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}
