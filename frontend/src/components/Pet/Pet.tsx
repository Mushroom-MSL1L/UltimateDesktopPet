import { MouseEvent as ReactMouseEvent, useCallback, useRef } from "react";
import { Quit as QuitFromApp } from "../../../wailsjs/go/app/App";
import "./Pet.css";
import { TalkBar } from "./TalkBar";
import { ResponseBubble } from "./ResponseBubble";

const ANCHOR_OFFSET_X = 16;

type AnchorPosition = { left: number; top: number };

type PetProps = {
  sprite: string;
  onOpenDialog: (anchor: AnchorPosition) => void;
  isQuickTalkOpen: boolean;
  onSendQuickMessage: (message: string) => void | Promise<void>;
  onRequestQuickTalk?: () => void;
  onCloseQuickTalk?: () => void;
  isChatBusy?: boolean;
  quickResponseMessage: string | null;
  isResponseBubbleOpen: boolean;
  onDismissResponseBubble: () => void;
  onSpriteMoveStart?: () => void;
  onSpriteMoveEnd?: () => void;
};

export function Pet({
  sprite,
  onOpenDialog,
  isQuickTalkOpen,
  onSendQuickMessage,
  onRequestQuickTalk,
  onCloseQuickTalk,
  isChatBusy,
  quickResponseMessage,
  isResponseBubbleOpen,
  onDismissResponseBubble,
  onSpriteMoveStart,
  onSpriteMoveEnd,
}: PetProps) {
  const spriteRef = useRef<HTMLImageElement | null>(null);

  const computeAnchorPosition = useCallback(() => {
    const spriteNode = spriteRef.current;
    if (!spriteNode) {
      return { left: ANCHOR_OFFSET_X, top: 0 };
    }
    const rect = spriteNode.getBoundingClientRect();
    return {
      left: rect.right + ANCHOR_OFFSET_X,
      top: rect.top,
    };
  }, []);

  const handleContextMenu = (event: ReactMouseEvent<HTMLDivElement>) => {
    event.preventDefault();
    event.stopPropagation();
    onRequestQuickTalk?.();
  };

  const handleQuit = () => {
    void QuitFromApp();
    console.log("Quit called from App");
  };

  const handleOpenDialog = () => {
    const anchor = computeAnchorPosition();
    onOpenDialog(anchor);
  };

  const handleShellMouseDown = (event: ReactMouseEvent<HTMLDivElement>) => {
    if (event.button !== 0) {
      return;
    }
    const target = event.target;
    if (!(target instanceof Element)) {
      return;
    }
    onSpriteMoveStart?.();

    const clickedOpenElement =
      target.closest(".pet-quick-panel") || target.closest(".response-bubble");

    if (!clickedOpenElement) {
      onDismissResponseBubble();
      onCloseQuickTalk?.();
    }
  };

  const handleShellMouseUp = () => {
    onSpriteMoveEnd?.();
  };

  const handleShellMouseLeave = () => {
    onSpriteMoveEnd?.();
  };

  return (
    <div
      className="pet-shell"
      onContextMenu={handleContextMenu}
      onMouseDown={handleShellMouseDown}
      onMouseUp={handleShellMouseUp}
      onMouseLeave={handleShellMouseLeave}
    >
      {sprite && (
        <img
          src={sprite}
          alt="Desktop pet"
          draggable={false}
          className="pet-sprite"
          ref={spriteRef}
        />
      )}
      <ResponseBubble
        message={quickResponseMessage ?? null}
        open={isResponseBubbleOpen}
        onDismiss={onDismissResponseBubble}
      />
      <div
        className={`pet-quick-panel${
          isQuickTalkOpen ? " pet-quick-panel--open" : ""
        }`}
        aria-hidden={!isQuickTalkOpen}
      >
        <div className="pet-talk-actions">
          <button
            type="button"
            className="pet-talk-actions__button"
            onClick={handleOpenDialog}
          >
            Open dialog
          </button>
          <button
            type="button"
            className="pet-talk-actions__button pet-talk-actions__button--danger"
            onClick={handleQuit}
          >
            Quit
          </button>
        </div>
        <TalkBar
          open={isQuickTalkOpen}
          onSend={onSendQuickMessage}
          disabled={isChatBusy}
        />
      </div>
    </div>
  );
}
