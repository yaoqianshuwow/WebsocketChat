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

function getToken(): string {
  return localStorage.getItem('token') || '';
}

function wrapFileUrl(url: string): string {
  if (url && url.includes('/files/')) {
    const token = getToken();
    if (token) {
      // Use same origin as API client (gateway port)
      const base = import.meta.env.VITE_API_BASE_URL as string || '';
      const apiOrigin = base ? base.replace(/\/api\/v1$/, '') : window.location.origin;
      return `${apiOrigin}/api/v1/message/viewFile?fileUrl=${encodeURIComponent(url)}&token=${encodeURIComponent(token)}`;
    }
  }
  return url;
}

export function pickAvatar(seed: number | string | undefined, avatar?: string) {
  if (avatar && !avatar.includes('�')) {
    return wrapFileUrl(avatar);
  }
  return AVATAR_POOL[stableIndex(seed)];
}

