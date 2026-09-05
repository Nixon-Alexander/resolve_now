import { apiPostJson } from "./client";
import type { ChatRequest, ChatResponse } from "../types";

export function searchChat(payload: ChatRequest): Promise<ChatResponse> {
  return apiPostJson<ChatResponse>("/chat/search", payload);
}
