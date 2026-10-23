# Spec Delta

## ADDED Requirements

### Requirement: Bean detail API

The system SHALL provide a public GET endpoint `/coffee/bean/:id` that returns one non-deleted coffee bean.

#### Scenario: Read one bean
- **WHEN** any user requests GET `/coffee/bean/:id` with a valid bean ID
- **THEN** system returns the bean as a JSON object

#### Scenario: Read an unknown bean
- **WHEN** any user requests GET `/coffee/bean/:id` with an ID that does not exist
- **THEN** system returns 404 Not Found

#### Scenario: Read a soft-deleted bean
- **WHEN** any user requests GET `/coffee/bean/:id` for a bean that has been soft-deleted
- **THEN** system returns 404 Not Found
