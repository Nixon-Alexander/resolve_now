import { apiPostForm } from "./client";
import type { AddDocumentsMeta, AddDocumentsResponse } from "../types";

export function uploadDocuments(
  meta: AddDocumentsMeta,
  files: File[]
): Promise<AddDocumentsResponse> {
  const form = new FormData();
  form.append("company_id", String(meta.company_id));
  form.append("document_type", meta.document_type);
  form.append("version", String(meta.version));
  form.append("status", meta.status);
  form.append("effective_from", meta.effective_from);
  form.append("effective_until", meta.effective_until);

  for (const file of files) {
    form.append("files", file);
  }

  return apiPostForm<AddDocumentsResponse>("/doc/add", form);
}
