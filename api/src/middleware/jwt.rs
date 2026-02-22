use axum::{
    extract::State,
    http::{Request, StatusCode},
    middleware::Next,
    response::Response,
};
use axum_extra::extract::cookie::CookieJar;
use crate::routes::AppState;

#[derive(Clone)]
pub struct UserId(pub i32);

pub async fn jwt_auth(
    State(state): State<AppState>,
    jar: CookieJar, // Extract the CookieJar
    mut req: Request<axum::body::Body>,
    next: Next,
) -> Result<Response, (StatusCode, String)> {
    let token = jar
        .get("access_token")
        .map(|cookie| cookie.value())
        .ok_or((StatusCode::UNAUTHORIZED, "Missing authentication cookie".into()))?;
    
    let claims = state
        .jwt
        .verify_token(token, false)
        .map_err(|_| (StatusCode::UNAUTHORIZED, "Invalid or expired token".into()))?;
    
    let user_id_str = claims
        .subject
        .ok_or((StatusCode::UNAUTHORIZED, "Missing subject in token".into()))?;
        
    let i32_user_id = user_id_str
        .parse::<i32>()
        .map_err(|e| (StatusCode::INTERNAL_SERVER_ERROR, format!("ID parse error: {}", e)))?;

    req.extensions_mut().insert(UserId(i32_user_id));

    Ok(next.run(req).await)
}