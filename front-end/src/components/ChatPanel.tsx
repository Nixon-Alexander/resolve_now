import { useState } from "react";
import type { FormEvent } from "react";
import { searchChat } from "../api/chat";
import { ApiError } from "../api/client";
import type { ChatResponse } from "../types";

interface ChatPanelProps {
  companyId: number | null;
}

interface Turn {
  question: string;
  response: ChatResponse;
}

export default function ChatPanel({ companyId }: ChatPanelProps) {
  const [message, setMessage] = useState("");
  const [turns, setTurns] = useState<Turn[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    if (!companyId) {
      setError("Choose a company first.");
      return;
    }
    if (!message.trim()) return;

    setLoading(true);
    setError(null);
    const question = message.trim();
    try {
      const response = await searchChat({ company_id: companyId, message: question });
      setTurns((prev) => [...prev, { question, response }]);
      setMessage("");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Something went wrong asking the docs.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="chat-panel">
      <div className="chat-log">
        {turns.length === 0 && (
          <div className="empty-state">
            <p>No questions asked yet.</p>
            <p className="empty-state-sub">
              Try something like &ldquo;what is the return window for damaged items?&rdquo;
            </p>
          </div>
        )}

        {turns.map((turn, i) => (
          <div className="chat-turn" key={i}>
            <p className="chat-question">{turn.question}</p>
            <div className="chat-answer">
              <p>{turn.response.answer}</p>
              {turn.response.sources.length > 0 && (
                <details className="chat-sources">
                  <summary>{turn.response.sources.length} source chunk(s)</summary>
                  <ul>
                    {turn.response.sources.map((s, j) => (
                      <li key={j} className="source-chunk">
                        <div className="source-chunk-head">
                          <span>Doc #{s.document_id} · chunk {s.chunk_index}</span>
                          <span>score {s.score.toFixed(3)}</span>
                        </div>
                        <p>{s.content}</p>
                      </li>
                    ))}
                  </ul>
                </details>
              )}
            </div>
          </div>
        ))}
      </div>

      {error && <p className="form-error">{error}</p>}

      <form className="chat-input-row" onSubmit={handleSubmit}>
        <input
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder={companyId ? "Ask a question about this company's policies…" : "Choose a company first"}
          disabled={!companyId}
        />
        <button type="submit" className="button button-primary" disabled={!companyId || loading}>
          {loading ? "Asking…" : "Ask"}
        </button>
      </form>
    </div>
  );
}
