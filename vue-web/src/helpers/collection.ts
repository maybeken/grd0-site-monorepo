import type { GalleryCollection } from "@/interfaces/Gallery";

export function formatCollectionName(gallery_collection: GalleryCollection, value?: string): string | void {
  if (!value) return;

  const collection_id = value.replace('/gallery/', '');
  let collection_name = collection_id;

  if (gallery_collection && gallery_collection[collection_id]) {
    collection_name = gallery_collection[collection_id].title || collection_name;
  }

  return collection_name.replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase())
}

export function getCollectionCover(gallery_collection: GalleryCollection, value: string): string | undefined {
  const collection_id = value.replace('/gallery/', '');

  if (gallery_collection && gallery_collection[collection_id]) {
    return gallery_collection[collection_id].cover || undefined;
  }

  return undefined;
}