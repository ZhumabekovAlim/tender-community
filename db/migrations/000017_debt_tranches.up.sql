CREATE TABLE debt_tranches
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    debt_id     INT,
    amount      INT,
    description VARCHAR(255),
    date        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (debt_id) REFERENCES transactions (id)
);