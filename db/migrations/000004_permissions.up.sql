CREATE TABLE permissions
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT NOT NULL,
    company_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (company_id) REFERENCES companies (id)
);