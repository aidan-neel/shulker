package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/aidan-neel/shulker/apps/proto/gen/go/auth/authconnect"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/blob/blobconnect"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/health/healthconnect"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/user/userconnect"
	"github.com/aidan-neel/shulker/apps/server/pkg/middleware"
	auth "github.com/aidan-neel/shulker/apps/server/src/auth"
	blob "github.com/aidan-neel/shulker/apps/server/src/blob"
	health "github.com/aidan-neel/shulker/apps/server/src/health"
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

	blobRepo := blob.NewBlobRepo(queries)
	userBlobRepo := blob.NewUserBlobRepo(queries)
	blobService := blob.NewService(blobRepo, userBlobRepo)

	mux := http.NewServeMux()

	healthPath, healthHandler := healthconnect.NewHealthServiceHandler(health.New())
	mux.Handle(healthPath, middleware.AuthMiddleware(healthHandler))

	userPath, userHandler := userconnect.NewUserServiceHandler(user.New(userService))
	mux.Handle(userPath, middleware.AuthMiddleware(userHandler))

	blobPath, blobHandler := blobconnect.NewBlobServiceHandler(blob.New(blobService))
	mux.Handle(blobPath, middleware.AuthMiddleware(blobHandler))

	authPath, authHandler := authconnect.NewAuthServiceHandler(auth.New(authService))
	mux.Handle(authPath, middleware.InjectResponseWriter(authHandler))

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
            { url: "/docs/swagger/blob/blob.swagger.json", name: "Blob" },
        ],
        "urls.primaryName": "Health",
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
        layout: "StandaloneLayout",
        persistAuthorization: true,
        securityDefinitions: {
            Bearer: {
                type: "apiKey",
                in: "header",
                name: "Authorization"
            }
        }
    })
}
</script>
</body>
</html>`))
	})

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Connect-Protocol-Version",
			"Connect-Timeout-Ms",
			"Grpc-Timeout",
			"Authorization",
		},
		AllowCredentials: true,
		ExposedHeaders: []string{
			"Grpc-Status",
			"Grpc-Message",
		},
	})

	log.Println("listening on :8080")
	if err := http.ListenAndServe(":8080", c.Handler(h2c.NewHandler(mux, &http2.Server{}))); err != nil {
		log.Fatal(err)
	}
}
