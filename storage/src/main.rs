use tonic::transport::Channel;
use user::user_service_client::UserServiceClient;
use user::{UserRequest, UserResponse};

pub mod user {
    tonic::include_proto!("user");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = UserServiceClient::connect("http://[::1]:50051").await?;

    let request = tonic::Request::new(UserRequest {
        id: 4
    });

    let response: tonic::Response<UserResponse> = client.get_user(request).await?;

    println!("Response: {:?}", response.into_inner());

    Ok(())
}
