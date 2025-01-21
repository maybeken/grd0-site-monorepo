import type { GalleryCategory } from "@/interfaces/Gallery";

export function formatCategoryName(gallery_category: GalleryCategory, value?: string): string | void {
  if (!value) return;

  const category_id = value.replace('/gallery/', '');
  let category_name = category_id;

  if (gallery_category && gallery_category[category_id]) {
    category_name = gallery_category[category_id].title || category_name;
  }

  return category_name.replace(/(^\w|\s\w)/g, (m: string) => m.toUpperCase())
}

export function getCategoryCover(gallery_category: GalleryCategory, value: string): string | undefined {
  const category_id = value.replace('/gallery/', '');

  if (gallery_category && gallery_category[category_id]) {
    return gallery_category[category_id].cover || undefined;
  }

  return undefined;
}