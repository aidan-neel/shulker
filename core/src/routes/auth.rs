use axum::{extract::State, response::Json, routing::get, Router};
use axum_macros::debug_handler;
use serde::Serialize;

#[derive(Serialize)]
pub struct AuthResponse {
    pub message: String,
}

#[debug_handler]
pub async fn register() -> Json<AuthResponse> {
    Json(AuthResponse {
        message: "register route".into(),
    })
}

pub fn router() -> Router<()> {
    Router::new()
        .route("/register", get(register))
}