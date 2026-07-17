## 1. Create folder structure

- [ ] 1.1 Create `src/components/shared/` directory
- [ ] 1.2 Create `src/components/layout/` directory
- [ ] 1.3 Create `src/components/blog/` directory
- [ ] 1.4 Create `src/components/gallery/` directory
- [ ] 1.5 Create `src/components/music/` directory
- [ ] 1.6 Create `src/components/map/` directory

## 2. Move shared components

- [ ] 2.1 Move `CDNImage.vue` to `shared/`
- [ ] 2.2 Move `CursorBlink.vue` to `shared/`
- [ ] 2.3 Move `Icon.vue` to `shared/`
- [ ] 2.4 Move `MarkdownDisplay.vue` to `shared/`
- [ ] 2.5 Move `NotFound.vue` to `shared/`
- [ ] 2.6 Move `Skeleton.vue` to `shared/`

## 3. Move layout components

- [ ] 3.1 Move `Navbar.vue` to `layout/`
- [ ] 3.2 Move `SiteFooter.vue` to `layout/`
- [ ] 3.3 Move `DropdownBtn.vue` to `layout/`
- [ ] 3.4 Move `MenuLinkBtn.vue` to `layout/`

## 4. Move blog components

- [ ] 4.1 Move `BlogCard.vue` to `blog/`
- [ ] 4.2 Move `BlogDetail.vue` to `blog/`
- [ ] 4.3 Move `BlogSummary.vue` to `blog/`
- [ ] 4.4 Move `ProfileIcon.vue` to `blog/`

## 5. Move gallery components

- [ ] 5.1 Move `GalleryCard.vue` to `gallery/`
- [ ] 5.2 Move `GalleryGrid.vue` to `gallery/`
- [ ] 5.3 Move `GalleryZoom.vue` to `gallery/`
- [ ] 5.4 Move `DropdownSelection.vue` to `gallery/`

## 6. Move music components

- [ ] 6.1 Move and rename `YouTube.vue` to `music/YouTubeEmbed.vue`
- [ ] 6.2 Move and rename `MusicProgressbar.vue` to `music/MusicProgressBar.vue`

## 7. Move map component

- [ ] 7.1 Move and rename `Map.vue` to `map/TravelersMap.vue`

## 8. Update imports

- [ ] 8.1 Search for explicit imports of renamed components (TravelersMap, YouTubeEmbed, MusicProgressBar) and update them
- [ ] 8.2 Search for explicit imports of moved components and update paths if needed
- [ ] 8.3 Run `pnpm dev` to trigger components.d.ts regeneration
- [ ] 8.4 Run `pnpm type-check` to verify no type errors
- [ ] 8.5 Run `pnpm lint` to verify no lint errors

## 9. Verify

- [ ] 9.1 Verify `pnpm build` succeeds
- [ ] 9.2 Verify dev server runs without errors
- [ ] 9.3 Spot-check key pages (blog, gallery, music, map) in browser
