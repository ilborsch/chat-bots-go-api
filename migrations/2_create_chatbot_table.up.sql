CREATE TABLE IF NOT EXISTS chat_bot (
    id INT AUTO_INCREMENT,
    assistant_id VARCHAR(50) NOT NULL UNIQUE,
    vector_store_id VARCHAR(50) NOT NULL UNIQUE,
    owner_id INT NOT NULL,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    instructions TEXT,
    PRIMARY KEY (id),
    FOREIGN KEY (owner_id) REFERENCES user(id),
    INDEX (owner_id)
);
