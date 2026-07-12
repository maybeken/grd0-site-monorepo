## 1. Backend Schema Update

- [x] 1.1 Add `WaterTemperature *int` field to `TastingNote` struct in `api/schema/coffee.go` with JSON tag `water_temperature`, placed after `BrewTime`
- [x] 1.2 Verify GORM auto-migration adds the nullable column (no manual migration needed)

## 2. Frontend Type Definition

- [x] 2.1 Add `water_temperature?: number | null` to `TastingNote` interface in `vue-web/src/interfaces/Coffee.d.ts`

## 3. Frontend Form Component

- [x] 3.1 Add water temperature numeric input field to `TastingForm.vue` in the brew recipe section (after brew time)
- [x] 3.2 Wire IDK toggle for water temperature field following existing pattern
- [x] 3.3 Ensure water temperature value is included in the upsert payload

## 4. Frontend Card Display

- [x] 4.1 Add water temperature display to `TastingCard.vue` showing value with "°C" suffix
- [x] 4.2 Display "IDK" when water temperature is null, following existing pattern

## 5. Verification

- [x] 5.1 Test creating a tasting with water temperature value via form
- [x] 5.2 Test creating a tasting with water temperature set to IDK (null)
- [x] 5.3 Verify existing tastings display "IDK" for water temperature (null values)
- [x] 5.4 Verify API returns water_temperature field in tasting list and upsert responses
