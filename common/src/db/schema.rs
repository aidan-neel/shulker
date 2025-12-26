use std::sync::{Arc, Mutex};

use rusqlite::Connection;

pub fn init_db(conn: Arc<Mutex<Connection>>) -> rusqlite::Result<()> {
    let conn = conn.lock().unwrap();
    
    conn.execute(
        "CREATE TABLE IF NOT EXISTS file (
            id INTEGER PRIMARY KEY,
            file_path TEXT NOT NULL,
            file_name TEXT NOT NULL,
            file_size INTEGER NOT NULL
        )",
        [],
    )?;

    conn.execute(
        "CREATE TABLE IF NOT EXISTS user (
            id INTEGER PRIMARY KEY,
            display_name TEXT NOT NULL,
            email TEXT NOT NULL,
            password_hash TEXT NOT NULL
        )",
        [],
    )?;
    Ok(())
}
