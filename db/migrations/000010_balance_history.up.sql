CREATE TABLE balance_history
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    amount      DECIMAL(10, 2) NOT NULL,
    description VARCHAR(255)   NOT NULL,
    user_id     INT            NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);