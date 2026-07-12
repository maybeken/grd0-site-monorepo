# Brew Equipment

## Purpose

Manage brew equipment entities used in coffee tasting logs. Equipment represents the brewing device (e.g., pour over dripper, espresso machine) associated with each tasting note.

## Requirements

### Requirement: Brew equipment entity
The system SHALL provide a `BrewEquipment` entity with the following fields:
- `id` (UUID, primary key, auto-generated)
- `name` (text, required, not null)
- `type` (text, optional, from constant array)
- `description` (text, optional)
- `created_at`, `updated_at`, `deleted_at` (timestamps, auto-managed)

The `type` field SHALL be one of the following constant values:
`["Pour Over", "Espresso", "Immersion", "AeroPress", "French Press", "Moka Pot", "Cold Brew", "Other"]`

All fields except `name` SHALL support null values to represent "I don't know".

#### Scenario: Create equipment with name and type
- **WHEN** admin creates equipment with `name = "Hario V60"` and `type = "Pour Over"`
- **THEN** system creates the equipment record with UUID primary key

#### Scenario: Create equipment with only name
- **WHEN** admin creates equipment with `name = "Unknown Dripper"` and all other fields null
- **THEN** system creates the equipment record with null optional fields

#### Scenario: Soft delete equipment
- **WHEN** admin deletes equipment
- **THEN** system sets `deleted_at` timestamp (soft delete), equipment is no longer returned in list queries

### Requirement: Equipment list API
The system SHALL provide a public GET endpoint `/coffee/equipment` that returns all non-deleted equipment.

#### Scenario: List equipment
- **WHEN** any user requests GET `/coffee/equipment`
- **THEN** system returns a JSON array of all non-deleted equipment ordered by `created_at DESC`

### Requirement: Equipment upsert API
The system SHALL provide an admin-protected PUT endpoint `/coffee/equipment` that creates or updates equipment.

#### Scenario: Create new equipment via upsert
- **WHEN** admin sends PUT `/coffee/equipment` with an equipment object without an existing ID
- **THEN** system creates a new equipment record and returns it

#### Scenario: Update existing equipment via upsert
- **WHEN** admin sends PUT `/coffee/equipment` with an equipment object containing an existing ID
- **THEN** system updates the existing equipment record and returns it

### Requirement: Equipment delete API
The system SHALL provide an admin-protected DELETE endpoint `/coffee/equipment/:id` that soft-deletes equipment.

#### Scenario: Delete existing equipment
- **WHEN** admin sends DELETE `/coffee/equipment/:id` with a valid equipment ID
- **THEN** system soft-deletes the equipment and returns 202 Accepted

#### Scenario: Delete non-existent equipment
- **WHEN** admin sends DELETE `/coffee/equipment/:id` with an ID that does not exist
- **THEN** system returns 404 Not Found
