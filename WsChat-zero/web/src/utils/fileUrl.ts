function getToken(): string {
  return localStorage.getItem('token') || '';
}

export function resolveFileUrl(url?: string) {
  if (!url) return '';
  const normalized = /__server_ip__/i.test(url)
    ? url.replace(/__server_ip__/ig, window.location.hostname || '127.0.0.1')
    : url;
  if (normalized.includes('/files/')) {
    const token = getToken();
    const base = (import.meta.env.VITE_API_BASE_URL as string) || '';
    const apiOrigin = base ? base.replace(/\/api\/v1$/, '') : window.location.origin;
    const resolved = `${apiOrigin}/api/v1/message/viewFile?fileUrl=${encodeURIComponent(normalized)}`;
    return token ? `${resolved}&token=${encodeURIComponent(token)}` : resolved;
  }
  return normalized;
}
