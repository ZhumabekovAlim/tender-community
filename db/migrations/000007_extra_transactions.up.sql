CREATE TABLE extra_transactions
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    user_id       INT,
    description  VARCHAR(255),
    total         DECIMAL(10, 2) NOT NULL,
    date          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status        TINYINT        NOT NULL CHECK (status IN (0, 1, 2)),
    FOREIGN KEY (user_id) REFERENCES users (id)
);