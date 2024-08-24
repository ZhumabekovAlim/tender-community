CREATE TABLE notify_tokens
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    user_id       INT,
    token  VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES users (id)
);