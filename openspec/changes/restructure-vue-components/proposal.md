## Why

The vue-web component structure is flat and layer-based, making it difficult to understand which components belong to which feature. With 21 components in the root `components/` directory, 13 are feature-specific but scattered across the codebase alongside their services, types, and views. This forces developers to jump between 5+ directories to understand a single feature.

## What Changes

- Reorganize `src/components/` from flat structure to feature-based folders
- Move feature-specific components into their respective feature folders (blog, gallery, music, map)
- Create `shared/` folder for components used across multiple features
- Create `layout/` folder for app shell components (Navbar, SiteFooter, and their dependencies)
- Rename 3 components for clarity:
  - `Map.vue` → `TravelersMap.vue` (less generic, matches feature name)
  - `YouTube.vue` → `YouTubeEmbed.vue` (describes what it is)
  - `MusicProgressbar.vue` → `MusicProgressBar.vue` (fix casing)
- Keep existing `coffee/` and `editor/` subfolders unchanged

## Capabilities

### New Capabilities
- `component-organization`: Feature-based folder structure for vue-web components with clear separation between shared, layout, and feature-specific components

### Modified Capabilities

## Impact

- **Code**: All component imports will need path updates (though `unplugin-vue-components` may handle this automatically)
- **Files affected**: 21 component files moved, 3 renamed, all dependent views/components updated
- **Breaking changes**: None (internal reorganization only, no API changes)
- **Dependencies**: No changes to package.json or external dependencies
- **Build**: No changes to vite.config.ts or build process (auto-import handles resolution)
