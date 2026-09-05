import { useEffect, useState } from "react";
import { getCompanies } from "../api/companies";
import { ApiError } from "../api/client";
import type { Company } from "../types";
import CompanyList from "../components/CompanyList";
import CompanyForm from "../components/CompanyForm";

export default function CompaniesPage() {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);

  async function load() {
    setLoading(true);
    setError(null);
    try {
      const data = await getCompanies();
      setCompanies(data ?? []);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Could not load companies.");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect -- initial data fetch on mount
    load();
  }, []);

  return (
    <div className="page">
      <header className="page-header">
        <div>
          <h1>Companies</h1>
          <p className="page-subtitle">
            Every company whose policy documents are on file for the chat assistant to search.
          </p>
        </div>
        <button className="button button-primary" onClick={() => setShowForm(true)}>
          New company
        </button>
      </header>

      {showForm && (
        <div className="panel">
          <CompanyForm
            onCreated={() => {
              setShowForm(false);
              load();
            }}
            onClose={() => setShowForm(false)}
          />
        </div>
      )}

      {loading && <p className="muted">Loading companies…</p>}
      {error && <p className="form-error">{error}</p>}
      {!loading && !error && <CompanyList companies={companies} />}
    </div>
  );
}
