import { useState } from "react";
import type { FormEvent } from "react";
import { addCompany } from "../api/companies";
import type { CompanyStatus } from "../types";
import { ApiError } from "../api/client";

interface CompanyFormProps {
  onCreated: () => void;
  onClose: () => void;
}

export default function CompanyForm({ onCreated, onClose }: CompanyFormProps) {
  const [name, setName] = useState("");
  const [status, setStatus] = useState<CompanyStatus>("active");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!name.trim()) {
      setError("Give the company a name before saving.");
      return;
    }
    setSubmitting(true);
    setError(null);
    try {
      await addCompany({ name: name.trim(), status });
      onCreated();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Could not save the company.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form className="panel-form" onSubmit={handleSubmit}>
      <div className="panel-form-header">
        <h2>New company</h2>
        <button type="button" className="icon-button" onClick={onClose} aria-label="Close">
          &times;
        </button>
      </div>

      <label className="field">
        <span>Company name</span>
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="e.g. Kopi Kenangan"
          autoFocus
        />
      </label>

      <label className="field">
        <span>Status</span>
        <select value={status} onChange={(e) => setStatus(e.target.value as CompanyStatus)}>
          <option value="active">Active</option>
          <option value="suspended">Suspended</option>
        </select>
      </label>

      {error && <p className="form-error">{error}</p>}

      <div className="panel-form-actions">
        <button type="button" className="button button-ghost" onClick={onClose}>
          Cancel
        </button>
        <button type="submit" className="button button-primary" disabled={submitting}>
          {submitting ? "Saving…" : "Save company"}
        </button>
      </div>
    </form>
  );
}
