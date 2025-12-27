use axum::{extract::State, response::Json, routing::post, Router};
use axum_macros::debug_handler;
use serde::Serialize;

#[derive(Serialize)]
pub struct AuthResponse {
    pub message: String,
}

#[derive(Deserialize)]
pub struct RegisterPayload {
    email: String,
    display_name: String,
    password: String,
}

// add to database here
// but use AuthService for:
//     access & refresh tokens
#[debug_handler]
pub async fn register(Json(register_payload): Json<RegisterPayload>) -> Json<AuthResponse> {
    

    Json(AuthResponse {
        message: "register route".into(),
    })
}

pub fn router() -> Router<()> {
    Router::new()
        .route("/register", post(register))
}