package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/aidan-neel/shulker/apps/server/pkg/utils"
)

type contextKey string

const userIDKey contextKey = "userID"
const responseWriterKey contextKey = "response_writer"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")

		if err == nil {
			userID, err := utils.ValidateToken(cookie.Value, "access")
			if err == nil {
				ctx := context.WithValue(r.Context(), userIDKey, userID)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func InjectResponseWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), responseWriterKey, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetAuthCookies(ctx context.Context, accessToken, refreshToken string) {
	w, ok := ctx.Value(responseWriterKey).(http.ResponseWriter)
	if !ok {
		log.Println("❌ no response writer in context")
		return
	}
	log.Println("✅ got response writer, setting cookies")

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   15 * 60,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})
}

func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
