import { useRef, useState } from "react";
import type { FormEvent } from "react";
import { uploadDocuments } from "../api/documents";
import { ApiError } from "../api/client";
import { DOC_TYPES, DOC_VERSION_STATUSES } from "../types";
import type { DocType, DocVersionStatus } from "../types";

interface DocumentUploadFormProps {
  companyId: number;
  onUploaded: () => void;
}

function todayIso(): string {
  return new Date().toISOString().slice(0, 10);
}

export default function DocumentUploadForm({ companyId, onUploaded }: DocumentUploadFormProps) {
  const [documentType, setDocumentType] = useState<DocType>("return_policy");
  const [version, setVersion] = useState(1);
  const [status, setStatus] = useState<DocVersionStatus>("active");
  const [effectiveFrom, setEffectiveFrom] = useState(todayIso());
  const [effectiveUntil, setEffectiveUntil] = useState("");
  const [files, setFiles] = useState<File[]>([]);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setSuccess(null);

    if (files.length === 0) {
      setError("Attach at least one file.");
      return;
    }

    setSubmitting(true);
    try {
      await uploadDocuments(
        {
          company_id: companyId,
          document_type: documentType,
          version,
          status,
          effective_from: `${effectiveFrom}T00:00:00Z`,
          effective_until: effectiveUntil ? `${effectiveUntil}T00:00:00Z` : `${effectiveFrom}T00:00:00Z`,
        },
        files
      );
      setSuccess(files.length === 1 ? "Document uploaded." : `${files.length} documents uploaded.`);
      setFiles([]);
      if (fileInputRef.current) fileInputRef.current.value = "";
      onUploaded();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Upload failed. Try again.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form className="upload-form" onSubmit={handleSubmit}>
      <div className="field-grid">
        <label className="field">
          <span>Document type</span>
          <select value={documentType} onChange={(e) => setDocumentType(e.target.value as DocType)}>
            {DOC_TYPES.map((t) => (
              <option key={t.value} value={t.value}>
                {t.label}
              </option>
            ))}
          </select>
        </label>

        <label className="field">
          <span>Version</span>
          <input
            type="number"
            min={1}
            value={version}
            onChange={(e) => setVersion(Number(e.target.value))}
          />
        </label>

        <label className="field">
          <span>Status</span>
          <select value={status} onChange={(e) => setStatus(e.target.value as DocVersionStatus)}>
            {DOC_VERSION_STATUSES.map((s) => (
              <option key={s.value} value={s.value}>
                {s.label}
              </option>
            ))}
          </select>
        </label>

        <label className="field">
          <span>Effective from</span>
          <input
            type="date"
            value={effectiveFrom}
            onChange={(e) => setEffectiveFrom(e.target.value)}
          />
        </label>

        <label className="field">
          <span>Effective until</span>
          <input
            type="date"
            value={effectiveUntil}
            onChange={(e) => setEffectiveUntil(e.target.value)}
            placeholder="Optional"
          />
        </label>
      </div>

      <label className="field">
        <span>Files</span>
        <input
          ref={fileInputRef}
          type="file"
          multiple
          onChange={(e) => setFiles(Array.from(e.target.files ?? []))}
        />
      </label>

      {files.length > 0 && (
        <ul className="file-chip-list">
          {files.map((f) => (
            <li key={f.name} className="file-chip">
              {f.name}
            </li>
          ))}
        </ul>
      )}

      {error && <p className="form-error">{error}</p>}
      {success && <p className="form-success">{success}</p>}

      <div className="panel-form-actions">
        <button type="submit" className="button button-primary" disabled={submitting}>
          {submitting ? "Uploading…" : "Upload documents"}
        </button>
      </div>
    </form>
  );
}
