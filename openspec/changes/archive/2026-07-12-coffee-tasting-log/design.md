## Context

The site (grd0.net) is a Vue 3 + Vite SPA backed by a Go/Echo API with SQLite via GORM. Existing features (blog, gallery, music, traveler's map) all follow the same full-stack pattern: schema → handler → route → service → view. The coffee tasting log is a new feature that fits this pattern but introduces a radar chart visualization (Chart.js) and a more complex card layout.

## Goals / Non-Goals

**Goals:**
- Full-stack coffee tasting log: record beans, equipment, and tasting notes
- Public read access, admin-only write access (JWT-protected)
- 1:1 square card layout with equipment image, bean/brew info, radar chart, and notes
- "I don't know" toggle for all optional fields, displayed as literal "IDK"
- Auto-calculated water ratio from coffee dose and water in
- Brew time input as mm:ss, stored as integer seconds
- Pinned tastings sort to top, then newest first

**Non-Goals:**
- Image upload for tastings (excluded for now)
- Bean detail page (bean info shown inline in each tasting card)
- Equipment management UI for end users (admin-only)
- Filtering/searching the tasting log
- Pagination (assume small dataset for now)
- Multi-user support (single admin)

## Decisions

### 1. All entities use UUID primary keys via BaseColumns

**Decision**: CoffeeBean, BrewEquipment, and TastingNote all embed `BaseColumns` (UUID PK, timestamps, soft delete).

**Rationale**: Consistency with existing schemas (Music, Blog, etc.). Simpler code — no need for a separate incremental base type. The bean ID being "incremental" is a nice-to-have that doesn't outweigh the consistency benefit.

**Alternative considered**: Integer auto-increment for beans/equipment (cleaner dropdown display). Rejected because it breaks the BaseColumns pattern and adds a second base type.

### 2. Equipment images as local static files with mapping array

**Decision**: Store equipment images in `vue-web/src/assets/coffee/equipment/`. Maintain a constant mapping array in code: `EQUIPMENT_IMAGE_MAP: Record<string, string>` linking equipment type to image path.

**Rationale**: No need for CDN upload for a fixed set of images. Mapping array keeps the coupling explicit and type-safe. Adding new equipment types just requires adding an image and a mapping entry.

**Alternative considered**: CDN-hosted images. Rejected — overkill for a fixed set of ~8 equipment types.

### 3. Chart.js + vue-chartjs for radar chart

**Decision**: Install `chart.js` and `vue-chartjs` for the taste profile radar chart.

**Rationale**: Chart.js is the standard lightweight charting library with good Vue integration. The radar chart is a built-in chart type. Bundle impact is acceptable (~60kb gzipped).

**Alternative considered**: Custom SVG rendering. Rejected — more code to maintain, harder to tune.

### 4. Taste profile as fixed columns, not dynamic

**Decision**: The 7 taste dimensions (fruity, sour, sweetness, nutty, spice, floral, green) are fixed columns on the `tasting_notes` table, each an optional integer 0-10.

**Rationale**: Fixed dimensions match the radar chart's fixed axes. Simpler than a dynamic key-value approach. If new dimensions are needed later, a migration can add them.

**Alternative considered**: Separate `taste_notes` table with (tasting_id, dimension, value). Rejected — over-engineered for a fixed set, adds join complexity.

### 5. Water ratio auto-calculation on frontend only

**Decision**: `ratio` is stored as a REAL column but computed entirely on the frontend from `water_in / coffee_dose` when both are present. The API does not calculate or enforce ratio — it stores whatever value the frontend sends.

**Rationale**: Keeps the API dumb (just stores what's sent). The frontend has the context to show "IDK" when either input is missing and to update ratio reactively as the user changes inputs. The stored ratio is a convenience for display.

### 6. Brew time stored as integer seconds

**Decision**: Frontend accepts mm:ss input, converts to total seconds for storage. On read, converts back to mm:ss for display.

**Rationale**: Integer seconds is simple to store, sort, and compare. The mm:ss format is purely a UI concern.

### 7. Card layout: 1:1 square, image top, two-column middle, notes bottom

**Decision**: Each tasting card uses `aspect-square` with three rows:
1. Equipment image (full width)
2. Left: date/rating + bean info + brew recipe. Right: radar chart
3. Notes (full width)

**Rationale**: Square cards create a consistent grid. The two-column middle balances text-heavy info with visual radar chart. Notes at bottom allow variable-length text without disrupting the layout.

## Risks / Trade-offs

### [Risk] Card content overflow in 1:1 square
Long bean names, many brew fields, or verbose notes could overflow the fixed square. → **Mitigation**: Use `overflow-hidden` with text truncation. Notes section can have a max-height with scroll. Font sizes responsive via Tailwind.

### [Risk] Radar chart rendering with all "IDK" values
If all 7 taste dimensions are "IDK", the radar chart has no data to render. → **Mitigation**: Show a placeholder message ("No taste data") instead of an empty chart.

### [Risk] Equipment image mapping gaps
If a new equipment type is added without an image, the card shows a broken image. → **Mitigation**: Fallback to a default "unknown" image if the mapping doesn't have an entry.

### [Trade-off] No pagination
Small dataset assumption means the entire tasting log loads at once. → **Mitigation**: Acceptable for personal use. Add pagination later if needed — API already supports offset/limit patterns from blog.

### [Trade-off] Ratio stored redundantly
Ratio can be derived from water_in / coffee_dose, but we store it anyway. → **Mitigation**: Frontend computes on save. If inputs change, ratio updates. Stored value is for display convenience and historical accuracy.
