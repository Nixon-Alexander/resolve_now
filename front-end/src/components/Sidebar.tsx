import { NavLink } from "react-router-dom";

export default function Sidebar() {
  return (
    <aside className="sidebar">
      <div className="sidebar-brand">
        <span className="sidebar-brand-mark">RN</span>
        <div className="sidebar-brand-text">
          <span className="sidebar-brand-name">Resolve Now</span>
          <span className="sidebar-brand-sub">Policy &amp; case console</span>
        </div>
      </div>

      <nav className="sidebar-nav">
        <span className="sidebar-nav-label">Workspace</span>
        <NavLink
          to="/companies"
          className={({ isActive }) => "sidebar-link" + (isActive ? " is-active" : "")}
        >
          Companies
        </NavLink>
        <NavLink
          to="/chat"
          className={({ isActive }) => "sidebar-link" + (isActive ? " is-active" : "")}
        >
          Ask the docs
        </NavLink>
      </nav>

      <div className="sidebar-footer">
        <p>
          Upload a company&rsquo;s policy documents, then ask questions and get
          answers grounded in those documents.
        </p>
      </div>
    </aside>
  );
}
