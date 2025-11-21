import {
  CSSProperties,
  MouseEvent as ReactMouseEvent,
  useCallback,
  useLayoutEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { Quit as QuitFromApp } from "../../../wailsjs/go/app/App";
import { AdjustWindowFromBottom } from "../../../wailsjs/go/window/WindowService";
import "./Pet.css";
import { TalkBar } from "./TalkBar";
import { ResponseBubble } from "./ResponseBubble";
import {
  WindowGetPosition,
  WindowSetPosition,
} from "../../../wailsjs/runtime/runtime";

const DEFAULT_WINDOW_SIZE = { width: 150, height: 150 };
const WINDOW_EXPAND_TARGET_SIZE = { width: 400, height: 300 }; // reserve space for side menu without clipping the pet
const MENU_OFFSET_X = 16;

type AnchorPosition = { left: number; top: number };

type PetProps = {
  sprite: string;
  onOpenDialog: (anchor: AnchorPosition) => void;
  isQuickTalkOpen: boolean;
  onSendQuickMessage: (message: string) => void | Promise<void>;
  onRequestQuickTalk?: () => void;
  isChatBusy?: boolean;
  quickResponseMessage: string | null;
  isResponseBubbleOpen: boolean;
  onDismissResponseBubble: () => void;
};

export function Pet({
  sprite,
  onOpenDialog,
  isQuickTalkOpen,
  onSendQuickMessage,
  onRequestQuickTalk,
  isChatBusy,
  quickResponseMessage,
  isResponseBubbleOpen,
  onDismissResponseBubble,
}: PetProps) {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [menuAnchorPosition, setMenuAnchorPosition] = useState<{
    top: number;
    left: number;
  } | null>(null);

  const spriteRef = useRef<HTMLImageElement | null>(null);
  const menuWindowExpandedRef = useRef(false);
  const menuPaperRef = useRef<HTMLDivElement | null>(null);

  const computeMenuAnchorPosition = useCallback(() => {
    const spriteNode = spriteRef.current;
    const rect = spriteNode!.getBoundingClientRect();
    return {
      left: rect.right + MENU_OFFSET_X,
      top: rect.top,
    };
  }, []);

  const handleContextMenu = async (event: ReactMouseEvent<HTMLDivElement>) => {
    event.preventDefault();
    event.stopPropagation();
    onRequestQuickTalk?.();
    setMenuAnchorPosition(computeMenuAnchorPosition());
    await AdjustWindowFromBottom(
      WINDOW_EXPAND_TARGET_SIZE.width,
      WINDOW_EXPAND_TARGET_SIZE.height
    );
    setIsMenuOpen(true);
    menuWindowExpandedRef.current = true;
  };

  const handleCloseMenu = async () => {
    setIsMenuOpen(false);
    await AdjustWindowFromBottom(
      DEFAULT_WINDOW_SIZE.width,
      DEFAULT_WINDOW_SIZE.height
    );
    menuWindowExpandedRef.current = false;
  };

  const handleQuit = () => {
    handleCloseMenu();
    void QuitFromApp();
    console.log("Quit called from App");
  };

  const handleOpenDialogFromMenu = () => {
    const anchor = computeMenuAnchorPosition();
    handleCloseMenu();
    onOpenDialog(anchor);
  };

  const handleMouseDown = (event: ReactMouseEvent<HTMLDivElement>) => {
    if (event.button !== 0) {
      event.stopPropagation();
      return;
    }

    if (isMenuOpen) {
      const menuNode = menuPaperRef.current;
      if (menuNode && menuNode.contains(event.target as Node)) {
        // Let menu items process this click.
        return;
      }
      onDismissResponseBubble!();
      const target = event.target;
      const clickedTalkBar =
        target instanceof Element && Boolean(target.closest(".talk-bar"));
      if (clickedTalkBar) {
        return;
      }
      event.preventDefault();
      void handleCloseMenu();
      return;
    }
  };

  return (
    <div
      className="pet-shell"
      onContextMenu={handleContextMenu}
      onMouseDown={handleMouseDown}
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
      <TalkBar
        open={isQuickTalkOpen}
        onSend={onSendQuickMessage}
        disabled={isChatBusy}
      />
    </div>
  );
}
