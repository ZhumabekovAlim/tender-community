CREATE TABLE tranches
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    transaction_id INT,
    amount         INT,
    description    VARCHAR(255),
    date           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions (id)
);