import type { SpriteAnimationKey } from "./spriteFrames";

export type WanderDirection =
  | "left"
  | "right"
  | "up"
  | "down"
  | "leftUp"
  | "leftDown"
  | "rightUp"
  | "rightDown";

const wanderDirections: WanderDirection[] = [
  "left",
  "right",
  "up",
  "down",
  "leftUp",
  "leftDown",
  "rightUp",
  "rightDown",
];

export const wanderDirectionVectors: Record<WanderDirection, { x: number; y: number }> = {
  left: { x: -1, y: 0 },
  right: { x: 1, y: 0 },
  up: { x: 0, y: -1 },
  down: { x: 0, y: 1 },
  leftUp: { x: -1, y: -1 },
  leftDown: { x: -1, y: 1 },
  rightUp: { x: 1, y: -1 },
  rightDown: { x: 1, y: 1 },
};

export const wanderAnimationByDirection: Record<WanderDirection, SpriteAnimationKey> = {
  left: "move_left",
  right: "move_right",
  up: "move_far",
  // We don't have a dedicated down animation; reuse left-facing frames per requirement.
  down: "move_left",
  leftUp: "move_left",
  leftDown: "move_left",
  rightUp: "move_right",
  rightDown: "move_right",
};

export function randomWanderDirection(previous?: WanderDirection): WanderDirection {
  if (!previous) {
    return wanderDirections[Math.floor(Math.random() * wanderDirections.length)];
  }

  const candidates = wanderDirections.filter((direction) => direction !== previous);
  return candidates[Math.floor(Math.random() * candidates.length)];
}
