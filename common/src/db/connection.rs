use rusqlite::{Connection, Result};

pub fn establish_connection() -> Result<Connection> {
    let path = "../data/db.sqlite3";  
    if let Some(parent) = std::path::Path::new(path).parent() {
        std::fs::create_dir_all(parent).map_err(|e| {
            eprintln!("Failed to create database directory: {}", e);
        }).ok(); 
    }
    let conn = Connection::open(path)?;
    Ok(conn)
}