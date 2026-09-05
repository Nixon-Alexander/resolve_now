import { apiGet, apiPostJson } from "./client";
import type { AddCompanyRequest, Company } from "../types";

export function getCompanies(): Promise<Company[]> {
  return apiGet<Company[]>("/get-companies");
}

export function getCompany(id: string): Promise<Company[]> {
  return apiGet<Company[]>(`/get-company/${id}`);
}

export function addCompany(payload: AddCompanyRequest): Promise<unknown> {
  return apiPostJson("/add-company", payload);
}
