CREATE TABLE IF NOT EXISTS user (
     id INT AUTO_INCREMENT,
     email VARCHAR(50) NOT NULL,
     plan VARCHAR(20) DEFAULT 'free_plan',
     plan_bought_date DATETIME NOT NULL,
     messages_left INT,
     bytes_data_left INT,
     bots_left INT,
     PRIMARY KEY (id),
     INDEX (email)
);
