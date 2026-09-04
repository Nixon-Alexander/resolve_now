package model

import "time"

type DocVersionStatus string

const (
	DOC_DRAFT    DocVersionStatus = "draft"
	DOC_ACTIVE   DocVersionStatus = "active"
	DOC_ARCHIVED DocVersionStatus = "archived"
)

type AddDocumentVersionRequest struct {
	DocumentID    int              `json:"document_id"`
	Version       int              `json:"version"`
	FilePath      string           `json:"file_path"`
	Status        DocVersionStatus `json:"status"`
	EffectiveFrom time.Time        `json:"effective_from"`
	EffeciveUntil time.Time        `json:"effective_until"`
}
