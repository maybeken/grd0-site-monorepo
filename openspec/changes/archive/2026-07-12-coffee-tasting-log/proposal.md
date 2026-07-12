## Why

Add a coffee tasting and brewing journal to track brew recipes, bean details, and taste profiles. This allows the site owner to record and review their coffee brewing history, with public visibility for site visitors.

## What Changes

- Add `/coffee` top-level route (public) displaying tasting log as 1:1 square cards
- Add three new database entities: `CoffeeBean`, `BrewEquipment`, `TastingNote`
- Add admin API endpoints for CRUD operations on beans, equipment, and tastings
- Add public API endpoint for listing tastings (with bean/equipment details)
- Add Vue components for tasting card display with radar chart visualization
- Add admin form for creating/editing tastings with "I don't know" toggles
- Add local static equipment images with mapping array configuration
- Install `chart.js` and `vue-chartjs` for radar chart rendering

## Capabilities

### New Capabilities

- `coffee-bean`: Coffee bean entity with name (required), origin, roaster, roast date, variety, process, altitude, description. All fields except name support "I don't know" (nullable).
- `brew-equipment`: Brew equipment entity with name, type (dropdown from constant array), description. Supports "I don't know" for optional fields.
- `coffee-tasting`: Tasting note entity linking bean and equipment, with brew recipe (grind size, grind setting, coffee dose, water in, coffee out, auto-calculated ratio, brew time), taste profile (7 dimensions on 0-10 scale), overall notes, rating (1-10), date, and pinned flag. All brew/taste fields support "I don't know" toggles.
- `coffee-list-view`: Public `/coffee` route displaying tasting cards in 1:1 square layout. Each card shows equipment image, bean/brew info (left), radar chart (right), and notes. Sorted by pinned status then date descending.

### Modified Capabilities

(none)

## Impact

- **Database**: New tables `coffee_beans`, `brew_equipment`, `tasting_notes` with UUID primary keys
- **API**: New routes under `/coffee/*` (public GET, admin PUT/DELETE)
- **Frontend**: New view `CoffeeView.vue`, components for card/form, service layer for API calls
- **Dependencies**: Add `chart.js` and `vue-chartjs` to vue-web
- **Assets**: Local static images for equipment types in `vue-web/src/assets/coffee/`
- **Router**: New `/coffee` route in `vue-web/src/router/index.ts`
