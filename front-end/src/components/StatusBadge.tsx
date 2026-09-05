interface StatusBadgeProps {
  status: string;
  tone?: "positive" | "neutral" | "warning";
}

export default function StatusBadge({ status, tone = "neutral" }: StatusBadgeProps) {
  return <span className={`badge badge-${tone}`}>{status}</span>;
}
