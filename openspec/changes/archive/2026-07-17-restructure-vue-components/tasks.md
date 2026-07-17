## 1. Create folder structure

- [x] 1.1 Create `src/components/shared/` directory
- [x] 1.2 Create `src/components/layout/` directory
- [x] 1.3 Create `src/components/blog/` directory
- [x] 1.4 Create `src/components/gallery/` directory
- [x] 1.5 Create `src/components/music/` directory
- [x] 1.6 Create `src/components/map/` directory

## 2. Move shared components

- [x] 2.1 Move `CDNImage.vue` to `shared/`
- [x] 2.2 Move `CursorBlink.vue` to `shared/`
- [x] 2.3 Move `Icon.vue` to `shared/`
- [x] 2.4 Move `MarkdownDisplay.vue` to `shared/`
- [x] 2.5 Move `NotFound.vue` to `shared/`
- [x] 2.6 Move `Skeleton.vue` to `shared/`

## 3. Move layout components

- [x] 3.1 Move `Navbar.vue` to `layout/`
- [x] 3.2 Move `SiteFooter.vue` to `layout/`
- [x] 3.3 Move `DropdownBtn.vue` to `layout/`
- [x] 3.4 Move `MenuLinkBtn.vue` to `layout/`

## 4. Move blog components

- [x] 4.1 Move `BlogCard.vue` to `blog/`
- [x] 4.2 Move `BlogDetail.vue` to `blog/`
- [x] 4.3 Move `BlogSummary.vue` to `blog/`
- [x] 4.4 Move `ProfileIcon.vue` to `blog/`

## 5. Move gallery components

- [x] 5.1 Move `GalleryCard.vue` to `gallery/`
- [x] 5.2 Move `GalleryGrid.vue` to `gallery/`
- [x] 5.3 Move `GalleryZoom.vue` to `gallery/`
- [x] 5.4 Move `DropdownSelection.vue` to `gallery/`

## 6. Move music components

- [x] 6.1 Move and rename `YouTube.vue` to `music/YouTubeEmbed.vue`
- [x] 6.2 Move and rename `MusicProgressbar.vue` to `music/MusicProgressBar.vue`

## 7. Move map component

- [x] 7.1 Move and rename `Map.vue` to `map/TravelersMap.vue`

## 8. Update imports

- [x] 8.1 Search for explicit imports of renamed components (TravelersMap, YouTubeEmbed, MusicProgressBar) and update them
- [x] 8.2 Search for explicit imports of moved components and update paths if needed
-- [x] 8.3 Run `pnpm dev` to trigger components.d.ts regeneration
-- [x] 8.4 Run `pnpm type-check` to verify no type errors
- [x] 8.5 Run `pnpm lint` to verify no lint errors

## 9. Verify

- [x] 9.1 Verify `pnpm build` succeeds
- [x] 9.2 Verify dev server runs without errors
- [x] 9.3 Spot-check key pages (blog, gallery, music, map) in browser
