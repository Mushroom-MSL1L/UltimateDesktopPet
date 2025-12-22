import { FormEvent, useEffect, useMemo, useRef, useState } from "react";

type TalkBarProps = {
  open: boolean;
  onSend: (message: string) => void | Promise<void>;
  disabled?: boolean;
  placeholder?: string;
  missingGeminiKey?: boolean;
  onOpenSettings?: () => void;
};

export function TalkBar({
  open,
  onSend,
  disabled = false,
  missingGeminiKey = false,
  onOpenSettings,
  placeholder = "Talk to your pet...",
}: TalkBarProps) {
  const [draft, setDraft] = useState("");
  const inputRef = useRef<HTMLInputElement | null>(null);

  useEffect(() => {
    if (!open) {
      setDraft("");
      return;
    }
    const timeoutId = window.setTimeout(() => {
      inputRef.current?.focus();
    }, 0);
    return () => window.clearTimeout(timeoutId);
  }, [open]);

  const isSendDisabled = useMemo(
    () => disabled || missingGeminiKey || !draft.trim(),
    [disabled, draft, missingGeminiKey]
  );

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (isSendDisabled) {
      return;
    }
    const trimmed = draft.trim();
    if (!trimmed) {
      return;
    }
    try {
      await onSend(trimmed);
      setDraft("");
    } catch (error) {
      console.error("Quick talk failed", error);
    }
  };

  return (
    <div
      className={`talk-bar${open ? " talk-bar--open" : ""}`}
      aria-hidden={!open}
    >
      {missingGeminiKey ? (
        <div
          className="talk-bar__notice"
          style={{
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
            gap: 8,
            background: "rgba(0,0,0,0.35)",
            color: "#f5f5f5",
            padding: "6px 10px",
            borderRadius: 10,
            marginBottom: 6,
          }}
        >
          <span style={{ fontSize: 12 }}>
            Add a Gemini API key in Settings to enable chat.
          </span>
          {onOpenSettings ? (
            <button
              type="button"
              onClick={onOpenSettings}
              className="talk-bar__notice-button"
              style={{
                background: "#60a5fa",
                color: "#0b1221",
                border: "none",
                borderRadius: 6,
                padding: "4px 10px",
                cursor: "pointer",
                fontSize: 12,
                fontWeight: 600,
              }}
            >
              Open settings
            </button>
          ) : null}
        </div>
      ) : null}

      <form className="talk-bar__form" onSubmit={handleSubmit}>
        <input
          ref={inputRef}
          type="text"
          className="talk-bar__input"
          value={draft}
          placeholder={placeholder}
          onChange={(event) => setDraft(event.target.value)}
          disabled={disabled || missingGeminiKey}
        />
        <div className="talk-bar__actions">
          <button
            type="submit"
            className="talk-bar__send-button"
            disabled={isSendDisabled}
          >
            Send
          </button>
        </div>
      </form>
    </div>
  );
}
