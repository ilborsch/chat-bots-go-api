CREATE TABLE IF NOT EXISTS file (
    id INT AUTO_INCREMENT,
    owner_id INT,
    chat_bot_id INT NOT NULL,
    openai_file_id VARCHAR(50) NOT NULL,
    filename VARCHAR(50) DEFAULT 'A file with no name given...',
    PRIMARY KEY (id),
    FOREIGN KEY (owner_id) REFERENCES user(id),
    FOREIGN KEY (chat_bot_id) REFERENCES chat_bot(id),
    INDEX (chat_bot_id)
)

