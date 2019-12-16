export const publicUrl = process.env.REACT_APP_PUBLIC_URL || window.location.origin;
export const publicWsUrl = publicUrl.replace('http', 'ws');
