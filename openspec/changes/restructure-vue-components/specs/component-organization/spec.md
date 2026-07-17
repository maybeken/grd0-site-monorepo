## ADDED Requirements

### Requirement: Feature-based component organization
The vue-web component directory SHALL be organized by feature rather than by layer. Each feature folder SHALL contain all components specific to that feature.

#### Scenario: Blog feature components
- **WHEN** a developer needs to find blog-related components
- **THEN** all blog components (BlogCard, BlogDetail, BlogSummary, ProfileIcon) SHALL be located in `src/components/blog/`

#### Scenario: Gallery feature components
- **WHEN** a developer needs to find gallery-related components
- **THEN** all gallery components (GalleryCard, GalleryGrid, GalleryZoom, DropdownSelection) SHALL be located in `src/components/gallery/`

#### Scenario: Music feature components
- **WHEN** a developer needs to find music-related components
- **THEN** all music components (MusicProgressBar, YouTubeEmbed) SHALL be located in `src/components/music/`

#### Scenario: Map feature components
- **WHEN** a developer needs to find map-related components
- **THEN** the map component (TravelersMap) SHALL be located in `src/components/map/`

### Requirement: Shared component folder
Components used across 2 or more features SHALL be placed in `src/components/shared/`.

#### Scenario: Cross-feature components
- **WHEN** a component is used by multiple features
- **THEN** it SHALL be located in `src/components/shared/` (CDNImage, CursorBlink, Icon, MarkdownDisplay, NotFound, Skeleton)

### Requirement: Layout component folder
App shell components (Navbar, SiteFooter, and their dependencies) SHALL be placed in `src/components/layout/`.

#### Scenario: Layout components
- **WHEN** a component is part of the app shell
- **THEN** it SHALL be located in `src/components/layout/` (Navbar, SiteFooter, DropdownBtn, MenuLinkBtn)

### Requirement: Component naming consistency
Component names SHALL be descriptive and follow consistent naming conventions.

#### Scenario: TravelersMap naming
- **WHEN** the map component is referenced
- **THEN** it SHALL be named `TravelersMap.vue` (not `Map.vue`)

#### Scenario: YouTubeEmbed naming
- **WHEN** the YouTube component is referenced
- **THEN** it SHALL be named `YouTubeEmbed.vue` (not `YouTube.vue`)

#### Scenario: MusicProgressBar naming
- **WHEN** the music progress component is referenced
- **THEN** it SHALL be named `MusicProgressBar.vue` (not `MusicProgressbar.vue`)

### Requirement: Existing feature folders preserved
Existing feature subfolders (coffee/, editor/) SHALL remain unchanged.

#### Scenario: Coffee folder
- **WHEN** the restructure is complete
- **THEN** `src/components/coffee/` SHALL contain all 7 coffee components unchanged

#### Scenario: Editor folder
- **WHEN** the restructure is complete
- **THEN** `src/components/editor/` SHALL contain PostTable.vue and SideBySideEditor.vue unchanged
