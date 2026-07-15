## MODIFIED Requirements

### Requirement: Inline quick-add for beans in tasting form
The system SHALL provide a "+ new" button next to the bean select dropdown in the tasting form. Clicking this button SHALL display an inline form to quickly create a new bean without leaving the tasting form.

#### Scenario: Show quick-add bean form
- **WHEN** user clicks "+ new" button next to bean select
- **THEN** system displays an inline form with fields for bean name, roaster, and origin

#### Scenario: Quick-add bean with minimal fields
- **WHEN** user enters only a bean name and clicks "Add"
- **THEN** system creates the bean, auto-selects it in the tasting form, and closes the inline form

#### Scenario: Quick-add bean with all fields
- **WHEN** user enters bean name, roaster, and origin, then clicks "Add"
- **THEN** system creates the bean with all provided fields, auto-selects it, and closes the inline form

#### Scenario: Cancel quick-add bean
- **WHEN** user clicks "Cancel" in the inline bean form
- **THEN** system closes the inline form without creating a bean

#### Scenario: Quick-add bean validation
- **WHEN** user attempts to add a bean without a name
- **THEN** the "Add" button is disabled

#### Scenario: Bean list refreshes after quick-add
- **WHEN** a new bean is created via quick-add
- **THEN** the bean list in the parent view is refreshed to include the new bean

### Requirement: Inline quick-add for equipment in tasting form
The system SHALL provide a "+ new" button next to the equipment select dropdown in the tasting form. Clicking this button SHALL display an inline form to quickly create new equipment without leaving the tasting form.

#### Scenario: Show quick-add equipment form
- **WHEN** user clicks "+ new" button next to equipment select
- **THEN** system displays an inline form with fields for equipment name and type

#### Scenario: Quick-add equipment with minimal fields
- **WHEN** user enters only an equipment name and clicks "Add"
- **THEN** system creates the equipment, auto-selects it in the tasting form, and closes the inline form

#### Scenario: Quick-add equipment with type
- **WHEN** user enters equipment name and selects a type, then clicks "Add"
- **THEN** system creates the equipment with the name and type, auto-selects it, and closes the inline form

#### Scenario: Cancel quick-add equipment
- **WHEN** user clicks "Cancel" in the inline equipment form
- **THEN** system closes the inline form without creating equipment

#### Scenario: Quick-add equipment validation
- **WHEN** user attempts to add equipment without a name
- **THEN** the "Add" button is disabled

#### Scenario: Equipment list refreshes after quick-add
- **WHEN** new equipment is created via quick-add
- **THEN** the equipment list in the parent view is refreshed to include the new equipment
