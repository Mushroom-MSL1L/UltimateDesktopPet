import {
  WindowGetPosition,
  WindowGetSize,
  WindowSetPosition,
} from "../../wailsjs/runtime/runtime";

export type WindowBounds = {
  left: number;
  top: number;
  width: number;
  height: number;
};

export type WindowMoveResult = {
  bounds: WindowBounds;
  didMove: boolean;
};

const clampNumber = (value: number, min: number, max: number) =>
  Math.min(Math.max(value, min), max);

function getScreenSize() {
  const screenWidth = window.screen?.availWidth ?? window.screen?.width ?? 0;
  const screenHeight = window.screen?.availHeight ?? window.screen?.height ?? 0;
  return { width: screenWidth, height: screenHeight };
}

export async function getWindowBounds(): Promise<WindowBounds | null> {
  try {
    const { w: width, h: height } = await WindowGetSize();
    const { x: left, y: top } = await WindowGetPosition();
    return { left, top, width, height };
  } catch (error) {
    console.warn("Failed to read window bounds", error);
    return null;
  }
}

export async function moveWindowWithinScreen(
  deltaX: number,
  deltaY: number
): Promise<WindowMoveResult | null> {
  const bounds = await getWindowBounds();
  if (!bounds) {
    return null;
  }

  const { width: screenWidth, height: screenHeight } = getScreenSize();
  const maxLeft = screenWidth > 0 ? Math.max(0, screenWidth - bounds.width) : undefined;
  const maxTop = screenHeight > 0 ? Math.max(0, screenHeight - bounds.height) : undefined;

  const targetLeft =
    typeof maxLeft === "number"
      ? clampNumber(bounds.left + deltaX, 0, maxLeft)
      : bounds.left + deltaX;
  const targetTop =
    typeof maxTop === "number"
      ? clampNumber(bounds.top + deltaY, 0, maxTop)
      : bounds.top + deltaY;

  const didMove = targetLeft !== bounds.left || targetTop !== bounds.top;

  try {
    await WindowSetPosition(targetLeft, targetTop);
    return {
      didMove,
      bounds: { ...bounds, left: targetLeft, top: targetTop },
    };
  } catch (error) {
    console.warn("Failed to move window", error);
    return null;
  }
}
