CREATE TABLE additional_expenses
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name      VARCHAR(255)   NOT NULL,
    amount      DECIMAL(10, 2) NOT NULL,
    transaction_id INT NOT NULL ,
    date        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (transaction_id) REFERENCES transactions (id)
);