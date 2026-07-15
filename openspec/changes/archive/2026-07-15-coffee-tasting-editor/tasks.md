## 1. Bean Form Component

- [x] 1.1 Create `vue-web/src/components/coffee/BeanForm.vue` with form fields for name, origin, roaster, roast date, variety, process, altitude, description
- [x] 1.2 Implement create/edit mode via optional `bean` prop
- [x] 1.3 Implement save via `upsertCoffeeBeanRaw` with roast date ISO conversion
- [x] 1.4 Emit `saved` and `cancel` events
- [x] 1.5 Disable save button when name is empty

## 2. Equipment Form Component

- [x] 2.1 Create `vue-web/src/components/coffee/EquipmentForm.vue` with form fields for name, type, description
- [x] 2.2 Populate type dropdown with `EQUIPMENT_TYPES` plus "Grinder"
- [x] 2.3 Implement create/edit mode via optional `equipment` prop
- [x] 2.4 Implement save via `upsertCoffeeEquipmentRaw`
- [x] 2.5 Emit `saved` and `cancel` events
- [x] 2.6 Disable save button when name is empty

## 3. Editor View Rewrite

- [x] 3.1 Rewrite `CoffeeEditorView.vue` with tabbed interface (Beans, Equipment, Tasting)
- [x] 3.2 Add bean list view with edit/delete buttons and "+ New Bean" button
- [x] 3.3 Add equipment list view with edit/delete buttons and "+ New Equipment" button
- [x] 3.4 Add tasting list view with edit/delete buttons and "+ New Tasting" button
- [x] 3.5 Implement fetch functions using `getCoffeeBeansRaw`, `getCoffeeEquipmentRaw`, `getCoffeeTastingsRaw`
- [x] 3.6 Implement delete handlers with confirmation dialogs using `deleteCoffeeBeanRaw`, `deleteCoffeeEquipmentRaw`, `deleteCoffeeTastingRaw`
- [x] 3.7 Refresh lists after save/delete operations

## 4. Tasting Form Quick-Add

- [x] 4.1 Add "+ new" button next to bean select in `TastingForm.vue`
- [x] 4.2 Add inline quick-add form for beans (name, roaster, origin)
- [x] 4.3 Implement `saveQuickBean` using `upsertCoffeeBeanRaw` with auto-select
- [x] 4.4 Add "+ new" button next to equipment select in `TastingForm.vue`
- [x] 4.5 Add inline quick-add form for equipment (name, type)
- [x] 4.6 Implement `saveQuickEquipment` using `upsertCoffeeEquipmentRaw` with auto-select
- [x] 4.7 Emit `beanAdded` and `equipmentAdded` events from `TastingForm`
- [x] 4.8 Handle `beanAdded` and `equipmentAdded` events in `CoffeeEditorView` to refresh lists

## 5. Verification

- [x] 5.1 Run `cd vue-web && pnpm type-check` to verify no type errors
- [x] 5.2 Run `cd vue-web && pnpm format` to apply code formatting
