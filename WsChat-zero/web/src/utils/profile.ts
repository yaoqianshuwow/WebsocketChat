export function buildProfilePath(userId?: number | string | null) {
  if (userId === undefined || userId === null || userId === '') {
    return '/profile';
  }
  return `/profile?userId=${encodeURIComponent(String(userId))}`;
}
