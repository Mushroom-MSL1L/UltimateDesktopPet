import { useEffect, useMemo, useState } from "react";
import { GetPetStatus } from "../../../wailsjs/go/pet/PetMeta";

type PetStatus = Awaited<ReturnType<typeof GetPetStatus>>;

type PetStatusBarsProps = {
  open: boolean;
};

type StatusMetricKey = "water" | "hunger" | "health" | "mood" | "energy";

type StatusMetric = {
  key: StatusMetricKey;
  label: string;
  fillColor: string;
};

const MAX_STATUS_VALUE = 100;
const STATUS_POLL_INTERVAL_MS = 2000;

const statusMetrics: StatusMetric[] = [
  { key: "water", label: "Water", fillColor: "rgba(59, 130, 246, 0.9)" },
  { key: "hunger", label: "Hunger", fillColor: "rgba(245, 158, 11, 0.9)" },
  { key: "health", label: "Health", fillColor: "rgba(34, 197, 94, 0.9)" },
  { key: "mood", label: "Mood", fillColor: "rgba(168, 85, 247, 0.9)" },
  { key: "energy", label: "Energy", fillColor: "rgba(236, 72, 153, 0.9)" },
];

const clampStatusValue = (value: number) =>
  Math.max(0, Math.min(MAX_STATUS_VALUE, value));

const readStatusNumber = (
  status: PetStatus | null,
  key: StatusMetricKey | "experience" | "money"
) => {
  const direct = (status as unknown as Record<string, unknown> | null)?.[key];
  if (typeof direct === "number") {
    return direct;
  }
  const nested = (
    status as unknown as { attributes?: Record<string, unknown> } | null
  )?.attributes?.[key];
  if (typeof nested === "number") {
    return nested;
  }
  return null;
};

export function PetStatusBars({ open }: PetStatusBarsProps) {
  const [status, setStatus] = useState<PetStatus | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!open) {
      return;
    }

    let cancelled = false;

    const fetchStatus = async () => {
      try {
        const snapshot = await GetPetStatus();
        if (cancelled) {
          return;
        }
        setStatus(snapshot);
        setError(null);
      } catch (err) {
        if (cancelled) {
          return;
        }
        console.error("GetPetStatus failed", err);
        setError("Status unavailable");
      }
    };

    void fetchStatus();
    const intervalId = window.setInterval(() => {
      void fetchStatus();
    }, STATUS_POLL_INTERVAL_MS);

    return () => {
      cancelled = true;
      window.clearInterval(intervalId);
    };
  }, [open]);

  const metaText = useMemo(() => {
    const experience = readStatusNumber(status, "experience");
    const money = readStatusNumber(status, "money");
    if (experience === null && money === null) {
      return null;
    }
    const moneyText =
      typeof money === "number" ? money.toLocaleString() : String(money ?? "—");
    const expText =
      typeof experience === "number" ? String(experience) : String("—");
    return `Exp ${expText} · $${moneyText}`;
  }, [status]);

  return (
    <div
      className={`pet-status${open ? " pet-status--open" : ""}`}
      aria-hidden={!open}
    >
      <div className="pet-status__panel">
        <div className="pet-status__header">
          <span className="pet-status__title">Status</span>
          {metaText ? (
            <span className="pet-status__meta" title={metaText}>
              {metaText}
            </span>
          ) : null}
        </div>

        <div className="pet-status__rows">
          {statusMetrics.map((metric) => {
            const rawValue = readStatusNumber(status, metric.key);
            const displayValue =
              typeof rawValue === "number" ? clampStatusValue(rawValue) : null;
            const percent =
              displayValue !== null
                ? (displayValue / MAX_STATUS_VALUE) * 100
                : 0;

            return (
              <div key={metric.key} className="pet-status-row">
                <div className="pet-status-row__label">
                  <span>{metric.label}</span>
                  <span className="pet-status-row__value">
                    {displayValue !== null
                      ? `${displayValue}/${MAX_STATUS_VALUE}`
                      : "—"}
                  </span>
                </div>
                <div
                  className="pet-status-row__bar"
                  role="progressbar"
                  aria-label={metric.label}
                  aria-valuenow={displayValue ?? 0}
                  aria-valuemin={0}
                  aria-valuemax={MAX_STATUS_VALUE}
                >
                  <div
                    className="pet-status-row__fill"
                    style={{
                      width: `${percent}%`,
                      backgroundColor: metric.fillColor,
                    }}
                  />
                </div>
              </div>
            );
          })}
        </div>

        {error ? <div className="pet-status__error">{error}</div> : null}
      </div>
    </div>
  );
}
