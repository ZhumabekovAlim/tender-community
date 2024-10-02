CREATE TABLE tenders
(
    id             INT AUTO_INCREMENT PRIMARY KEY,
    type           VARCHAR(255)   NOT NULL,
    tender_number  VARCHAR(255),
    user_id        INT,
    company_id     INT,
    organization   VARCHAR(255),
    total          DECIMAL(10, 2) NOT NULL,
    commission      DECIMAL(10, 2) NOT NULL,
    completed_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    date           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status         TINYINT        NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (company_id) REFERENCES companies (id)
);