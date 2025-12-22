import { MouseEvent as ReactMouseEvent, useCallback, useRef } from "react";
import { Quit as QuitFromApp } from "../../../wailsjs/go/app/App";
import "./Pet.css";
import { TalkBar } from "./TalkBar";
import { ResponseBubble } from "./ResponseBubble";
import { PetStatusBars } from "./PetStatusBars";

const ANCHOR_OFFSET_X = 16;

type AnchorPosition = { left: number; top: number };

type ItemEffect = { key: number; src: string; alt: string };

type PetProps = {
  sprite: string;
  onOpenDialog: (anchor: AnchorPosition) => void;
  onOpenShop: (anchor: AnchorPosition) => void;
  onOpenConfig: (anchor: AnchorPosition) => void;
  itemEffect?: ItemEffect | null;
  isQuickTalkOpen: boolean;
  onSendQuickMessage: (message: string) => void | Promise<void>;
  onRequestQuickTalk?: () => void;
  onCloseQuickTalk?: () => void;
  isChatBusy?: boolean;
  isGeminiKeyMissing?: boolean;
  quickResponseMessage: string | null;
  isResponseBubbleOpen: boolean;
  onDismissResponseBubble: () => void;
  onSpriteMoveStart?: () => void;
  onSpriteMoveEnd?: () => void;
  onUserInteraction?: () => void;
};

export function Pet({
  sprite,
  onOpenDialog,
  onOpenShop,
  onOpenConfig,
  itemEffect,
  isQuickTalkOpen,
  onSendQuickMessage,
  onRequestQuickTalk,
  onCloseQuickTalk,
  isChatBusy,
  isGeminiKeyMissing = false,
  quickResponseMessage,
  isResponseBubbleOpen,
  onDismissResponseBubble,
  onSpriteMoveStart,
  onSpriteMoveEnd,
  onUserInteraction,
}: PetProps) {
  const spriteRef = useRef<HTMLImageElement | null>(null);

  const notifyInteraction = useCallback(() => {
    onUserInteraction?.();
  }, [onUserInteraction]);

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
    notifyInteraction();
    onRequestQuickTalk?.();
  };

  const handleQuit = () => {
    notifyInteraction();
    void QuitFromApp();
    console.log("Quit called from App");
  };

  const handleOpenDialog = () => {
    notifyInteraction();
    const anchor = computeAnchorPosition();
    onOpenDialog(anchor);
  };

  const handleOpenShop = () => {
    notifyInteraction();
    const anchor = computeAnchorPosition();
    onOpenShop(anchor);
  };

  const handleOpenConfig = () => {
    notifyInteraction();
    const anchor = computeAnchorPosition();
    onOpenConfig(anchor);
  };

  const handleShellMouseDown = (event: ReactMouseEvent<HTMLDivElement>) => {
    if (event.button !== 0) {
      return;
    }
    notifyInteraction();
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
    notifyInteraction();
    onSpriteMoveEnd?.();
  };

  const handleShellMouseLeave = () => {
    notifyInteraction();
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
      {itemEffect ? (
        <img
          key={itemEffect.key}
          src={itemEffect.src}
          alt={itemEffect.alt}
          draggable={false}
          className="pet-item-effect"
        />
      ) : null}
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
            className="pet-talk-actions__button"
            onClick={handleOpenShop}
          >
            Shop
          </button>
          <button
            type="button"
            className="pet-talk-actions__button"
            onClick={handleOpenConfig}
          >
            Settings
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
          disabled={isChatBusy || isGeminiKeyMissing}
          missingGeminiKey={isGeminiKeyMissing}
          onOpenSettings={handleOpenConfig}
        />
        <PetStatusBars open={isQuickTalkOpen} />
      </div>
    </div>
  );
}
