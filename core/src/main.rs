use tonic::{transport::Server, Request, Response, Status};
use user::user_service_server::{UserService, UserServiceServer};
use user::{UserRequest, UserResponse};

pub mod user {
    tonic::include_proto!("user");
}

#[derive(Default)]
pub struct MyUserService;

#[tonic::async_trait]
impl UserService for MyUserService {
    async fn get_user(&self, request: Request<UserRequest>) -> Result<Response<UserResponse>, Status> {
        let req = request.into_inner();
        let user = UserResponse {
            id: req.id,
            name: "John Doe".to_string(),
            email: "johndoe@gmail.com".to_string(),
        };
        Ok(Response::new(user))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let user_service = MyUserService::default();

    println!("UserService listening on {}", addr);

    Server::builder()
        .add_service(UserServiceServer::new(user_service))
        .serve(addr)
        .await?;

    Ok(())
}
