CREATE TABLE IF NOT EXISTS document_chunks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    document_version_id BIGINT UNSIGNED NOT NULL,
    chunk_index INT UNSIGNED NOT NULL,
    content TEXT NOT NULL,
    token_count INT UNSIGNED,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_version_id)
        REFERENCES document_versions(id),
    UNIQUE KEY unique_chunk (
        document_version_id,
        chunk_index
    )
);