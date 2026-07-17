## Context

The vue-web project currently has 21 components in a flat `src/components/` directory. While `coffee/` and `editor/` subfolders exist, 13 feature-specific components are mixed with shared components, making it difficult to understand feature boundaries. The project uses `unplugin-vue-components` for auto-import, which means component resolution is automatic regardless of folder structure.

## Goals / Non-Goals

**Goals:**
- Organize components by feature for better discoverability
- Separate shared components from feature-specific ones
- Rename 3 components for clarity (Map → TravelersMap, YouTube → YouTubeEmbed, MusicProgressbar → MusicProgressBar)
- Maintain auto-import functionality

**Non-Goals:**
- Restructure services, helpers, interfaces, or views (out of scope)
- Change component functionality or behavior
- Update component APIs or props
- Modify build configuration

## Decisions

### Decision 1: Feature-based folders over layer-based
**Choice:** Organize by feature (blog/, gallery/, music/, map/) rather than by layer (all components flat).

**Rationale:** Feature-based organization makes it easier to understand what belongs to a feature. When working on the blog, all blog components are in one place. This scales better as the project grows.

**Alternatives considered:**
- Keep flat structure: Rejected because it's already causing confusion
- Full feature folders (components + services + views together): Too disruptive, out of scope

### Decision 2: Shared vs Layout separation
**Choice:** Create two special folders: `shared/` for cross-feature components and `layout/` for app shell components.

**Rationale:** Navbar and SiteFooter are not "shared" in the same way as Icon or Skeleton — they're part of the app shell. Separating them makes the distinction clear.

### Decision 3: Rename scope
**Choice:** Rename only 3 components (Map, YouTube, MusicProgressbar).

**Rationale:** These names are either too generic (Map), vague (YouTube), or have inconsistent casing (MusicProgressbar). Other names like CDNImage are acceptable despite being implementation details.

**Alternatives considered:**
- Rename all abbreviations (Btn → Button): Rejected per user preference to keep "Btn"
- Rename CDNImage: Rejected, current name is acceptable

### Decision 4: Auto-import compatibility
**Choice:** Rely on `unplugin-vue-components` to handle path resolution automatically.

**Rationale:** The project already uses auto-import, so moving files won't break imports. No need to manually update import statements in most cases.

**Risks:**
- Some explicit imports may need updating → Mitigation: Search for explicit imports and update them
- components.d.ts may need regeneration → Mitigation: Run dev server to trigger regeneration

## Risks / Trade-offs

**Risk:** Breaking explicit imports
- **Mitigation:** Search for all explicit imports of renamed/moved components and update them

**Risk:** Merge conflicts with other branches
- **Mitigation:** Coordinate with team, restructure in a single PR

**Risk:** Forgetting to update components.d.ts
- **Mitigation:** Run `pnpm dev` after restructure to trigger auto-regeneration

**Trade-off:** More nesting for consistency
- Some folders (map/) will have only 1 file, but this maintains the pattern
