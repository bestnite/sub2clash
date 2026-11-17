export function appUrl(path: string, pageUrl: string = window.location.href) {
  const base = new URL(pageUrl);
  base.pathname = base.pathname.replace(/\/?$/, "/");
  base.search = "";
  base.hash = "";
  return new URL(path, base).href;
}

export function base64EncodeUnicode(str: string) {
  const bytes = new TextEncoder().encode(str);
  let binary = "";
  bytes.forEach((b) => (binary += String.fromCharCode(b)));
  return btoa(binary);
}

export function base64decodeUnicode(str: string) {
  const binaryString = atob(str);
  const bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  return new TextDecoder().decode(bytes);
}
