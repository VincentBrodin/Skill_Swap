CREATE TABLE users(
    user_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    username TEXT NOT NULL,
    description TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);
