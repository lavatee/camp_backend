CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(200) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(20) NOT NULL,
    tag VARCHAR(20) NOT NULL UNIQUE,
    about VARCHAR(200),
    photo_url VARCHAR(200),
    language VARCHAR(200)
);

CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    tree INT DEFAULT 0,
    last_tree_update TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users_in_chat (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    chat_id INT NOT NULL,
    translating BOOLEAN
);

CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    text VARCHAR(255) NOT NULL,
    chat_id INT NOT NULL,
    user_id INT NOT NULL,
    sent_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN
);

CREATE TABLE rooms (
    id SERIAL PRIMARY KEY
);

CREATE TABLE users_in_room (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    room_id INT NOT NULL,
    translating BOOLEAN
);

ALTER TABLE users_in_chat ADD FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE users_in_chat ADD FOREIGN KEY (chat_id) REFERENCES chats(id);
ALTER TABLE friends ADD FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE users_in_room ADD FOREIGN KEY (user_id) REFERENCES users(id);
ALTER TABLE users_in_room ADD FOREIGN KEY (room_id) REFERENCES rooms(id);
