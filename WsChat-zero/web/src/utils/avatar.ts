const AVATAR_POOL = ['/fallback-avatar.jpg', '/fallback-avatar-2.jpg'];

function stableIndex(seed: number | string | undefined) {
  if (seed === undefined || seed === null) return 0;
  const text = String(seed);
  let hash = 0;
  for (let i = 0; i < text.length; i += 1) {
    hash = (hash * 31 + text.charCodeAt(i)) >>> 0;
  }
  return hash % AVATAR_POOL.length;
}

export function pickAvatar(seed: number | string | undefined, avatar?: string) {
  if (avatar && !avatar.includes('�')) {
    return avatar;
  }
  return AVATAR_POOL[stableIndex(seed)];
}

