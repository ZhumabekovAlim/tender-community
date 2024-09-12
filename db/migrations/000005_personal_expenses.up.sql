CREATE TABLE personal_expenses
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    amount      DECIMAL(10, 2) NOT NULL,
    reason      VARCHAR(255)   NOT NULL,
    description TEXT,
    category_id INT NOT NULL ,
    date        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories (id)
);