## Why

Water temperature is a critical brew parameter for coffee extraction. The tasting record currently captures grind size, dose, water volume, and brew time, but omits water temperature — leaving a gap in the brew recipe that makes it harder to reproduce or analyze results.

## What Changes

- Add an optional, nullable `water_temperature` field (integer, °C) to the `TastingNote` entity alongside the existing brew recipe fields
- Extend the tasting list and upsert API payloads to include `water_temperature`
- Add a `water_temperature` input to the tasting form with "I don't know" toggle support
- Display `water_temperature` on the tasting card, showing "IDK" when null

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `coffee-tasting`: The tasting note entity gains a `water_temperature` field (optional, nullable, integer °C). The form, card display, and "IDK" null-display behavior extend to cover this new field.

## Impact

- **API**: `api/schema/coffee.go` — add `WaterTemperature` column to `TastingNote` struct; GORM auto-migration adds the column. No handler changes needed (existing upsert/list logic is field-agnostic).
- **Vue types**: `vue-web/src/interfaces/Coffee.d.ts` — add `water_temperature` to `TastingNote` interface.
- **Vue components**: `TastingForm.vue` — add numeric input with IDK toggle. `TastingCard.vue` — render the value or "IDK".
- **Database**: SQLite auto-migration adds a nullable integer column; no data migration required.
