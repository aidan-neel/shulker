mod grpc;

use tonic::{transport::Server};
use upload::upload_service_server::UploadServiceServer;
use grpc::upload::UploadServiceImpl;

pub mod upload {
    tonic::include_proto!("upload");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let upload_service = UploadServiceImpl::default();

    println!("UploadService listening on {}", addr);

    Server::builder()
        .add_service(UploadServiceServer::new(upload_service))
        .serve(addr)
        .await?;

    Ok(())
}
