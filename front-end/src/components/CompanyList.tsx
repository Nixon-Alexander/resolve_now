import { Link } from "react-router-dom";
import type { Company } from "../types";
import StatusBadge from "./StatusBadge";

interface CompanyListProps {
  companies: Company[];
}

function formatDate(value: string): string {
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return value;
  return d.toLocaleDateString(undefined, { year: "numeric", month: "short", day: "numeric" });
}

export default function CompanyList({ companies }: CompanyListProps) {
  if (companies.length === 0) {
    return (
      <div className="empty-state">
        <p>No companies yet.</p>
        <p className="empty-state-sub">Add a company to start uploading its policy documents.</p>
      </div>
    );
  }

  return (
    <ul className="docket">
      {companies.map((company) => (
        <li key={company.id} className="docket-row">
          <Link to={`/companies/${company.id}`} className="docket-row-link">
            <span className="docket-row-name">{company.name}</span>
            <span className="docket-row-meta">
              Added {formatDate(company.created_at)}
            </span>
            <StatusBadge
              status={company.status}
              tone={company.status === "active" ? "positive" : "warning"}
            />
          </Link>
        </li>
      ))}
    </ul>
  );
}
