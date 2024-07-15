CREATE TABLE transactions
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT            NOT NULL,
    company_id INT            NOT NULL,
    amount     DECIMAL(10, 2) NOT NULL,
    date       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status     TINYINT        NOT NULL CHECK (status IN (0, 1, 2)),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (company_id) REFERENCES companies (id)
);