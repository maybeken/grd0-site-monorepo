# Coffee Tasting Editor

## Purpose

Provide a tabbed editor interface for managing coffee beans, brewing equipment, and tasting notes.

## Requirements

### Requirement: Coffee section navigation

The coffee editor SHALL be reached through the editor console as three sections with listings at `/editor/coffee/beans`, `/editor/coffee/equipment`, and `/editor/coffee/tastings`. Each section SHALL provide create, edit, and delete controls, and editing SHALL navigate to that section's identifier-based edit route. Entering `/editor/coffee` SHALL show the Beans listing.

#### Scenario: Default section selection
- **WHEN** user navigates to `/editor/coffee`
- **THEN** the Beans listing is shown

#### Scenario: Section switching
- **WHEN** user selects a different coffee section in the editor sidebar
- **THEN** that section's listing is displayed

#### Scenario: Edit route
- **WHEN** user activates the edit control on an item
- **THEN** the editor SHALL navigate to that section's edit route for the item's identifier

#### Scenario: Direct edit link
- **WHEN** user opens an item's edit route directly
- **THEN** the editor SHALL load that item without visiting the listing first

### Requirement: Bean management

The system SHALL provide bean management in the Beans section that displays a list of all beans and allows creating, editing, and deleting beans. Creating a bean SHALL use an inline form on the listing; editing a bean SHALL open the bean's edit route.

#### Scenario: Display bean list
- **WHEN** user views the Beans section
- **THEN** system displays all beans with name, roaster, and origin information

#### Scenario: Create new bean
- **WHEN** user clicks "+ New Bean" button
- **THEN** system displays a bean form with fields for name, origin, roaster, roast date, variety, process, altitude, and description

#### Scenario: Edit existing bean
- **WHEN** user clicks "Edit" button on a bean item
- **THEN** system navigates to the bean's edit route
- **THEN** system displays the bean form pre-populated with the bean's data

#### Scenario: Unknown bean
- **WHEN** user opens a bean edit route whose identifier does not exist
- **THEN** system displays a not-found message and a link back to the Beans listing

#### Scenario: Delete bean
- **WHEN** user clicks "Delete" button on a bean item
- **THEN** system shows a confirmation dialog and deletes the bean upon confirmation

#### Scenario: Bean form validation
- **WHEN** user attempts to save a bean without a name
- **THEN** the save button is disabled

#### Scenario: Bean saved successfully
- **WHEN** user saves a valid bean from the create form
- **THEN** system closes the form and refreshes the bean list

#### Scenario: Bean edit saved
- **WHEN** user saves a valid bean from its edit route
- **THEN** system reports success and remains on the edit route

### Requirement: Equipment management

The system SHALL provide equipment management in the Equipment section that displays a list of all equipment and allows creating, editing, and deleting equipment. Creating equipment SHALL use an inline form on the listing; editing equipment SHALL open the equipment's edit route.

#### Scenario: Display equipment list
- **WHEN** user views the Equipment section
- **THEN** system displays all equipment with name and type information

#### Scenario: Create new equipment
- **WHEN** user clicks "+ New Equipment" button
- **THEN** system displays an equipment form with fields for name, type, and description

#### Scenario: Edit existing equipment
- **WHEN** user clicks "Edit" button on an equipment item
- **THEN** system navigates to the equipment's edit route
- **THEN** system displays the equipment form pre-populated with the equipment's data

#### Scenario: Unknown equipment
- **WHEN** user opens an equipment edit route whose identifier does not exist
- **THEN** system displays a not-found message and a link back to the Equipment listing

#### Scenario: Delete equipment
- **WHEN** user clicks "Delete" button on an equipment item
- **THEN** system shows a confirmation dialog and deletes the equipment upon confirmation

#### Scenario: Equipment form validation
- **WHEN** user attempts to save equipment without a name
- **THEN** the save button is disabled

#### Scenario: Equipment saved successfully
- **WHEN** user saves valid equipment from the create form
- **THEN** system closes the form and refreshes the equipment list

#### Scenario: Equipment edit saved
- **WHEN** user saves valid equipment from its edit route
- **THEN** system reports success and remains on the edit route

### Requirement: Tasting management

The system SHALL provide tasting management in the Tastings section that displays a list of all tastings and allows creating, editing, and deleting tastings. Creating a tasting SHALL use an inline form on the listing; editing a tasting SHALL open the tasting's edit route with its bean and equipment options available.

#### Scenario: Display tasting list
- **WHEN** user views the Tastings section
- **THEN** system displays all tastings with bean name, equipment name, date, and rating

#### Scenario: Create new tasting
- **WHEN** user clicks "+ New Tasting" button
- **THEN** system displays the tasting form

#### Scenario: Edit existing tasting
- **WHEN** user clicks "Edit" button on a tasting item
- **THEN** system navigates to the tasting's edit route
- **THEN** system displays the tasting form pre-populated with the tasting's data

#### Scenario: Edit route options
- **WHEN** a tasting edit route loads
- **THEN** the form SHALL offer the current beans and equipment as selections

#### Scenario: Unknown tasting
- **WHEN** user opens a tasting edit route whose identifier does not exist
- **THEN** system displays a not-found message and a link back to the Tastings listing

#### Scenario: Delete tasting
- **WHEN** user clicks "Delete" button on a tasting item
- **THEN** system shows a confirmation dialog and deletes the tasting upon confirmation

#### Scenario: Tasting saved successfully
- **WHEN** user saves a valid tasting from the create form
- **THEN** system closes the form and refreshes the tasting list

#### Scenario: Tasting edit saved
- **WHEN** user saves a valid tasting from its edit route
- **THEN** system reports success and remains on the edit route

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
