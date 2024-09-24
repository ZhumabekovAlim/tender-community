CREATE TABLE transactions
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    type           VARCHAR(255)   NOT NULL,
    tender_number  VARCHAR(255),
    user_id        INT,
    company_id     INT,
    organization   VARCHAR(255),
    amount         DECIMAL(10, 2) NOT NULL,
    total          DECIMAL(10, 2) NOT NULL,
    sell           DECIMAL(10, 2) NOT NULL,
    product_name   VARCHAR(255),
    completed_date TIMESTAMP,
    date           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status         TINYINT        NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (company_id) REFERENCES companies (id)
);