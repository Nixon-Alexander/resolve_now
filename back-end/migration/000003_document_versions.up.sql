CREATE TABLE IF NOT EXISTS document_versions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    document_id BIGINT UNSIGNED NOT NULL,
    version INT UNSIGNED NOT NULL,
    file_path TEXT,
    status ENUM(
        'draft',
        'active',
        'archived'
    ) DEFAULT 'draft',
    effective_from DATETIME,
    effective_until DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY unique_document_version (
        document_id,
        version
    ),

    FOREIGN KEY (document_id)
        REFERENCES documents(id)
);
