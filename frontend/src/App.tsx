import { CSSProperties, useEffect, useMemo, useState } from "react";
import { PetSprite } from "../wailsjs/go/app/App";
import { WindowSetAlwaysOnTop } from "../wailsjs/runtime/runtime";
import { Pet } from "./components/Pet";

const noDragStyle: CSSProperties = {
  ["--wails-draggable" as any]: "no-drag",
};

const TRANSPARENT_BACKGROUND = "rgba(0, 0, 0, 0)";
const DEV_TINT_BACKGROUND = "rgba(20, 20, 20, 0.33)";

const baseAppShellStyle: CSSProperties = {
  display: "flex",
  flexDirection: "column",
  alignItems: "flex-start",
  justifyContent: "flex-start",
  padding: "16px",
  boxSizing: "border-box",
  userSelect: "none",
  pointerEvents: "none",
  backdropFilter: "blur(6px)",
  position: "fixed",
  inset: 0,
  transform: "translateZ(0)",
  willChange: "transform, opacity",
  background: TRANSPARENT_BACKGROUND,
  backfaceVisibility: "hidden",
};

const baseHintStyle: CSSProperties = {
  display: "none",
  opacity: 0,
  animation: "none",
};

function App() {
  const [sprite, setSprite] = useState<string>("");
  const [isFloating, setIsFloating] = useState(true);
  const [showHint, setShowHint] = useState(true);

  const isDevMode = import.meta.env.DEV;
  const windowBackground = isDevMode
    ? DEV_TINT_BACKGROUND
    : TRANSPARENT_BACKGROUND;

  const appShellStyle = useMemo<CSSProperties>(
    () => ({ ...baseAppShellStyle, background: windowBackground }),
    [windowBackground]
  );

  const hintStyle = useMemo<CSSProperties>(
    () => ({
      ...noDragStyle,
      ...baseHintStyle,
    }),
    []
  );

  useEffect(() => {
    let isMounted = true;

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

    const hintTimeout = window.setTimeout(() => setShowHint(false), 6000);

    return () => {
      isMounted = false;
      window.clearTimeout(hintTimeout);

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
  }, [windowBackground]);

  const toggleFloating = () => setIsFloating((prev) => !prev);

  return (
    <div style={appShellStyle}>
      <Pet
        sprite={sprite}
        isFloating={isFloating}
        onToggleFloating={toggleFloating}
      />

      {showHint && (
        <div style={hintStyle}>Drag me anywhere · Double-click to pause</div>
      )}
    </div>
  );
}

export default App;
