// Mirrors back-end/model/*.go

export type CompanyStatus = "active" | "suspended";

export interface Company {
  id: string;
  name: string;
  status: CompanyStatus;
  created_at: string;
  updated_at: string;
}

export interface AddCompanyRequest {
  name: string;
  status: CompanyStatus;
}

export type DocType =
  | "return_policy"
  | "warranty_policy"
  | "shipping_policy"
  | "refund_policy"
  | "complaint_sop"
  | "faq"
  | "other";

export const DOC_TYPES: { value: DocType; label: string }[] = [
  { value: "return_policy", label: "Return policy" },
  { value: "warranty_policy", label: "Warranty policy" },
  { value: "shipping_policy", label: "Shipping policy" },
  { value: "refund_policy", label: "Refund policy" },
  { value: "complaint_sop", label: "Complaint SOP" },
  { value: "faq", label: "FAQ" },
  { value: "other", label: "Other" },
];

export type DocVersionStatus = "draft" | "active" | "archived";

export const DOC_VERSION_STATUSES: { value: DocVersionStatus; label: string }[] = [
  { value: "draft", label: "Draft" },
  { value: "active", label: "Active" },
  { value: "archived", label: "Archived" },
];

export interface AddDocumentsMeta {
  company_id: number;
  document_type: DocType;
  version: number;
  status: DocVersionStatus;
  effective_from: string; // ISO date
  effective_until: string; // ISO date
}

export interface AddDocumentsResponse {
  company_id: number;
  name: number;
  document_type: DocType;
  file_type: string;
  mime_type: string;
  created_at: string;
}

export interface ChatRequest {
  company_id: number;
  message: string;
  limit?: number;
}

export interface ChatSourceChunk {
  document_id: number;
  document_version_id: number;
  chunk_index: number;
  content: string;
  score: number;
}

export interface ChatResponse {
  answer: string;
  sources: ChatSourceChunk[];
}

export type StatusEnum = "SUCCESS" | "FAILED";

export interface ErrorResponse {
  message: string;
  status: StatusEnum;
  error_code: number;
}

export interface SuccessResponse<T> {
  message: string;
  status: StatusEnum;
  code: number;
  data: T;
}
