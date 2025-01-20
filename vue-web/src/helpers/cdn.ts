const ASSET_URL = import.meta.env.VITE_ASSETS_URL;

export function generateSrcset(uri: string, quality: number = 85, format: string = 'avif'): string {
  let srcset: string[] = [];

  for (const size of [240, 320, 640, 960, 1024]) {
    srcset = [...srcset, `${generateResizedSrc(uri, size, quality, format)} ${size}w`]
  }

  return srcset.join(',');
}

export function generateResizedSrc(uri: string, size: number = 640, quality: number = 85, format: string = 'avif'): string {
  return `${ASSET_URL}/cdn-cgi/image/width=${size},quality=${quality},format=${format}${uri}`;
}