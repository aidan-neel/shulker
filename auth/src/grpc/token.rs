use std::sync::{Arc, Mutex};
use rusqlite::Connection;
use crate::proto::{GetTokenRequest, GetTokenResponse, RefreshTokenRequest, RefreshTokenResponse, Token};
use crate::proto::token_service_server::TokenService;

#[derive(Clone)]
pub struct TokenServiceImpl {
    pub conn: Arc<Mutex<Connection>>
}

#[tonic::async_trait]
impl TokenService for TokenServiceImpl {
    async fn get_token(
        &self,
        request: tonic::Request<GetTokenRequest>
    ) -> Result<tonic::Response<GetTokenResponse>, tonic::Status> {
        let mut data = request.into_inner();
        Ok(tonic::Response::new(GetTokenResponse {
            token: Some(Token {
                access_token: "hey".to_string(),
                refresh_token: "hello".to_string(),
            }),
        }))
    }

    async fn refresh_token(
        &self,
        request: tonic::Request<RefreshTokenRequest>
    ) -> Result<tonic::Response<RefreshTokenResponse>, tonic::Status> {
        let mut data = request.into_inner();
        Ok(tonic::Response::new(RefreshTokenResponse {
            token: Some(Token {
                access_token: "hey".to_string(),
                refresh_token: "hello".to_string(),
            }),
        }))
    }
}