import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { getCompany } from "../api/companies";
import { ApiError } from "../api/client";
import type { Company } from "../types";
import StatusBadge from "../components/StatusBadge";
import DocumentUploadForm from "../components/DocumentUploadForm";

export default function CompanyDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [company, setCompany] = useState<Company | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [uploadCount, setUploadCount] = useState(0);

  useEffect(() => {
    if (!id) return;
    // eslint-disable-next-line react-hooks/set-state-in-effect -- data fetch on mount/id change
    setLoading(true);
    setError(null);
    getCompany(id)
      .then((rows) => setCompany(rows?.[0] ?? null))
      .catch((err) => setError(err instanceof ApiError ? err.message : "Could not load this company."))
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <div className="page"><p className="muted">Loading…</p></div>;
  if (error) return <div className="page"><p className="form-error">{error}</p></div>;
  if (!company) return <div className="page"><p className="muted">Company not found.</p></div>;

  return (
    <div className="page">
      <Link to="/companies" className="back-link">
        ← All companies
      </Link>

      <header className="page-header">
        <div>
          <h1>{company.name}</h1>
          <p className="page-subtitle">Company #{company.id}</p>
        </div>
        <StatusBadge
          status={company.status}
          tone={company.status === "active" ? "positive" : "warning"}
        />
      </header>

      <div className="two-col">
        <section className="panel">
          <h2>Upload policy documents</h2>
          <p className="muted">
            Files are chunked and indexed automatically, so the chat assistant can cite them
            when answering questions for this company.
          </p>
          <DocumentUploadForm
            companyId={Number(company.id)}
            onUploaded={() => setUploadCount((c) => c + 1)}
          />
        </section>

        <section className="panel">
          <h2>Ask this company&rsquo;s documents</h2>
          <p className="muted">
            Once documents are uploaded, head to the chat page and pick{" "}
            <strong>{company.name}</strong> to ask questions grounded in what you just uploaded.
          </p>
          {uploadCount > 0 && (
            <p className="form-success">
              {uploadCount} upload{uploadCount > 1 ? "s" : ""} sent so far this session.
            </p>
          )}
          <Link to="/chat" state={{ companyId: Number(company.id) }} className="button button-primary">
            Go to chat
          </Link>
        </section>
      </div>
    </div>
  );
}
