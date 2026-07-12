## Context

The `TastingNote` entity captures brew recipe parameters (grind size, dose, water volume, brew time) alongside taste profile ratings. Water temperature is a key extraction variable but is currently missing. The existing architecture supports optional, nullable brew parameters with "I don't know" (IDK) toggle behavior on the frontend and null-display on the card view.

The Go API uses GORM with SQLite and auto-migration. The Vue frontend renders tasting cards and provides a form with IDK toggles for optional fields.

## Goals / Non-Goals

**Goals:**
- Add `water_temperature` as an optional, nullable integer field (°C) to the tasting note entity
- Extend the existing form to accept water temperature input with IDK toggle
- Display water temperature on the tasting card, showing "IDK" when null
- Maintain consistency with existing brew parameter patterns

**Non-Goals:**
- Temperature unit conversion (°F/°C) — always store and display in °C
- Validation ranges (e.g., rejecting unrealistic temperatures) — follow existing pattern of accepting any integer
- Backend validation beyond nullability — frontend handles IDK toggle logic

## Decisions

**Field type and placement:**
- Use `*int` (nullable integer) in Go struct, matching `coffee_dose`, `water_in`, `brew_time`
- Place after `brew_time` in the struct and form — groups with other brew recipe fields
- JSON key: `water_temperature` (snake_case, consistent with existing fields)

**UI treatment:**
- Form: numeric input with IDK toggle, same pattern as other optional brew fields
- Card: display value followed by "°C" unit, or "IDK" when null
- No special formatting or color coding — keep it simple

**Migration:**
- Rely on GORM `AutoMigrate` to add the nullable column — no manual migration script needed
- Existing records will have `NULL` for this field, which displays as "IDK" — acceptable

**API changes:**
- No handler modifications required — existing `UpsertCoffeeTasting` and `GetCoffeeTastings` are field-agnostic
- TypeScript interface gains `water_temperature?: number | null`

## Risks / Trade-offs

**[Risk] No temperature range validation** → Mitigation: Follow existing pattern (no validation on `coffee_dose`, `water_in`, etc.). If invalid data becomes a problem, add validation later.

**[Risk] Unit confusion (°C vs °F)** → Mitigation: Always use °C. If international users need °F, add a frontend display preference later (non-goal for this change).

**[Trade-off] Nullable integer vs. sentinel value** → Chose nullable to match existing pattern and support IDK toggle cleanly. Sentinel values (e.g., -1) would complicate queries and display logic.
