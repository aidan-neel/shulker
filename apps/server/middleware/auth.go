package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/aidan-neel/shulker/apps/server/jwt"
)

type contextKey string

const (
	userIDKey         contextKey = "userID"
	responseWriterKey contextKey = "responseWriter"
)

// Auth extracts the user from the access_token cookie if present.
// It does not enforce authentication — use RequireAuth for that.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err == nil {
			if userID, err := jwt.ValidateToken(cookie.Value, "access"); err == nil {
				ctx := context.WithValue(r.Context(), userIDKey, userID)
				r = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, r)
	})
}

// InjectResponseWriter injects the ResponseWriter into the context so handlers
// can set cookies (required by Connect RPC which doesn't expose the writer directly).
func InjectResponseWriter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), responseWriterKey, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SetAuthCookies writes access and refresh token cookies onto the response.
func SetAuthCookies(ctx context.Context, accessToken, refreshToken string) {
	w, ok := ctx.Value(responseWriterKey).(http.ResponseWriter)
	if !ok {
		log.Println("no response writer in context")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   15 * 60,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})
}

// GetUserID returns the authenticated user's ID from the request context.
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
