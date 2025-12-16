import {
  CSSProperties,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { ChatWithPet } from "../wailsjs/go/chat/ChatMeta";
import { LoadSystemConfig } from "../wailsjs/go/configs/ConfigService";
import { LoadItemFrameByID } from "../wailsjs/go/items/ItemsMeta";
import { UseItemByID } from "../wailsjs/go/pet/PetMeta";
import {
  AdjustWindowFromBottom,
  AdjustWindowFromLeftBottom,
  SetDevicePixelRatio,
} from "../wailsjs/go/window/WindowService";
import { WindowSetAlwaysOnTop } from "../wailsjs/runtime/runtime";
import { Pet } from "./components/Pet/Pet";
import {
  PetDialog,
  type ConversationMessage,
} from "./components/Pet/PetDialog";
import { PetShopDialog, type ShopItem } from "./components/Pet/PetShopDialog";
import { SystemConfigDialog } from "./components/SystemConfigDialog";
import {
  SpriteAnimationKey,
  type SpriteFramesCache,
  ensureAnimationFrames,
  preloadAnimations,
} from "./utils/spriteFrames";
import {
  WanderDirection,
  randomWanderDirection,
  wanderAnimationByDirection,
  wanderDirectionVectors,
} from "./utils/wander";
import { moveWindowWithinScreen } from "./utils/windowPosition";

const TRANSPARENT_BACKGROUND = "rgba(0, 0, 0, 0)";
const DEV_TINT_BACKGROUND = "rgba(20, 20, 20, 0.33)";

const baseAppShellStyle: CSSProperties = {
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  justifyContent: "flex-start",
  userSelect: "none",
  pointerEvents: "none",
  backdropFilter: "blur(6px)",
  position: "fixed",
  inset: 0,
  transform: "translateZ(0)",
  background: TRANSPARENT_BACKGROUND,
  backfaceVisibility: "hidden",
};

const PET_WINDOW_DEFAULT_SIZE = { width: 150, height: 150 };
const QUICK_TALK_WINDOW_SIZE = { width: 320, height: 550 };
const DIALOG_WINDOW_SIZE = { width: 900, height: 500 };
const SHOP_WINDOW_SIZE = { width: 900, height: 500 };
const CONFIG_WINDOW_SIZE = { width: 1200, height: 640 };
const SPRITE_FRAME_DURATION_MS = 150;
const ITEM_EFFECT_DURATION_MS = 900;
const IDLE_WANDER_DELAY_MS = 5000;
const WANDER_STEP_DISTANCE_PX = 8;
const WANDER_STEP_INTERVAL_MS = 120;
const WANDER_DIRECTION_CHANGE_MS = 2600;
const INTERACTION_DEBOUNCE_MS = 120;

const getErrorMessage = (error: unknown) => {
  if (error instanceof Error) {
    return error.message;
  }
  if (typeof error === "string") {
    return error;
  }
  return "Something went wrong.";
};

const isGeminiKeyError = (message: string) =>
  message.toLowerCase().includes("gemini api key missing");

function App() {
  const [sprite, setSprite] = useState<string>("");
  const [petFrames, setPetFrames] = useState<string[]>([]);
  const [isSpriteMoving, setIsSpriteMoving] = useState(false);
  const [isWandering, setIsWandering] = useState(false);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isShopOpen, setIsShopOpen] = useState(false);
  const [isConfigOpen, setIsConfigOpen] = useState(false);
  const [isSendingMessage, setIsSendingMessage] = useState(false);
  const [isGeminiKeyMissing, setIsGeminiKeyMissing] = useState(false);
  const [conversation, setConversation] = useState<ConversationMessage[]>([
    {
      id: "intro",
      role: "pet",
      text: "Hey there! I'm still learning to chat, but I'd love to hear from you.",
    },
  ]);
  const [dialogAnchor, setDialogAnchor] = useState<{
    left: number;
    top: number;
  } | null>(null);
  const [shopAnchor, setShopAnchor] = useState<{
    left: number;
    top: number;
  } | null>(null);
  const [configAnchor, setConfigAnchor] = useState<{
    left: number;
    top: number;
  } | null>(null);
  const [isQuickTalkOpen, setIsQuickTalkOpen] = useState(false);
  const [quickResponseMessage, setQuickResponseMessage] = useState<
    string | null
  >(null);
  const [isResponseBubbleOpen, setIsResponseBubbleOpen] = useState(false);
  const [itemEffect, setItemEffect] = useState<{
    key: number;
    src: string;
    alt: string;
  } | null>(null);
  const framesCacheRef = useRef<SpriteFramesCache>({});
  const isMountedRef = useRef(true);
  const idleTimerRef = useRef<number | null>(null);
  const wanderStepTimerRef = useRef<number | null>(null);
  const wanderDirectionTimerRef = useRef<number | null>(null);
  const itemUseTimerRef = useRef<number | null>(null);
  const itemUseRequestIdRef = useRef(0);
  const wanderDirectionRef = useRef<WanderDirection>("left");
  const isWanderingRef = useRef(false);
  const lastInteractionRef = useRef<number>(Date.now());
  const interactionPausedRef = useRef(false);
  const animationRequestIdRef = useRef(0);
  const scheduleIdleTimerRef = useRef<() => void>(() => {});
  const startWanderingRef = useRef<() => void>(() => {});
  const wasConfigOpenRef = useRef(false);

  const isDevMode = import.meta.env.DEV;
  const windowBackground = isDevMode
    ? DEV_TINT_BACKGROUND
    : TRANSPARENT_BACKGROUND;

  const appShellStyle = useMemo<CSSProperties>(() => {
    const alignment =
      isDialogOpen || isShopOpen || isConfigOpen ? "flex-start" : "center";
    return {
      ...baseAppShellStyle,
      background: windowBackground,
      alignItems: alignment,
    };
  }, [windowBackground, isConfigOpen, isDialogOpen, isShopOpen]);

  const refreshGeminiKeyStatus = useCallback(async () => {
    try {
      const cfg = await LoadSystemConfig();
      setIsGeminiKeyMissing(!cfg.geminiAPIKey?.trim());
    } catch (error) {
      console.warn("LoadSystemConfig failed while checking Gemini key", error);
    }
  }, []);

  useEffect(() => {
    void refreshGeminiKeyStatus();
  }, [refreshGeminiKeyStatus]);

  useEffect(() => {
    if (wasConfigOpenRef.current && !isConfigOpen) {
      void refreshGeminiKeyStatus();
    }
    wasConfigOpenRef.current = isConfigOpen;
  }, [isConfigOpen, refreshGeminiKeyStatus]);

  const clearIdleTimer = useCallback(() => {
    if (idleTimerRef.current !== null) {
      window.clearTimeout(idleTimerRef.current);
      idleTimerRef.current = null;
    }
  }, []);

  const clearWanderTimers = useCallback(() => {
    if (wanderStepTimerRef.current !== null) {
      window.clearInterval(wanderStepTimerRef.current);
      wanderStepTimerRef.current = null;
    }
    if (wanderDirectionTimerRef.current !== null) {
      window.clearTimeout(wanderDirectionTimerRef.current);
      wanderDirectionTimerRef.current = null;
    }
  }, []);

  const clearItemUseTimer = useCallback(() => {
    if (itemUseTimerRef.current !== null) {
      window.clearTimeout(itemUseTimerRef.current);
      itemUseTimerRef.current = null;
    }
  }, []);

  const applyAnimationFrames = useCallback(
    async (animation: SpriteAnimationKey) => {
      const requestId = animationRequestIdRef.current + 1;
      animationRequestIdRef.current = requestId;

      const frames = await ensureAnimationFrames(
        animation,
        framesCacheRef.current
      );

      if (
        !isMountedRef.current ||
        requestId !== animationRequestIdRef.current
      ) {
        return;
      }

      setPetFrames((previous) => {
        if (frames.length > 0) {
          return frames;
        }
        const fallback =
          framesCacheRef.current.stand ??
          framesCacheRef.current.move_left ??
          previous;
        return fallback ?? previous;
      });
    },
    []
  );

  const setWanderDirection = useCallback(
    (direction: WanderDirection) => {
      wanderDirectionRef.current = direction;
      void applyAnimationFrames(wanderAnimationByDirection[direction]);
    },
    [applyAnimationFrames]
  );

  const stopWandering = useCallback(() => {
    if (!isWanderingRef.current) {
      return;
    }
    clearWanderTimers();
    isWanderingRef.current = false;
    setIsWandering(false);
    void applyAnimationFrames("stand");
  }, [applyAnimationFrames, clearWanderTimers]);

  const performWanderStep: () => Promise<void> = useCallback(async () => {
    if (!isWanderingRef.current) {
      return;
    }
    const direction = wanderDirectionRef.current;
    const delta = wanderDirectionVectors[direction];
    const moveResult = await moveWindowWithinScreen(
      delta.x * WANDER_STEP_DISTANCE_PX,
      delta.y * WANDER_STEP_DISTANCE_PX
    );
    if (!moveResult) {
      stopWandering();
      scheduleIdleTimerRef.current();
      return;
    }
    if (!moveResult.didMove) {
      setWanderDirection(randomWanderDirection(direction));
    }
  }, [setWanderDirection, stopWandering]);

  const scheduleDirectionChange: () => void = useCallback(() => {
    if (wanderDirectionTimerRef.current !== null) {
      window.clearTimeout(wanderDirectionTimerRef.current);
      wanderDirectionTimerRef.current = null;
    }
    wanderDirectionTimerRef.current = window.setTimeout(() => {
      const nextDirection = randomWanderDirection(
        wanderDirectionRef.current ?? "left"
      );
      setWanderDirection(nextDirection);
      scheduleDirectionChange();
    }, WANDER_DIRECTION_CHANGE_MS);
  }, [setWanderDirection]);

  const startWandering: () => void = useCallback(() => {
    if (interactionPausedRef.current || isWanderingRef.current) {
      return;
    }
    clearIdleTimer();
    clearWanderTimers();
    isWanderingRef.current = true;
    setIsWandering(true);
    const nextDirection = randomWanderDirection();
    setWanderDirection(nextDirection);

    wanderStepTimerRef.current = window.setInterval(() => {
      void performWanderStep();
    }, WANDER_STEP_INTERVAL_MS);

    scheduleDirectionChange();
  }, [
    clearIdleTimer,
    clearWanderTimers,
    performWanderStep,
    scheduleDirectionChange,
    setWanderDirection,
  ]);

  const scheduleIdleTimer: () => void = useCallback(() => {
    if (interactionPausedRef.current) {
      clearIdleTimer();
      return;
    }
    clearIdleTimer();
    idleTimerRef.current = window.setTimeout(() => {
      startWanderingRef.current();
    }, IDLE_WANDER_DELAY_MS);
  }, [clearIdleTimer]);

  startWanderingRef.current = startWandering;
  scheduleIdleTimerRef.current = scheduleIdleTimer;

  const registerInteraction = useCallback(() => {
    const now = Date.now();
    if (now - lastInteractionRef.current < INTERACTION_DEBOUNCE_MS) {
      return;
    }
    lastInteractionRef.current = now;
    stopWandering();
    scheduleIdleTimer();
  }, [scheduleIdleTimer, stopWandering]);

  useEffect(() => {
    isMountedRef.current = true;
    SetDevicePixelRatio(window.devicePixelRatio);

    const htmlElement = document.documentElement;
    const rootElement = document.getElementById("root");

    const previousBackgrounds = {
      bodyBackground: document.body.style.background,
      bodyBackgroundColor: document.body.style.backgroundColor,
      htmlBackground: htmlElement.style.background,
      htmlBackgroundColor: htmlElement.style.backgroundColor,
      rootBackground: rootElement?.style.background ?? "",
      rootBackgroundColor: rootElement?.style.backgroundColor ?? "",
    };

    const applyBackground = (
      element: HTMLElement | null,
      background: string
    ) => {
      if (!element) {
        return;
      }

      element.style.background = background;
      element.style.backgroundColor = background;
    };

    applyBackground(document.body, windowBackground);
    applyBackground(htmlElement, windowBackground);
    applyBackground(rootElement, windowBackground);

    WindowSetAlwaysOnTop(true);

    preloadAnimations([
      "stand",
      "move_left",
      "move_right",
      "move_far",
      "drag",
      "drop",
    ])
      .then((cache) => {
        if (!isMountedRef.current) {
          return;
        }
        framesCacheRef.current = cache;
        const hasStandFrames = (cache.stand?.length ?? 0) > 0;
        const initialAnimation: SpriteAnimationKey = hasStandFrames
          ? "stand"
          : "move_left";
        void applyAnimationFrames(initialAnimation);
      })
      .catch((err) => console.error("Pet frames failed to load", err));

    return () => {
      isMountedRef.current = false;

      document.body.style.background = previousBackgrounds.bodyBackground;
      document.body.style.backgroundColor =
        previousBackgrounds.bodyBackgroundColor;
      htmlElement.style.background = previousBackgrounds.htmlBackground;
      htmlElement.style.backgroundColor =
        previousBackgrounds.htmlBackgroundColor;
      if (rootElement) {
        rootElement.style.background = previousBackgrounds.rootBackground;
        rootElement.style.backgroundColor =
          previousBackgrounds.rootBackgroundColor;
      }
    };
  }, [applyAnimationFrames, windowBackground]);

  useEffect(() => {
    interactionPausedRef.current =
      isDialogOpen ||
      isShopOpen ||
      isConfigOpen ||
      isQuickTalkOpen ||
      isSpriteMoving;
    if (interactionPausedRef.current) {
      stopWandering();
      clearIdleTimer();
      return;
    }
    scheduleIdleTimer();
  }, [
    clearIdleTimer,
    isDialogOpen,
    isShopOpen,
    isQuickTalkOpen,
    isConfigOpen,
    isSpriteMoving,
    scheduleIdleTimer,
    stopWandering,
  ]);

  useEffect(() => {
    const handleUserAction = () => registerInteraction();
    const events: (keyof WindowEventMap)[] = [
      "pointerdown",
      "pointermove",
      "touchstart",
      "keydown",
      "wheel",
    ];
    events.forEach((eventName) =>
      window.addEventListener(eventName, handleUserAction, { passive: true })
    );

    return () => {
      events.forEach((eventName) =>
        window.removeEventListener(eventName, handleUserAction)
      );
    };
  }, [registerInteraction]);

  useEffect(() => {
    if (petFrames.length === 0) {
      setSprite("");
      return;
    }

    let frameIndex = 0;
    setSprite(petFrames[frameIndex]);

    if (petFrames.length === 1) {
      return;
    }

    const intervalId = window.setInterval(() => {
      frameIndex = (frameIndex + 1) % petFrames.length;
      setSprite(petFrames[frameIndex]);
    }, SPRITE_FRAME_DURATION_MS);

    return () => {
      window.clearInterval(intervalId);
    };
  }, [petFrames]);

  useEffect(() => {
    return () => {
      clearIdleTimer();
      clearWanderTimers();
      clearItemUseTimer();
    };
  }, [clearIdleTimer, clearItemUseTimer, clearWanderTimers]);

  const handleOpenDialog = useCallback(
    (anchor: { left: number; top: number }) => {
      registerInteraction();
      setDialogAnchor(anchor);
      setShopAnchor(null);
      setConfigAnchor(null);
      setIsShopOpen(false);
      setIsConfigOpen(false);
      setIsQuickTalkOpen(false);
      void AdjustWindowFromLeftBottom(
        DIALOG_WINDOW_SIZE.width,
        DIALOG_WINDOW_SIZE.height
      );
      setIsDialogOpen(true);
    },
    [registerInteraction]
  );
  const handleCloseDialog = useCallback(() => {
    registerInteraction();
    setIsDialogOpen(false);
    setDialogAnchor(null);
    const shouldReopenCard = !isSpriteMoving && !isWandering;
    setIsQuickTalkOpen(shouldReopenCard);
    const targetSize = shouldReopenCard
      ? QUICK_TALK_WINDOW_SIZE
      : PET_WINDOW_DEFAULT_SIZE;
    void AdjustWindowFromLeftBottom(targetSize.width, targetSize.height);
  }, [isSpriteMoving, isWandering, registerInteraction]);

  const handleOpenShop = useCallback(
    (anchor: { left: number; top: number }) => {
      registerInteraction();
      setShopAnchor(anchor);
      setDialogAnchor(null);
      setConfigAnchor(null);
      setIsDialogOpen(false);
      setIsConfigOpen(false);
      setIsQuickTalkOpen(false);
      void AdjustWindowFromLeftBottom(
        SHOP_WINDOW_SIZE.width,
        SHOP_WINDOW_SIZE.height
      );
      setIsShopOpen(true);
    },
    [registerInteraction]
  );

  const handleCloseShop = useCallback(() => {
    registerInteraction();
    setIsShopOpen(false);
    setShopAnchor(null);
    const shouldReopenCard = !isSpriteMoving && !isWandering;
    setIsQuickTalkOpen(shouldReopenCard);
    const targetSize = shouldReopenCard
      ? QUICK_TALK_WINDOW_SIZE
      : PET_WINDOW_DEFAULT_SIZE;
    void AdjustWindowFromLeftBottom(targetSize.width, targetSize.height);
  }, [isSpriteMoving, isWandering, registerInteraction]);

  const handleOpenConfig = useCallback(
    (anchor: { left: number; top: number }) => {
      registerInteraction();
      setConfigAnchor(anchor);
      setDialogAnchor(null);
      setShopAnchor(null);
      setIsDialogOpen(false);
      setIsShopOpen(false);
      setIsQuickTalkOpen(false);
      void AdjustWindowFromLeftBottom(
        CONFIG_WINDOW_SIZE.width,
        CONFIG_WINDOW_SIZE.height
      );
      setIsConfigOpen(true);
    },
    [registerInteraction]
  );

  const handleCloseConfig = useCallback(() => {
    registerInteraction();
    setIsConfigOpen(false);
    setConfigAnchor(null);
    const shouldReopenCard = !isSpriteMoving && !isWandering;
    setIsQuickTalkOpen(shouldReopenCard);
    const targetSize = shouldReopenCard
      ? QUICK_TALK_WINDOW_SIZE
      : PET_WINDOW_DEFAULT_SIZE;
    void AdjustWindowFromLeftBottom(targetSize.width, targetSize.height);
  }, [isSpriteMoving, isWandering, registerInteraction]);

  const openConfigFromDialog = useCallback(() => {
    const anchor = dialogAnchor ?? configAnchor ?? { left: 0, top: 0 };
    handleOpenConfig(anchor);
  }, [configAnchor, dialogAnchor, handleOpenConfig]);

  const playUseItemAnimation = useCallback(
    async (itemFrame: string | null, itemName?: string) => {
      const requestId = itemUseRequestIdRef.current + 1;
      itemUseRequestIdRef.current = requestId;

      clearItemUseTimer();

      if (itemFrame) {
        setItemEffect({
          key: requestId,
          src: itemFrame,
          alt: itemName ? `Using ${itemName}` : "Using item",
        });
      } else {
        setItemEffect(null);
      }

      const dropFrames = await ensureAnimationFrames(
        "drop",
        framesCacheRef.current
      );

      if (!isMountedRef.current || requestId !== itemUseRequestIdRef.current) {
        return;
      }

      const dragFrame = dropFrames[3];
      if (dragFrame) {
        setPetFrames([dragFrame]);
      } else {
        void applyAnimationFrames("stand");
      }

      const durationMs = ITEM_EFFECT_DURATION_MS;

      itemUseTimerRef.current = window.setTimeout(() => {
        if (
          !isMountedRef.current ||
          requestId !== itemUseRequestIdRef.current
        ) {
          return;
        }
        setItemEffect(null);
        void applyAnimationFrames("stand");
      }, durationMs);
    },
    [applyAnimationFrames, clearItemUseTimer]
  );

  const handlePurchaseItem = useCallback(
    async (item: ShopItem, itemFrame: string | null) => {
      registerInteraction();
      await UseItemByID(item.id);

      let frame = itemFrame;
      if (!frame) {
        try {
          frame = await LoadItemFrameByID(item.id);
        } catch (error) {
          console.warn(
            "LoadItemFrameByID failed after purchase",
            item.id,
            error
          );
          frame = null;
        }
      }

      await playUseItemAnimation(frame, item.name);
    },
    [playUseItemAnimation, registerInteraction]
  );

  const sendMessage = useCallback(
    async (message: string) => {
      registerInteraction();
      if (isSendingMessage) {
        return null;
      }
      const trimmed = message.trim();
      if (!trimmed) {
        return null;
      }

      const userMessage: ConversationMessage = {
        id: `user-${Date.now()}`,
        role: "user",
        text: trimmed,
      };

      setConversation((previous) => [...previous, userMessage]);
      setIsSendingMessage(true);

      try {
        const response = await ChatWithPet(trimmed);
        setIsGeminiKeyMissing(false);
        const replyText = response ?? "I'm thinking about what to say...";
        setConversation((previous) => [
          ...previous,
          {
            id: `pet-${Date.now()}`,
            role: "pet",
            text: replyText,
          },
        ]);
        return replyText;
      } catch (error) {
        console.error("ChatWithPet failed", error);
        const message = getErrorMessage(error);
        const missingKey = isGeminiKeyError(message);
        if (missingKey) {
          setIsGeminiKeyMissing(true);
        }
        const fallbackText = missingKey
          ? "I need a Gemini API key before I can chat. Open Settings to add it."
          : "I got distracted and missed that. Can you try again?";
        setConversation((previous) => [
          ...previous,
          {
            id: `pet-${Date.now()}`,
            role: "pet",
            text: fallbackText,
          },
        ]);
        return fallbackText;
      } finally {
        setIsSendingMessage(false);
      }
    },
    [isSendingMessage, registerInteraction]
  );

  const handleSendDialogMessage = useCallback(
    async (message: string) => {
      await sendMessage(message);
    },
    [sendMessage]
  );

  const handleSendQuickMessage = useCallback(
    async (message: string) => {
      const reply = await sendMessage(message);
      if (!reply) {
        return;
      }
      setQuickResponseMessage(reply);
      setIsResponseBubbleOpen(true);
    },
    [sendMessage]
  );

  const handleShowQuickTalk = useCallback(() => {
    registerInteraction();
    if (isDialogOpen || isShopOpen || isConfigOpen) {
      return;
    }
    setIsQuickTalkOpen((previous) => {
      if (!previous && !isDialogOpen && !isShopOpen && !isConfigOpen) {
        void AdjustWindowFromBottom(
          QUICK_TALK_WINDOW_SIZE.width,
          QUICK_TALK_WINDOW_SIZE.height
        );
      }
      return true;
    });
  }, [isConfigOpen, isDialogOpen, isShopOpen, registerInteraction]);

  const handleDismissResponseBubble = useCallback(() => {
    registerInteraction();
    setIsResponseBubbleOpen(false);
    setQuickResponseMessage(null);
  }, [registerInteraction]);

  const handleHideQuickTalk = useCallback(() => {
    registerInteraction();
    setIsQuickTalkOpen((previous) => {
      if (previous && !isDialogOpen && !isShopOpen && !isConfigOpen) {
        void AdjustWindowFromBottom(
          PET_WINDOW_DEFAULT_SIZE.width,
          PET_WINDOW_DEFAULT_SIZE.height
        );
      }
      return false;
    });
    handleDismissResponseBubble();
  }, [
    handleDismissResponseBubble,
    isDialogOpen,
    isConfigOpen,
    isShopOpen,
    registerInteraction,
  ]);

  const handleSpriteMoveStart = useCallback(() => {
    registerInteraction();
    setIsSpriteMoving(true);
  }, [registerInteraction]);

  const handleSpriteMoveEnd = useCallback(() => {
    registerInteraction();
    setIsSpriteMoving(false);
  }, [registerInteraction]);

  return (
    <div style={appShellStyle}>
      <Pet
        sprite={sprite}
        onOpenDialog={handleOpenDialog}
        onOpenShop={handleOpenShop}
        onOpenConfig={handleOpenConfig}
        itemEffect={itemEffect}
        isQuickTalkOpen={isQuickTalkOpen}
        onSendQuickMessage={handleSendQuickMessage}
        onRequestQuickTalk={handleShowQuickTalk}
        onCloseQuickTalk={handleHideQuickTalk}
        isChatBusy={isSendingMessage}
        isGeminiKeyMissing={isGeminiKeyMissing}
        quickResponseMessage={quickResponseMessage}
        isResponseBubbleOpen={isResponseBubbleOpen}
        onDismissResponseBubble={handleDismissResponseBubble}
        onSpriteMoveStart={handleSpriteMoveStart}
        onSpriteMoveEnd={handleSpriteMoveEnd}
        onUserInteraction={registerInteraction}
      />

      <PetDialog
        open={isDialogOpen}
        onClose={handleCloseDialog}
        onSend={handleSendDialogMessage}
        isBusy={isSendingMessage}
        isGeminiKeyMissing={isGeminiKeyMissing}
        messages={conversation}
        anchorPosition={dialogAnchor}
        onOpenSettings={openConfigFromDialog}
      />

      <PetShopDialog
        open={isShopOpen}
        onClose={handleCloseShop}
        anchorPosition={shopAnchor}
        onPurchaseItem={handlePurchaseItem}
      />

      <SystemConfigDialog
        open={isConfigOpen}
        onClose={handleCloseConfig}
        anchorPosition={configAnchor}
      />
    </div>
  );
}

export default App;
