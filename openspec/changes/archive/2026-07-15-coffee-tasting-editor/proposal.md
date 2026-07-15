## Why

The coffee tasting editor (`/coffee/editor`) previously only supported creating new tasting notes. There was no UI for managing the underlying beans and equipment that tastings reference. Users had to interact with the API directly to add, edit, or delete beans and equipment, which is impractical for daily use. Additionally, when creating a tasting, if the desired bean or equipment didn't exist yet, the user had to navigate away to another tool or tab.

## What Changes

- Add a tabbed editor interface at `/coffee/editor` with "Beans", "Equipment", and "Tasting" tabs
- Add list views for beans and equipment with inline edit and delete controls
- Create reusable `BeanForm` and `EquipmentForm` components for create/edit operations
- Add inline quick-add forms inside `TastingForm` for creating beans and equipment without leaving the tasting form
- Auto-select newly created beans/equipment in the tasting form after quick-add
- Refresh parent lists when quick-add creates new items

## Capabilities

### New Capabilities

- `coffee-tasting-editor`: Admin interface for managing coffee beans, brew equipment, and tasting notes with tabbed navigation, list views with CRUD controls, and inline quick-add for beans/equipment within the tasting form.

### Modified Capabilities

- `coffee-tasting`: The tasting form gains inline quick-add capability for beans and equipment, emitting `beanAdded` and `equipmentAdded` events to notify the parent when new items are created.

## Impact

- **Frontend components**: New `BeanForm.vue` and `EquipmentForm.vue` in `vue-web/src/components/coffee/`
- **Frontend views**: `CoffeeEditorView.vue` rewritten with tabbed interface and list management
- **Frontend components**: `TastingForm.vue` extended with quick-add inline forms and new emits
- **API**: No changes — uses existing `upsertCoffeeBeanRaw`, `upsertCoffeeEquipmentRaw`, `deleteCoffeeBeanRaw`, `deleteCoffeeEquipmentRaw`, `getCoffeeTastingsRaw`
- **Types**: No changes — existing `CoffeeBean`, `BrewEquipment`, `TastingNote` interfaces suffice
