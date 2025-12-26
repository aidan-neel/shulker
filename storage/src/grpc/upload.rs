use crate::upload::{UploadRequest, UploadResponse};
use crate::upload::upload_service_server::UploadService;
use tokio_stream::StreamExt;    
use tokio::fs::{create_dir_all, OpenOptions};
use tokio::io::AsyncWriteExt;

#[derive(Default)]
pub struct UploadServiceImpl;

#[tonic::async_trait]
impl UploadService for UploadServiceImpl {
    async fn upload_file(
        &self,
        request: tonic::Request<tonic::Streaming<UploadRequest>>,
    ) -> Result<tonic::Response<UploadResponse>, tonic::Status> {
        let mut stream = request.into_inner();
        let mut final_data = Vec::new();

        while let Some(result) = stream.next().await {
            match result {
                Ok(chunk) => {
                    let folder_path = "./data/user1/";  
                    let file_path = format!("{}{}", folder_path, chunk.filename);
                    create_dir_all(folder_path)
                        .await
                        .map_err(|e| tonic::Status::internal(format!("Failed to create directory: {}", e)))?;
                    
                    let mut file = OpenOptions::new()
                        .create(true)
                        .append(true)
                        .open(file_path)
                        .await
                        .map_err(|e| tonic::Status::internal(format!("File error: {}", e)))?;

                    file.write_all(&chunk.data)
                        .await
                        .map_err(|e| tonic::Status::internal(format!("Failed to write chunk: {}", e)))?;

                    final_data.extend(chunk.data);
                }
                Err(err) => {
                    return Err(tonic::Status::internal(format!("Stream error {:?}", err)))
                }
            }
        }

        Ok(tonic::Response::new(UploadResponse {
            id: format!("Saved {} bytes", final_data.len())
        }))
    }
}