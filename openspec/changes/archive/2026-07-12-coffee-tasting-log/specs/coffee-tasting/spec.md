## ADDED Requirements

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
- **WHEN** admin creates a tasting with all fields populated including ratio
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

#### Scenario: Taste profile values in range
- **WHEN** admin sets `taste_fruity = 7`
- **THEN** system stores the value 7

### Requirement: Tasting list API
The system SHALL provide a public GET endpoint `/coffee/tastings` that returns all non-deleted tastings with joined bean and equipment details.

Results SHALL be ordered by `pinned DESC`, then `tasted_at DESC`.

#### Scenario: List tastings
- **WHEN** any user requests GET `/coffee/tastings`
- **THEN** system returns a JSON array of tastings with nested bean and equipment objects, ordered by pinned status then date

#### Scenario: Pinned tastings appear first
- **WHEN** tasting A is pinned (tasted Jul 5) and tasting B is not pinned (tasted Jul 10)
- **THEN** tasting A appears before tasting B in the list

#### Scenario: Tastings include bean and equipment details
- **WHEN** a tasting references a bean and equipment
- **THEN** the response includes the full bean and equipment objects nested within the tasting

### Requirement: Tasting upsert API
The system SHALL provide an admin-protected PUT endpoint `/coffee/tasting` that creates or updates a tasting note.

#### Scenario: Create new tasting via upsert
- **WHEN** admin sends PUT `/coffee/tasting` with a tasting object without an existing ID
- **THEN** system creates a new tasting record and returns it

#### Scenario: Update existing tasting via upsert
- **WHEN** admin sends PUT `/coffee/tasting` with a tasting object containing an existing ID
- **THEN** system updates the existing tasting record and returns it

### Requirement: Tasting delete API
The system SHALL provide an admin-protected DELETE endpoint `/coffee/tasting/:id` that soft-deletes a tasting.

#### Scenario: Delete existing tasting
- **WHEN** admin sends DELETE `/coffee/tasting/:id` with a valid tasting ID
- **THEN** system soft-deletes the tasting and returns 202 Accepted

#### Scenario: Delete non-existent tasting
- **WHEN** admin sends DELETE `/coffee/tasting/:id` with an ID that does not exist
- **THEN** system returns 404 Not Found

### Requirement: "I don't know" toggle behavior
The frontend SHALL display a checkbox toggle labeled "I don't know" next to each optional field. When checked, the field SHALL be disabled and the value SHALL be sent as null. When unchecked, the field SHALL be enabled and the value SHALL be sent as entered.

#### Scenario: Toggle checked disables field
- **WHEN** user checks "I don't know" for grind_size
- **THEN** the grind_size input is disabled and the value is set to null

#### Scenario: Toggle unchecked enables field
- **WHEN** user unchecks "I don't know" for grind_size
- **THEN** the grind_size input is enabled and accepts user input

### Requirement: Brew time input format
The frontend SHALL accept brew time input in mm:ss format. The input SHALL only accept digits and a colon separator. On save, the value SHALL be converted to total seconds for API submission.

#### Scenario: Valid mm:ss input
- **WHEN** user enters "03:30" in brew time field
- **THEN** system converts to 210 seconds for API submission

#### Scenario: Invalid input rejected
- **WHEN** user enters "abc" in brew time field
- **THEN** system prevents submission or shows validation error

### Requirement: Display "IDK" for null values
The frontend SHALL display the literal text "IDK" for any field that has a null value in the tasting card view.

#### Scenario: Null field displays IDK
- **WHEN** a tasting has `grind_size = null`
- **THEN** the card displays "IDK" for the grind size field

#### Scenario: Non-null field displays value
- **WHEN** a tasting has `grind_size = "Medium-Fine"`
- **THEN** the card displays "Medium-Fine" for the grind size field

### Requirement: Auto-calculate water ratio on frontend
The frontend SHALL automatically calculate the water ratio as `water_in / coffee_dose` when both values are provided. The calculated ratio SHALL be included in the API request payload.

#### Scenario: Both inputs provided
- **WHEN** user enters `water_in = 270` and `coffee_dose = 18`
- **THEN** frontend calculates `ratio = 15.0` and includes it in the save request

#### Scenario: One input missing
- **WHEN** user enters `water_in = 270` but `coffee_dose` is null or "IDK"
- **THEN** frontend sends `ratio = null` in the save request

#### Scenario: Ratio updates reactively
- **WHEN** user changes `coffee_dose` from 18 to 20 while `water_in = 270`
- **THEN** frontend recalculates `ratio = 13.5` and updates the displayed value
