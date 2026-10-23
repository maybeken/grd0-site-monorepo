# Spec Delta

## ADDED Requirements

### Requirement: Equipment detail API

The system SHALL provide a public GET endpoint `/coffee/equipment/:id` that returns one non-deleted brew equipment record.

#### Scenario: Read one equipment record
- **WHEN** any user requests GET `/coffee/equipment/:id` with a valid equipment ID
- **THEN** system returns the equipment as a JSON object

#### Scenario: Read unknown equipment
- **WHEN** any user requests GET `/coffee/equipment/:id` with an ID that does not exist
- **THEN** system returns 404 Not Found

#### Scenario: Read soft-deleted equipment
- **WHEN** any user requests GET `/coffee/equipment/:id` for equipment that has been soft-deleted
- **THEN** system returns 404 Not Found
