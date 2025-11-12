import {
  CSSProperties,
  MouseEvent as ReactMouseEvent,
  useCallback,
  useEffect,
  useLayoutEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { Menu, MenuItem } from "@mui/material";
import { Quit as QuitFromApp } from "../../wailsjs/go/app/App";
import {
  AdjustWindowfromLeftBottom,
  UpdatePetHitbox,
} from "../../wailsjs/go/window/WindowService";
import {
  WindowGetPosition,
  WindowSetPosition,
} from "../../wailsjs/runtime/runtime";

const baseContainerStyle: CSSProperties = {
  position: "relative",
  display: "flex",
  alignItems: "flex-start",
  justifyContent: "flex-start",
  padding: "16px 26px 32px",
  pointerEvents: "auto",
  transform: "translateZ(0)",
  willChange: "transform",
  backfaceVisibility: "hidden",
  transition: "none",
  userSelect: "none",
  background: "transparent",
  backgroundColor: "transparent",
};

const baseSpriteStyle: CSSProperties = {
  width: "220px",
  maxWidth: "220px",
  minWidth: "140px",
  maxHeight: "220px",
  pointerEvents: "none",
  transformOrigin: "center bottom",
  display: "block",
  background: "transparent",
  backgroundColor: "transparent",
  animation: "none",
  transform: "none",
  willChange: "auto",
  backfaceVisibility: "hidden",
  transition: "none",
};

const DEFAULT_WINDOW_SIZE = { width: 280, height: 280 };
const FALLBACK_MENU_ANCHOR = { left: 32, top: 28 };
const WINDOW_EXPAND_TARGET_SIZE = { width: 500, height: 320 }; // reserve space for side menu without clipping the pet
const MENU_OFFSET_X = 16;

type PetProps = {
  sprite: string;
  isFloating: boolean;
  onToggleFloating: () => void;
};

export function Pet({ sprite, isFloating, onToggleFloating }: PetProps) {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [menuAnchorPosition, setMenuAnchorPosition] = useState<{
    top: number;
    left: number;
  } | null>(null);
  const [isGrabbing, setIsGrabbing] = useState(false);
  const [isDragging, setIsDragging] = useState(false);

  const spriteRef = useRef<HTMLImageElement | null>(null);
  const spriteSizeRef = useRef<{ width: number; height: number }>({
    ...DEFAULT_WINDOW_SIZE,
  });
  const dragStateRef = useRef<{
    startScreenX: number;
    startScreenY: number;
    windowX: number;
    windowY: number;
    ready: boolean;
  } | null>(null);
  const draggingRef = useRef(false);
  const menuWindowExpandedRef = useRef(false);

  const syncHitbox = useCallback(() => {
    if (typeof window === "undefined") {
      return;
    }
    const node = spriteRef.current;
    if (!node) {
      return;
    }
    const rect = node.getBoundingClientRect();
    spriteSizeRef.current = { width: rect.width, height: rect.height };
    void UpdatePetHitbox({
      left: rect.left,
      top: rect.top,
      width: rect.width,
      height: rect.height,
      devicePixelRatio: window.devicePixelRatio || 1,
    }).catch((error: unknown) =>
      console.error("UpdatePetHitbox failed", error)
    );
  }, []);

  useLayoutEffect(() => {
    if (!sprite) {
      return;
    }
    if (typeof window === "undefined") {
      return;
    }
    const frame = window.requestAnimationFrame(() => {
      syncHitbox();
    });
    return () => window.cancelAnimationFrame(frame);
  }, [sprite, syncHitbox]);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    const handleResize = () => {
      syncHitbox();
    };
    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, [syncHitbox]);

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    const handleMouseMove = (event: MouseEvent) => {
      if (!draggingRef.current) {
        return;
      }
      event.preventDefault();
      const dragState = dragStateRef.current;
      if (!dragState || !dragState.ready) {
        return;
      }
      const deltaX = event.screenX - dragState.startScreenX;
      const deltaY = event.screenY - dragState.startScreenY;
      WindowSetPosition(
        Math.round(dragState.windowX + deltaX),
        Math.round(dragState.windowY + deltaY)
      );
    };

    const handleMouseUp = () => {
      if (!draggingRef.current) {
        return;
      }
      draggingRef.current = false;
      dragStateRef.current = null;
      setIsDragging(false);
      setIsGrabbing(false);
      syncHitbox();
    };

    window.addEventListener("mousemove", handleMouseMove);
    window.addEventListener("mouseup", handleMouseUp);

    return () => {
      window.removeEventListener("mousemove", handleMouseMove);
      window.removeEventListener("mouseup", handleMouseUp);
    };
  }, [syncHitbox]);

  const containerStyle = useMemo<CSSProperties>(
    () => ({
      ...baseContainerStyle,
      cursor: isGrabbing ? "grabbing" : "grab",
    }),
    [isGrabbing]
  );

  type CSSVars = { [key: `--${string}`]: string };

  const spriteStyle = useMemo(() => {
    const style: React.CSSProperties & CSSVars = {
      ...baseSpriteStyle,

      imageRendering: "pixelated",

      // vendor-prefixed props that aren't in the type: use localized "any" key cast
      ["WebkitImageRendering" as any]: "pixelated",
      ["MozImageRendering" as any]: "pixelated",
      ["msInterpolationMode" as any]: "nearest-neighbor",

      ...(isFloating ? {} : { animationPlayState: "paused" }),
    };

    return style;
  }, [isFloating]);

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
    setMenuAnchorPosition(computeMenuAnchorPosition());
    await AdjustWindowfromLeftBottom(
      WINDOW_EXPAND_TARGET_SIZE.width,
      WINDOW_EXPAND_TARGET_SIZE.height
    );
    setIsMenuOpen(true);
  };

  const handleCloseMenu = async () => {
    setIsMenuOpen(false);
    await AdjustWindowfromLeftBottom(
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

  const handleMouseDown = (event: ReactMouseEvent<HTMLDivElement>) => {
    if (event.button !== 0) {
      event.stopPropagation();
      return;
    }

    event.preventDefault();
    handleCloseMenu();
    draggingRef.current = true;
    setIsGrabbing(true);
    setIsDragging(true);
    dragStateRef.current = {
      startScreenX: event.screenX,
      startScreenY: event.screenY,
      windowX: 0,
      windowY: 0,
      ready: false,
    };

    void WindowGetPosition()
      .then((pos) => {
        if (!draggingRef.current || !dragStateRef.current) {
          return;
        }
        dragStateRef.current = {
          ...dragStateRef.current,
          windowX: pos.x,
          windowY: pos.y,
          ready: true,
        };
      })
      .catch((error: unknown) => {
        console.error("WindowGetPosition failed", error);
        dragStateRef.current = null;
        draggingRef.current = false;
        setIsDragging(false);
        setIsGrabbing(false);
      });
  };

  useEffect(() => {
    if (isDragging) {
      return;
    }
    syncHitbox();
  }, [isDragging, syncHitbox]);

  return (
    <div
      style={containerStyle}
      onDoubleClick={onToggleFloating}
      onContextMenu={handleContextMenu}
      onMouseDown={handleMouseDown}
    >
      {sprite && (
        <img
          src={sprite}
          alt="Desktop pet"
          draggable={false}
          style={spriteStyle}
          ref={spriteRef}
          onLoad={syncHitbox}
        />
      )}

      <Menu
        open={isMenuOpen}
        onClose={handleCloseMenu}
        anchorReference="anchorPosition"
        anchorPosition={menuAnchorPosition ?? undefined}
        transformOrigin={{ vertical: "bottom", horizontal: "left" }}
        disablePortal
        sx={{
          position: "fixed",
        }}
        slotProps={{
          paper: {
            elevation: 6,
            sx: {
              "--wails-draggable": "no-drag",
              minWidth: 140,
              borderRadius: 1.5,
              bgcolor: "rgba(30, 30, 30, 0.95)",
              color: "#f5f5f5",
              py: 0.5,
              backdropFilter: "blur(6px)",
            },
          },
          list: {
            dense: true,
            onMouseDown: (event: ReactMouseEvent) => {
              event.stopPropagation();
            },
            sx: {
              "--wails-draggable": "no-drag",
              py: 0,
            },
          },
        }}
      >
        <MenuItem
          onClick={handleQuit}
          sx={{
            "--wails-draggable": "no-drag",
            fontWeight: 500,
            px: 2,
            py: 1.25,
            "&:hover": {
              bgcolor: "rgba(255,255,255,0.12)",
            },
          }}
        >
          Quit
        </MenuItem>
      </Menu>
    </div>
  );
}
