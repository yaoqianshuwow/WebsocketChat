function isBrowser() {
  return typeof window !== 'undefined';
}

function getToken(): string {
  if (!isBrowser()) return '';
  return localStorage.getItem('token') || '';
}

function normalizePlaceholder(url: string) {
  if (!url) return url;
  if (/__server_ip__/i.test(url)) {
    const hostname = isBrowser() ? window.location.hostname || '127.0.0.1' : '127.0.0.1';
    return url.replace(/__server_ip__/ig, hostname);
  }
  return url;
}

function proxyThroughGateway(url: string): string {
  const normalized = normalizePlaceholder(url);
  if (!normalized) return normalized;
  // Only proxy URLs that go to the file service
  if (normalized.includes('/files/')) {
    const baseUrl = import.meta.env.VITE_API_BASE_URL as string || '';
    const apiOrigin = baseUrl ? baseUrl.replace(/\/api\/v1$/, '') : window.location.origin;
    const token = getToken();
    if (token) {
      return `${apiOrigin}/api/v1/message/viewFile?fileUrl=${encodeURIComponent(normalized)}&token=${encodeURIComponent(token)}`;
    }
    return `${apiOrigin}/api/v1/message/viewFile?fileUrl=${encodeURIComponent(normalized)}`;
  }
  return normalized;
}

export function resolveAvatarUrl(avatar?: string) {
  if (!avatar || avatar.includes('�')) {
    return '';
  }
  return proxyThroughGateway(avatar);
}
