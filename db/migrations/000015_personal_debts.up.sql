CREATE TABLE personal_debts
(
    id   INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL ,
    amount double(10,2) NOT NULL ,
    type varchar(20) NOT NULL ,
    get_date date,
    return_date date ,
    status INT NOT NULL ,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
