export interface Asset {
  filename: string,
  exif?: {
    datetime?: string,
    shutter?: string,
    fstop?: number,
    iso?: number,
    focal?: number,
    equipment?: {
      camera?: string,
      lens?: string,
    }
  }
};
  
export interface AssetFileList {
  [key: string]: Asset[],
}