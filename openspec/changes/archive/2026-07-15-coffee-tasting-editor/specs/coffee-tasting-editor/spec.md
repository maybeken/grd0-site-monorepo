## ADDED Requirements

### Requirement: Tabbed editor interface
The system SHALL provide a tabbed editor view at `/coffee/editor` with three tabs: "Beans", "Equipment", and "Tasting". Each tab SHALL display a list of existing items and provide controls to add, edit, or delete items.

#### Scenario: Default tab selection
- **WHEN** user navigates to the editor view
- **THEN** the "Beans" tab is selected by default

#### Scenario: Tab switching
- **WHEN** user clicks on a different tab
- **THEN** the corresponding tab content is displayed

### Requirement: Bean management
The system SHALL provide a bean management interface within the "Beans" tab that displays a list of all beans and allows creating, editing, and deleting beans.

#### Scenario: Display bean list
- **WHEN** user views the "Beans" tab
- **THEN** system displays all beans with name, roaster, and origin information

#### Scenario: Create new bean
- **WHEN** user clicks "+ New Bean" button
- **THEN** system displays a bean form with fields for name, origin, roaster, roast date, variety, process, altitude, and description

#### Scenario: Edit existing bean
- **WHEN** user clicks "Edit" button on a bean item
- **THEN** system displays the bean form pre-populated with the bean's data

#### Scenario: Delete bean
- **WHEN** user clicks "Delete" button on a bean item
- **THEN** system shows a confirmation dialog and deletes the bean upon confirmation

#### Scenario: Bean form validation
- **WHEN** user attempts to save a bean without a name
- **THEN** the save button is disabled

#### Scenario: Bean saved successfully
- **WHEN** user saves a valid bean
- **THEN** system closes the form and refreshes the bean list

### Requirement: Equipment management
The system SHALL provide an equipment management interface within the "Equipment" tab that displays a list of all equipment and allows creating, editing, and deleting equipment.

#### Scenario: Display equipment list
- **WHEN** user views the "Equipment" tab
- **THEN** system displays all equipment with name and type information

#### Scenario: Create new equipment
- **WHEN** user clicks "+ New Equipment" button
- **THEN** system displays an equipment form with fields for name, type, and description

#### Scenario: Edit existing equipment
- **WHEN** user clicks "Edit" button on an equipment item
- **THEN** system displays the equipment form pre-populated with the equipment's data

#### Scenario: Delete equipment
- **WHEN** user clicks "Delete" button on an equipment item
- **THEN** system shows a confirmation dialog and deletes the equipment upon confirmation

#### Scenario: Equipment form validation
- **WHEN** user attempts to save equipment without a name
- **THEN** the save button is disabled

#### Scenario: Equipment saved successfully
- **WHEN** user saves valid equipment
- **THEN** system closes the form and refreshes the equipment list

### Requirement: Tasting management
The system SHALL provide a tasting management interface within the "Tasting" tab that displays a list of all tastings and allows creating, editing, and deleting tastings.

#### Scenario: Display tasting list
- **WHEN** user views the "Tasting" tab
- **THEN** system displays all tastings with bean name, equipment name, date, and rating

#### Scenario: Create new tasting
- **WHEN** user clicks "+ New Tasting" button
- **THEN** system displays the tasting form

#### Scenario: Edit existing tasting
- **WHEN** user clicks "Edit" button on a tasting item
- **THEN** system displays the tasting form pre-populated with the tasting's data

#### Scenario: Delete tasting
- **WHEN** user clicks "Delete" button on a tasting item
- **THEN** system shows a confirmation dialog and deletes the tasting upon confirmation

#### Scenario: Tasting saved successfully
- **WHEN** user saves a valid tasting
- **THEN** system closes the form and refreshes the tasting list

### Requirement: Bean form component
The system SHALL provide a reusable `BeanForm` component that accepts an optional bean prop for editing and emits `saved` and `cancel` events.

#### Scenario: Bean form in create mode
- **WHEN** `BeanForm` is rendered without a bean prop
- **THEN** form displays with empty fields and title "New Bean"

#### Scenario: Bean form in edit mode
- **WHEN** `BeanForm` is rendered with a bean prop
- **THEN** form displays pre-populated with the bean's data and title "Edit Bean"

#### Scenario: Bean form save event
- **WHEN** user successfully saves a bean
- **THEN** component emits `saved` event

#### Scenario: Bean form cancel event
- **WHEN** user clicks "Cancel" button
- **THEN** component emits `cancel` event

### Requirement: Equipment form component
The system SHALL provide a reusable `EquipmentForm` component that accepts an optional equipment prop for editing and emits `saved` and `cancel` events.

#### Scenario: Equipment form in create mode
- **WHEN** `EquipmentForm` is rendered without an equipment prop
- **THEN** form displays with empty fields and title "New Equipment"

#### Scenario: Equipment form in edit mode
- **WHEN** `EquipmentForm` is rendered with an equipment prop
- **THEN** form displays pre-populated with the equipment's data and title "Edit Equipment"

#### Scenario: Equipment form save event
- **WHEN** user successfully saves equipment
- **THEN** component emits `saved` event

#### Scenario: Equipment form cancel event
- **WHEN** user clicks "Cancel" button
- **THEN** component emits `cancel` event

### Requirement: Equipment type options
The system SHALL provide equipment type options including all standard brewing methods plus "Grinder" for the equipment form and quick-add form.

#### Scenario: Equipment type dropdown options
- **WHEN** user views the equipment type dropdown
- **THEN** system displays options: "Pour Over", "Espresso", "French Press", "Moka Pot", "Syphon", "Other", and "Grinder"

### Requirement: Roast date handling
The system SHALL handle roast date input as a date picker and convert it to ISO format for API submission.

#### Scenario: Roast date input
- **WHEN** user selects a roast date in the bean form
- **THEN** system stores the date and converts to ISO format on save

#### Scenario: Roast date display in edit mode
- **WHEN** editing a bean with an existing roast date
- **THEN** system displays the date in the date picker (YYYY-MM-DD format)
