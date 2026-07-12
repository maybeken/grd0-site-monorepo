## 1. Database & API Schema

- [x] 1.1 Create `api/schema/coffee.go` with `CoffeeBean`, `BrewEquipment`, and `TastingNote` structs embedding `BaseColumns`
- [x] 1.2 Add all three schemas to `api/database/migration.go` tables array

## 2. API Handlers

- [x] 2.1 Create `api/handler/coffee.go` with handlers: `GetCoffeeBeans`, `UpsertCoffeeBean`, `DeleteCoffeeBean`
- [x] 2.2 Add handlers: `GetCoffeeEquipment`, `UpsertCoffeeEquipment`, `DeleteCoffeeEquipment`
- [x] 2.3 Add handlers: `GetCoffeeTastings`, `UpsertCoffeeTasting`, `DeleteCoffeeTasting`
- [x] 2.4 Register all coffee routes in `api/handler/handler.go` `RegisterRouter` function (public GET, admin PUT/DELETE)

## 3. Frontend Dependencies & Types

- [x] 3.1 Install `chart.js` and `vue-chartjs` in `vue-web/`
- [x] 3.2 Create `vue-web/src/interfaces/Coffee.d.ts` with `CoffeeBean`, `BrewEquipment`, `TastingNote` TypeScript interfaces
- [x] 3.3 Create `vue-web/src/services/coffee.ts` with API service functions using `dataInstance` and `adminInstance`

## 4. Frontend Constants & Assets

- [x] 4.1 Create `vue-web/src/helpers/coffee.ts` with constant arrays: `EQUIPMENT_TYPES`, `GRIND_SIZES`, `EQUIPMENT_IMAGE_MAP`
- [x] 4.2 Add placeholder SVG equipment images to `vue-web/src/assets/coffee/equipment/` (user will replace manually)

## 5. Frontend Components

- [x] 5.1 Create `vue-web/src/components/coffee/TasteRadarChart.vue` — Chart.js radar chart component for 7 taste dimensions
- [x] 5.2 Create `vue-web/src/components/coffee/TastingCard.vue` — 1:1 square card with equipment image, bean/brew info, radar chart, notes
- [x] 5.3 Create `vue-web/src/components/coffee/IdkToggle.vue` — reusable "I don't know" checkbox toggle component
- [x] 5.4 Create `vue-web/src/components/coffee/BrewTimeInput.vue` — mm:ss format input with conversion to seconds
- [x] 5.5 Create `vue-web/src/components/coffee/TastingForm.vue` — admin form for creating/editing tastings with all fields and IDK toggles
- [x] 5.6 Implement auto-calculation of `ratio` from `water_in / coffee_dose` in TastingForm (reactive update when inputs change)

## 6. Frontend Views & Routing

- [x] 6.1 Create `vue-web/src/views/CoffeeView.vue` — public tasting log with card grid, sorting, loading/empty states
- [x] 6.2 Add `/coffee` route to `vue-web/src/router/index.ts`

## 7. Verification

- [x] 7.1 Run `cd vue-web && pnpm type-check` to verify no type errors
- [x] 7.2 Run `cd vue-web && pnpm lint` to verify no lint errors
- [x] 7.3 Run `cd api && go build -o ../grd0_blog_api .` to verify API compiles

## 8. Testing

- [x] 8.1 Install Playwright for headless browser testing: `cd vue-web && pnpm add -D @playwright/test`
- [x] 8.2 Create `vue-web/tests/coffee-page.spec.ts` to verify the coffee page layout matches spec
- [x] 8.3 Write Playwright test to verify tasting cards render with 1:1 aspect ratio
- [x] 8.4 Write Playwright test to verify radar chart displays 7 taste dimensions
- [x] 8.5 Write Playwright test to verify pinned tastings appear first in the list
- [x] 8.6 Write Playwright test to verify "IDK" displays for null fields
- [x] 8.7 Write Playwright test to verify equipment images render correctly
- [x] 8.8 Run `cd vue-web && pnpm exec playwright test` to execute visual verification tests

## 9. Documentation

- [x] 9.1 Create `docs/api/coffee.json` in Restfox collection format documenting all coffee API endpoints (GET/PUT/DELETE for beans, equipment, tastings)
