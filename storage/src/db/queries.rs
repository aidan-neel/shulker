use std::sync::{Arc, Mutex};

use rusqlite::{Connection, params};
use crate::models::file::File;
use tokio::task;

pub async fn insert_file_async(conn: Arc<Mutex<Connection>>, file: File) -> rusqlite::Result<()> {  
    let conn = conn.clone();
    task::spawn_blocking(move || {
        let conn = conn.lock().unwrap();
        conn.execute(
            "INSERT INTO file (file_path, file_name, file_size) VALUES (?1, ?2, ?3)",
            params![file.file_path, file.file_name, file.file_size as i64],
        )?;
        Ok(())
    })
    .await
    .unwrap() // propagate panics
}