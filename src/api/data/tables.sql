CREATE TABLE users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT NOT NULL,
    password TEXT NOT NULL,
    salt TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    created_by TEXT NOT NULL,
    status_id INTEGER NOT NULL
);


select * from users;