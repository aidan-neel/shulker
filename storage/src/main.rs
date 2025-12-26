mod grpc;

use tonic::transport::Server;
use upload::upload_service_server::UploadServiceServer;
use grpc::upload::UploadServiceImpl;

use std::sync::{Arc, Mutex};
use common::db::connection::establish_connection;
use rusqlite::Connection;

use common::db::schema::init_db;

pub mod upload {
    tonic::include_proto!("upload");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;

    let conn = Arc::new(Mutex::new(establish_connection()?));
    let _ = init_db(conn.clone());
    let upload_service = UploadServiceImpl { conn: conn.clone() };

    println!("UploadService listening on {}", addr);

    Server::builder()
        .add_service(UploadServiceServer::new(upload_service))
        .serve(addr)
        .await?;

    Ok(())  
}
