# Spec Delta

## ADDED Requirements

### Requirement: Tasting detail API

The system SHALL provide a public GET endpoint `/coffee/tasting/:id` that returns one non-deleted tasting with its bean and equipment details nested.

#### Scenario: Read one tasting
- **WHEN** any user requests GET `/coffee/tasting/:id` with a valid tasting ID
- **THEN** system returns the tasting as a JSON object with nested bean and equipment objects
- **THEN** a grinder reference, when present, SHALL be returned as a nested equipment object

#### Scenario: Read an unknown tasting
- **WHEN** any user requests GET `/coffee/tasting/:id` with an ID that does not exist
- **THEN** system returns 404 Not Found

#### Scenario: Read a soft-deleted tasting
- **WHEN** any user requests GET `/coffee/tasting/:id` for a tasting that has been soft-deleted
- **THEN** system returns 404 Not Found
