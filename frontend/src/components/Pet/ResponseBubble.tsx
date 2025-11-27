import { useEffect, useState } from "react";

const FADE_OUT_DURATION_MS = 180;

type ResponseBubbleProps = {
  /**
   * Text that should appear inside the bubble.
   */
  message: string | null;
  /**
   * Whether the bubble should be rendered and visible.
   */
  open: boolean;
  /**
   * Time in milliseconds before the bubble fades away automatically.
   */
  duration?: number;
  /**
   * Called when the bubble auto-dismisses.
   */
  onDismiss?: () => void;
};

export function ResponseBubble({
  message,
  open,
  duration = 40000,
  onDismiss,
}: ResponseBubbleProps) {
  const [renderedMessage, setRenderedMessage] = useState<string | null>(() =>
    open && message ? message : null
  );

  useEffect(() => {
    if (open && message) {
      setRenderedMessage(message);
    }
  }, [open, message]);

  useEffect(() => {
    if (open || !renderedMessage) {
      return;
    }
    const timeoutId = window.setTimeout(() => {
      setRenderedMessage(null);
    }, FADE_OUT_DURATION_MS);
    return () => window.clearTimeout(timeoutId);
  }, [open, renderedMessage]);

  useEffect(() => {
    if (!open || !renderedMessage) {
      return;
    }
    const timeoutId = window.setTimeout(() => {
      onDismiss?.();
    }, duration);
    return () => window.clearTimeout(timeoutId);
  }, [open, renderedMessage, duration, onDismiss]);

  if (!renderedMessage) {
    return null;
  }

  return (
    <div
      className={`response-bubble${open ? " response-bubble--open" : ""}`}
      role="status"
      aria-live="polite"
    >
      <span>{renderedMessage}</span>
    </div>
  );
}
