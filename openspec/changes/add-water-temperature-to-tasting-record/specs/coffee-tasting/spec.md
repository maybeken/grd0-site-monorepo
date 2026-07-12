## MODIFIED Requirements

### Requirement: Tasting note entity
The system SHALL provide a `TastingNote` entity with the following fields:
- `id` (UUID, primary key, auto-generated)
- `bean_id` (UUID, foreign key to coffee_bean)
- `equipment_id` (UUID, foreign key to brew_equipment)

Brew recipe fields (all optional, nullable):
- `grind_size` (text, from constant array)
- `grind_setting` (real, decimal number)
- `coffee_dose` (integer, grams)
- `water_in` (integer, ml)
- `coffee_out` (integer, ml, yield)
- `ratio` (real, auto-calculated as water_in / coffee_dose)
- `brew_time` (integer, seconds)
- `water_temperature` (integer, °C)

Taste profile fields (all optional, nullable, 0-10 scale):
- `taste_fruity`, `taste_sour`, `taste_sweetness`, `taste_nutty`, `taste_spice`, `taste_floral`, `taste_green` (integer)

Overall fields:
- `overall_notes` (text, optional)
- `rating` (integer, 1-10, optional)
- `tasted_at` (date, required)
- `pinned` (boolean, default false)

- `created_at`, `updated_at`, `deleted_at` (timestamps, auto-managed)

The `grind_size` field SHALL be one of the following constant values:
`["Extra Fine", "Fine", "Medium-Fine", "Medium", "Medium-Coarse", "Coarse", "Extra Coarse"]`

All brew recipe and taste profile fields SHALL support null values to represent "I don't know".

#### Scenario: Create tasting with all fields
- **WHEN** admin creates a tasting with all fields populated including ratio and water_temperature
- **THEN** system creates the record with all values stored as provided

#### Scenario: Create tasting with minimal fields
- **WHEN** admin creates a tasting with only `bean_id`, `equipment_id`, and `tasted_at`
- **THEN** system creates the record with all optional fields null

#### Scenario: Ratio stored as provided
- **WHEN** admin creates a tasting with `water_in = 270`, `coffee_dose = 18`, and `ratio = 15.0`
- **THEN** system stores `ratio = 15.0` as provided (frontend is responsible for calculation)

#### Scenario: Brew time stored as seconds
- **WHEN** admin creates a tasting with brew time 3 minutes 30 seconds
- **THEN** system stores `brew_time = 210` (3 * 60 + 30)

#### Scenario: Water temperature stored as integer
- **WHEN** admin creates a tasting with `water_temperature = 93`
- **THEN** system stores `water_temperature = 93`

#### Scenario: Water temperature null when unknown
- **WHEN** admin creates a tasting with `water_temperature = null` (IDK toggle checked)
- **THEN** system stores `water_temperature = null`

#### Scenario: Taste profile values in range
- **WHEN** admin sets `taste_fruity = 7`
- **THEN** system stores the value 7
