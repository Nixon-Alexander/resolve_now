import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import { getCompanies } from "../api/companies";
import { ApiError } from "../api/client";
import type { Company } from "../types";
import ChatPanel from "../components/ChatPanel";

interface LocationState {
  companyId?: number;
}

export default function ChatPage() {
  const location = useLocation();
  const [companies, setCompanies] = useState<Company[]>([]);
  const [selectedId, setSelectedId] = useState<number | null>(
    (location.state as LocationState | null)?.companyId ?? null
  );
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    getCompanies()
      .then((data) => {
        setCompanies(data ?? []);
        if (!selectedId && data && data.length > 0) {
          setSelectedId(Number(data[0].id));
        }
      })
      .catch((err) => setError(err instanceof ApiError ? err.message : "Could not load companies."))
      .finally(() => setLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="page">
      <header className="page-header">
        <div>
          <h1>Ask the docs</h1>
          <p className="page-subtitle">
            Questions are answered from the policy documents uploaded for the selected company.
          </p>
        </div>
      </header>

      {loading && <p className="muted">Loading companies…</p>}
      {error && <p className="form-error">{error}</p>}

      {!loading && !error && (
        <>
          <label className="field company-picker">
            <span>Company</span>
            <select
              value={selectedId ?? ""}
              onChange={(e) => setSelectedId(e.target.value ? Number(e.target.value) : null)}
            >
              <option value="" disabled>
                Choose a company
              </option>
              {companies.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </label>

          <ChatPanel companyId={selectedId} />
        </>
      )}
    </div>
  );
}
