pub mod auth;

use axum::Router;

pub fn create_router() -> Router {
    Router::new()
        .nest("/auth", auth::router())
        .fallback(handler_404)
}

async fn handler_404() -> impl axum::response::IntoResponse {
    (axum::http::StatusCode::NOT_FOUND, "Route not found - check your nesting!")
}