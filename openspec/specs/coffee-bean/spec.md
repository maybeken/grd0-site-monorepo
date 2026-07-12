# Coffee Bean

## Purpose

Manage coffee bean entities used in coffee tasting logs. A bean represents a specific coffee with origin, roaster, and processing details that is referenced by tasting notes.

## Requirements

### Requirement: Coffee bean entity
The system SHALL provide a `CoffeeBean` entity with the following fields:
- `id` (UUID, primary key, auto-generated)
- `name` (text, required, not null)
- `origin` (text, optional)
- `roaster` (text, optional)
- `roast_date` (date, optional)
- `variety` (text, optional)
- `process` (text, optional)
- `altitude` (text, optional)
- `description` (text, optional)
- `created_at`, `updated_at`, `deleted_at` (timestamps, auto-managed)

All fields except `name` SHALL support null values to represent "I don't know".

#### Scenario: Create a bean with only name
- **WHEN** admin creates a coffee bean with `name = "Ethiopian Yirgacheffe"` and all other fields null
- **THEN** system creates the bean record with UUID primary key and null optional fields

#### Scenario: Create a bean with all fields
- **WHEN** admin creates a coffee bean with all fields populated
- **THEN** system creates the bean record with all values stored

#### Scenario: Bean name is required
- **WHEN** admin attempts to create a coffee bean without a name
- **THEN** system rejects the request with a validation error

#### Scenario: Soft delete a bean
- **WHEN** admin deletes a coffee bean
- **THEN** system sets `deleted_at` timestamp (soft delete), bean is no longer returned in list queries

### Requirement: Bean list API
The system SHALL provide a public GET endpoint `/coffee/beans` that returns all non-deleted coffee beans.

#### Scenario: List beans
- **WHEN** any user requests GET `/coffee/beans`
- **THEN** system returns a JSON array of all non-deleted beans ordered by `created_at DESC`

#### Scenario: List beans excludes soft-deleted
- **WHEN** a bean has been soft-deleted
- **THEN** that bean is not included in the list response

### Requirement: Bean upsert API
The system SHALL provide an admin-protected PUT endpoint `/coffee/bean` that creates or updates a coffee bean.

#### Scenario: Create new bean via upsert
- **WHEN** admin sends PUT `/coffee/bean` with a bean object without an existing ID
- **THEN** system creates a new bean record and returns it

#### Scenario: Update existing bean via upsert
- **WHEN** admin sends PUT `/coffee/bean` with a bean object containing an existing ID
- **THEN** system updates the existing bean record and returns it

### Requirement: Bean delete API
The system SHALL provide an admin-protected DELETE endpoint `/coffee/bean/:id` that soft-deletes a coffee bean.

#### Scenario: Delete existing bean
- **WHEN** admin sends DELETE `/coffee/bean/:id` with a valid bean ID
- **THEN** system soft-deletes the bean and returns 202 Accepted

#### Scenario: Delete non-existent bean
- **WHEN** admin sends DELETE `/coffee/bean/:id` with an ID that does not exist
- **THEN** system returns 404 Not Found
