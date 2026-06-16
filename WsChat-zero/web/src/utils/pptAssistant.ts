export const PPT_ASSISTANT_NAME = 'PPT小助手';
export const PPT_ASSISTANT_ROUTE = '/ppt-assistant';
export const PPT_ASSISTANT_DEFAULT_URL = 'http://127.0.0.1:8123/';

export function getPptAssistantUrl() {
  return import.meta.env.VITE_PPT_ASSISTANT_URL || PPT_ASSISTANT_DEFAULT_URL;
}
