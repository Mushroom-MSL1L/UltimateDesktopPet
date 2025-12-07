import {
  PetFramesMoveFar,
  PetFramesMoveLeft,
  PetFramesMoveRight,
  PetFramesStand,
} from "../../wailsjs/go/pet/PetMeta";

export type SpriteAnimationKey = "stand" | "move_left" | "move_right" | "move_far";

type SpriteFramesCache = Partial<Record<SpriteAnimationKey, string[]>>;

type FrameLoader = () => Promise<string[] | null | undefined>;

const animationLoaders: Record<SpriteAnimationKey, FrameLoader> = {
  stand: PetFramesStand,
  move_left: PetFramesMoveLeft,
  move_right: PetFramesMoveRight,
  move_far: PetFramesMoveFar,
};

async function loadFrames(animation: SpriteAnimationKey): Promise<string[]> {
  try {
    const frames = await animationLoaders[animation]();
    return frames ?? [];
  } catch (error) {
    console.warn(`Failed to load frames for ${animation}`, error);
    return [];
  }
}

export async function ensureAnimationFrames(
  animation: SpriteAnimationKey,
  cache: SpriteFramesCache
): Promise<string[]> {
  if (!cache[animation]?.length) {
    cache[animation] = await loadFrames(animation);
  }
  return cache[animation] ?? [];
}

export async function preloadAnimations(
  animations: SpriteAnimationKey[]
): Promise<SpriteFramesCache> {
  const cache: SpriteFramesCache = {};
  await Promise.all(
    animations.map(async (animation) => {
      cache[animation] = await loadFrames(animation);
    })
  );
  return cache;
}

export function getCachedFrames(
  animation: SpriteAnimationKey,
  cache: SpriteFramesCache
): string[] {
  return cache[animation] ?? [];
}

export type { SpriteFramesCache };
