package model

import "time"

type DocType string

const (
	RETURN_POLICY   DocType = "return_policy"
	WARRANTY_POLICY DocType = "warranty_policy"
	SHIPPING_POLICY DocType = "shipping_policy"
	REFUND_POLICY   DocType = "refund_policy"
	COMPLAINT_SOP   DocType = "complaint_sop"
	FAQ             DocType = "faq"
	OTHER           DocType = "other"
)

type AddDocumentsRequest struct {
	CompanyId      int              `form:"company_id"`
	DocumentType   DocType          `form:"document_type"`
	Version        int              `form:"version"`
	Status         DocVersionStatus `form:"status"`
	EffectiveFrom  time.Time        `form:"effective_from"`
	EffectiveUntil time.Time        `form:"effective_until"`
}

type AddDocumentsResponse struct {
	CompanyId    int       `json:"company_id"`
	Name         int       `json:"name"`
	DocumentType DocType   `json:"document_type"`
	FilePath     string    `json:"file_type"`
	MimeType     string    `json:"mime_type"`
	CreatedAt    time.Time `json:"created_at"`
}

type UploadedFile struct {
	FilePath string
	FileName string
	MimeType string
}

type DocumentSection struct {
	Title   string
	Content string
}

type DocumentChunk struct {
	Index        int
	SectionTitle string
	Content      string
}

type ChunkSearchResult struct {
	ChunkID           uint64
	DocumentID        int64
	DocumentVersionID int64
	ChunkIndex        int64
	Content           string
	Score             float32
}

type ChatRequest struct {
	CompanyId int64  `json:"company_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Limit     int    `json:"limit"` // opsional, default 5
}

type ChatResponse struct {
	Answer  string            `json:"answer"`  // chunk paling relevan (top match)
	Sources []ChatSourceChunk `json:"sources"` // semua chunk yang ketemu, urut dari paling relevan
}

type ChatSourceChunk struct {
	DocumentID        int64   `json:"document_id"`
	DocumentVersionID int64   `json:"document_version_id"`
	ChunkIndex        int64   `json:"chunk_index"`
	Content           string  `json:"content"`
	Score             float32 `json:"score"`
}
