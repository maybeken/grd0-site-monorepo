## 1. Setup & Infrastructure

- [x] 1.1 Install vue-i18n and yaml dependencies in vue-web
- [x] 1.2 Create src/i18n.ts with vue-i18n plugin initialization
- [x] 1.3 Create src/locales/ directory structure
- [x] 1.4 Update main.ts to import and use i18n plugin
- [x] 1.5 Create vite-plugin-gettext-yaml-i18n.ts to transform gettext-style YAML to vue-i18n format at build time

## 2. Extract Strings to YAML

- [x] 2.1 Create src/locales/en.yaml with gettext conventions (msgid, msgstr, msgctxt, #., #:)
- [x] 2.2 Extract navigation strings from Navbar.vue (menu items, tagline, external links)
- [x] 2.3 Extract footer strings from SiteFooter.vue (copyright, license text)
- [x] 2.4 Extract gallery strings from GalleryView.vue and GalleryCard.vue
- [x] 2.5 Extract map strings from MapView.vue and TravelersMap.vue
- [x] 2.6 Extract music strings from MusicView.vue
- [x] 2.7 Extract blog strings from BlogView.vue, BlogDetail.vue, BlogSummary.vue
- [x] 2.8 Extract blog editor strings from BlogEditorView.vue and SideBySideEditor.vue
- [x] 2.9 Extract coffee view strings from CoffeeView.vue
- [x] 2.10 Extract coffee editor strings from CoffeeEditorView.vue
- [x] 2.11 Extract tasting form strings from TastingForm.vue
- [x] 2.12 Extract tasting card strings from TastingCard.vue
- [x] 2.13 Extract bean form strings from BeanForm.vue
- [x] 2.14 Extract equipment form strings from EquipmentForm.vue
- [x] 2.15 Extract coffee constants to YAML (grind sizes, equipment types, taste dimensions)
- [x] 2.16 Extract 404/not found strings from NotFound.vue and NotFoundView.vue
- [x] 2.17 Extract dot matrix tool strings from DotMatrix.vue
- [x] 2.18 Extract tldraw tool strings from TldrawView.vue
- [x] 2.19 Extract redirect strings from RedirectView.vue
- [x] 2.20 Extract common strings (save, cancel, edit, delete, search, etc.)
- [x] 2.21 Extract dropdown and pagination strings from DropdownBtn.vue, PostTable.vue

## 3. Refactor Coffee Constants

- [x] 3.1 Refactor GRIND_SIZES in helpers/coffee.ts to use { value, labelKey } structure
- [x] 3.2 Refactor EQUIPMENT_TYPES in helpers/coffee.ts to use { value, labelKey } structure
- [x] 3.3 Refactor TASTE_DIMENSIONS in helpers/coffee.ts to use { value, labelKey } structure
- [x] 3.4 Update TastingForm.vue to use $t(grindSize.labelKey) for display
- [x] 3.5 Update TastingCard.vue to use $t() for taste dimension labels
- [x] 3.6 Update TasteRadarChart.vue to use $t() for chart labels
- [x] 3.7 Update BeanForm.vue and EquipmentForm.vue dropdowns to use $t() for labels

## 4. Update Router Configuration

- [x] 4.1 Update router meta.title to use translation keys in router/index.ts
- [x] 4.2 Update document title composition logic to translate meta.title
- [x] 4.3 Verify all route titles are translated correctly

## 5. Update Date Formatting

- [x] 5.1 Update TastingCard.vue to use locale-aware date formatting instead of hardcoded en-GB
- [x] 5.2 Verify date display works correctly with vue-i18n

## 6. Replace Strings in Components

- [x] 6.1 Replace hardcoded strings in Navbar.vue with $t() calls
- [x] 6.2 Replace hardcoded strings in SiteFooter.vue with $t() calls
- [x] 6.3 Replace hardcoded strings in GalleryView.vue with $t() calls
- [x] 6.4 Replace hardcoded strings in MapView.vue with $t() calls
- [x] 6.5 Replace hardcoded strings in MusicView.vue with $t() calls
- [x] 6.6 Replace hardcoded strings in BlogView.vue, BlogDetail.vue, BlogSummary.vue with $t() calls
- [x] 6.7 Replace hardcoded strings in BlogEditorView.vue and SideBySideEditor.vue with $t() calls
- [x] 6.8 Replace hardcoded strings in CoffeeView.vue with $t() calls
- [x] 6.9 Replace hardcoded strings in CoffeeEditorView.vue with $t() calls
- [x] 6.10 Replace hardcoded strings in TastingForm.vue with $t() calls
- [x] 6.11 Replace hardcoded strings in TastingCard.vue with $t() calls
- [x] 6.12 Replace hardcoded strings in BeanForm.vue with $t() calls
- [x] 6.13 Replace hardcoded strings in EquipmentForm.vue with $t() calls
- [x] 6.14 Replace hardcoded strings in NotFound.vue and NotFoundView.vue with $t() calls
- [x] 6.15 Replace hardcoded strings in DotMatrix.vue with $t() calls
- [x] 6.16 Replace hardcoded strings in TldrawView.vue with $t() calls
- [x] 6.17 Replace hardcoded strings in RedirectView.vue with $t() calls
- [x] 6.18 Replace hardcoded strings in DropdownBtn.vue and PostTable.vue with $t() calls

## 7. Verification

- [x] 7.1 Search for remaining hardcoded strings in src/**/*.vue and src/**/*.ts
- [x] 7.2 Verify no hardcoded UI strings remain (excluding comments, brand names, technical constants)
- [x] 7.3 Run type-check to ensure no TypeScript errors
- [x] 7.4 Run lint to ensure code style compliance
- [x] 7.5 Manual testing: verify all pages display English text correctly
- [x] 7.6 Manual testing: verify form placeholders, alerts, and confirmations work
- [x] 7.7 Manual testing: verify router page titles update correctly
