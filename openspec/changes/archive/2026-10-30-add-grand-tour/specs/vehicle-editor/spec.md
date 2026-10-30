# Spec Delta

## Purpose

Gives the authenticated author a place to create, edit and remove vehicle entries and to manage each vehicle's cover assets without hand-editing data.

## ADDED Requirements

### Requirement: Vehicle listing

The editor SHALL provide a listing at `/editor/vehicles` showing every stored vehicle with its cover thumbnail beside its name, its identity and its cover state, and SHALL allow deletion with confirmation.

#### Scenario: Listing entries
- **WHEN** the author opens the vehicle listing
- **THEN** each vehicle appears with its cover thumbnail beside its name — the cutout, the silhouette, or the "Unknown" placeholder — along with make, model, year and whether a cover is assigned

#### Scenario: Delete with confirmation
- **WHEN** the author deletes a vehicle and confirms
- **THEN** the vehicle is removed and the listing refreshes without it

#### Scenario: Delete cancelled
- **WHEN** the author dismisses the confirmation prompt
- **THEN** the vehicle remains and no request is sent

### Requirement: Vehicle workspace

The editor SHALL provide a workspace at `/editor/vehicles/:id` for creating and editing a single vehicle. Opening the route directly SHALL load the vehicle without visiting the listing first.

#### Scenario: Edit existing vehicle
- **WHEN** the author opens an existing vehicle's workspace
- **THEN** the form is populated with its stored values

#### Scenario: Invalid identifier
- **WHEN** the route identifier does not match a stored vehicle
- **THEN** the workspace reports that the vehicle could not be loaded and offers a way back to the listing

### Requirement: Field validation

The form SHALL require make, model and year, SHALL allow every other field to be empty, and MUST NOT save while a required field is missing.

#### Scenario: Save blocked
- **WHEN** the author submits with an empty model
- **THEN** the save is prevented and the missing field is identified

#### Scenario: Optional fields empty
- **WHEN** the author saves with only the required fields populated
- **THEN** the vehicle is stored successfully

### Requirement: Single country selection

The form SHALL present the country label and the ISO country code as one autocompleted field. Suggestions SHALL be drawn from the bundled flag collection, SHALL match on either country name or code, and SHALL show the matching flag. Selecting a suggestion SHALL store the country label and the upper-case ISO code together. Clearing the field SHALL clear both values. The field SHALL remain optional.

#### Scenario: Selecting a country
- **WHEN** the author types part of a country name or code and selects a suggestion
- **THEN** the form stores that country's label and ISO code together and shows the country's flag

#### Scenario: Clearing the country
- **WHEN** the author clears the field
- **THEN** both the label and the ISO code are removed and a subsequent save sends neither

#### Scenario: Leaving without choosing
- **WHEN** the author types text that matches no suggestion and leaves the field
- **THEN** the field reverts to the stored country and no unmatched label is saved

### Requirement: Cover preview

The workspace SHALL preview the vehicle's current cover and SHALL indicate when no generated cover exists.

#### Scenario: Generated cover present
- **WHEN** the vehicle has a cutout or silhouette asset
- **THEN** the workspace shows the cutout, or the silhouette when only that exists, as the preview

#### Scenario: No generated cover
- **WHEN** the vehicle has neither a cutout nor a silhouette asset
- **THEN** the workspace shows the "Unknown" placeholder and labels the cover as not assigned

### Requirement: Cover generation hand-off

The workspace SHALL let the author supply a photograph and obtain generated cover assets, and SHALL persist the resulting asset keys only when the vehicle is saved.

#### Scenario: Generate then save
- **WHEN** the author supplies a photograph and saves the vehicle
- **THEN** the generated cover keys are stored on that vehicle

#### Scenario: Generate then discard
- **WHEN** the author supplies a photograph but leaves without saving
- **THEN** no cover keys are attached to the vehicle

### Requirement: Editor access control

The vehicle editor SHALL require authentication. An unauthenticated visitor MUST NOT reach the listing or workspace.

#### Scenario: Unauthenticated visitor
- **WHEN** an unauthenticated visitor opens `/editor/vehicles`
- **THEN** the existing editor authentication flow applies and the listing does not render

### Requirement: Load and failure states

The listing and workspace SHALL show loading feedback while requests are in flight and SHALL present a retryable error when a request fails.

#### Scenario: Failed load
- **WHEN** the vehicle listing request fails
- **THEN** the page shows an error with a retry action instead of an empty list
