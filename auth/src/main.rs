mod grpc;

use tonic::transport::Server;
use proto::token_service_server::TokenServiceServer;
use grpc::token::TokenServiceImpl;

use std::sync::{Arc, Mutex};
use common::db::connection::establish_connection;
use rusqlite::Connection;

use common::db::schema::init_db;

pub mod proto {
    tonic::include_proto!("auth");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;

    let conn = Arc::new(Mutex::new(establish_connection()?));
    let _ = init_db(conn.clone());
    let token_service = TokenServiceImpl { conn: conn.clone() };

    println!("TokenService listening on {}", addr);

    Server::builder()
        .add_service(TokenServiceServer::new(token_service))
        .serve(addr)
        .await?;

    Ok(())  
}
