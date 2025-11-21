import {
  CSSProperties,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from "react";
import { ChatWithPet, PetSprite } from "../wailsjs/go/app/App";
import {
  AdjustWindowFromBottom,
  SetDevicePixelRatio,
} from "../wailsjs/go/window/WindowService";
import { WindowSetAlwaysOnTop } from "../wailsjs/runtime/runtime";
import { Pet } from "./components/Pet/Pet";
import {
  PetDialog,
  type ConversationMessage,
} from "./components/Pet/PetDialog";

const noDragStyle: CSSProperties = {
  ["--wails-draggable" as any]: "no-drag",
};

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
const DIALOG_WINDOW_SIZE = { width: 900, height: 500 };

function App() {
  const [sprite, setSprite] = useState<string>("");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isSendingMessage, setIsSendingMessage] = useState(false);
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
  const [isQuickTalkOpen, setIsQuickTalkOpen] = useState(false);
  const [quickResponseMessage, setQuickResponseMessage] = useState<
    string | null
  >(null);
  const [isResponseBubbleOpen, setIsResponseBubbleOpen] = useState(false);

  const isDevMode = import.meta.env.DEV;
  const windowBackground = isDevMode
    ? DEV_TINT_BACKGROUND
    : TRANSPARENT_BACKGROUND;

  const appShellStyle = useMemo<CSSProperties>(
    () => ({ ...baseAppShellStyle, background: windowBackground }),
    [windowBackground]
  );

  useEffect(() => {
    let isMounted = true;
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

    console.log(
      "body background:",
      getComputedStyle(document.body).backgroundColor
    );

    WindowSetAlwaysOnTop(true);

    PetSprite()
      .then((data) => {
        if (isMounted) {
          setSprite(data);
        }
      })
      .catch((err) => console.error("Pet sprite failed to load", err));

    return () => {
      isMounted = false;

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
  }, []);

  const handleOpenDialog = useCallback(
    (anchor: { left: number; top: number }) => {
      setDialogAnchor(anchor);
      void AdjustWindowFromBottom(
        DIALOG_WINDOW_SIZE.width,
        DIALOG_WINDOW_SIZE.height
      );
      setIsDialogOpen(true);
    },
    []
  );
  const handleCloseDialog = useCallback(() => {
    setIsDialogOpen(false);
    setDialogAnchor(null);
    void AdjustWindowFromBottom(
      PET_WINDOW_DEFAULT_SIZE.width,
      PET_WINDOW_DEFAULT_SIZE.height
    );
  }, []);

  const sendMessage = useCallback(
    async (message: string) => {
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
        const fallbackText =
          "I got distracted and missed that. Can you try again?";
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
    [isSendingMessage]
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
    setIsQuickTalkOpen(true);
  }, []);

  const handleDismissResponseBubble = () => {
    setIsResponseBubbleOpen(false);
    setQuickResponseMessage(null);
  };

  return (
    <div style={appShellStyle}>
      <Pet
        sprite={sprite}
        onOpenDialog={handleOpenDialog}
        isQuickTalkOpen={isQuickTalkOpen}
        onSendQuickMessage={handleSendQuickMessage}
        onRequestQuickTalk={handleShowQuickTalk}
        isChatBusy={isSendingMessage}
        quickResponseMessage={quickResponseMessage}
        isResponseBubbleOpen={isResponseBubbleOpen}
        onDismissResponseBubble={handleDismissResponseBubble}
      />

      <PetDialog
        open={isDialogOpen}
        onClose={handleCloseDialog}
        onSend={handleSendDialogMessage}
        isBusy={isSendingMessage}
        messages={conversation}
        anchorPosition={dialogAnchor}
      />
    </div>
  );
}

export default App;
