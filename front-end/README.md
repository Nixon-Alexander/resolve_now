# Resolve Now — Front End

React + TypeScript (Vite) front end for the `resolve_now` Go backend. It gives
you three things:

1. **Companies** — list companies and add new ones.
2. **Documents** — upload a company's policy documents (return, warranty,
   shipping, refund, complaint SOP, FAQ, other) with version/status/effective
   dates, matching `POST /api/v1/doc/add`.
3. **Ask the docs** — chat-style search that calls `POST /api/v1/chat/search`
   and shows the grounded answer plus the source chunks it was built from.

## Backend endpoints used

| Method | Path                        | Used for                     |
|--------|-----------------------------|-------------------------------|
| GET    | `/api/v1/get-companies`     | Companies list                |
| GET    | `/api/v1/get-company/:id`   | Company detail                |
| POST   | `/api/v1/add-company`       | Create company                |
| POST   | `/api/v1/doc/add`           | Upload documents (multipart)  |
| POST   | `/api/v1/chat/search`       | Ask a question (RAG search)   |

## Getting started

```bash
npm install
cp .env.example .env   # then edit VITE_API_BASE_URL if needed
npm run dev
```

By default the app talks to `http://localhost:8080/api/v1` — set
`VITE_API_BASE_URL` in `.env` to point at wherever the Go backend (`make start`
in `back-end/`) is actually running, including the `/api/v1` prefix.

## Scripts

- `npm run dev` — start the Vite dev server
- `npm run build` — type-check (`tsc -b`) and build for production into `dist/`
- `npm run preview` — preview the production build locally
- `npm run lint` — run ESLint

## Project structure

```
src/
  api/         fetch wrappers per resource (companies, documents, chat)
  types/       TypeScript types mirroring the Go models
  components/  Sidebar, forms, lists, chat panel, status badge
  pages/       CompaniesPage, CompanyDetailPage, ChatPage
  App.tsx      routes
  main.tsx     entry point (BrowserRouter)
  index.css    design tokens + all styles (no CSS framework)
```
