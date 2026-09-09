CREATE TABLE IF NOT EXISTS documents (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    company_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    document_type ENUM(
        'return_policy',
        'warranty_policy',
        'shipping_policy',
        'refund_policy',
        'complaint_sop',
        'faq',
        'other'
    ) NOT NULL,
    file_path TEXT,
    mime_type VARCHAR(100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id)
        REFERENCES companies(id)
);