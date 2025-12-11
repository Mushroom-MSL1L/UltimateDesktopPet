import { FormEvent, useEffect, useMemo, useRef, useState } from "react";

type TalkBarProps = {
  open: boolean;
  onSend: (message: string) => void | Promise<void>;
  disabled?: boolean;
  placeholder?: string;
};

export function TalkBar({
  open,
  onSend,
  disabled = false,
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
    () => disabled || !draft.trim(),
    [disabled, draft]
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
      <form className="talk-bar__form" onSubmit={handleSubmit}>
        <input
          ref={inputRef}
          type="text"
          className="talk-bar__input"
          value={draft}
          placeholder={placeholder}
          onChange={(event) => setDraft(event.target.value)}
          disabled={disabled}
        />
        <div className="talk-bar__actions">
          <button
            type="submit"
            className="talk-bar__send-button"
            disabled={isSendDisabled}
          >
            {disabled ? "..." : "Send"}
          </button>
        </div>
      </form>
    </div>
  );
}
