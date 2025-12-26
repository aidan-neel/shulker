use tonic::transport::Channel;
use tokio_stream::iter;

pub mod proto {
    tonic::include_proto!("upload"); 
}

use proto::upload_service_client::UploadServiceClient;
use proto::{UploadRequest, UploadResponse};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = UploadServiceClient::connect("http://[::1]:50051").await?;

    let chunk = UploadRequest {
        data: vec![1, 2, 3, 4, 5],
        filename: "example.txt".to_string(),
    };
    
    let request_stream = iter(vec![chunk]);

    let response: tonic::Response<UploadResponse> = client.upload_file(request_stream).await?;
    println!("Response: {:?}", response.into_inner());

    Ok(())
}
