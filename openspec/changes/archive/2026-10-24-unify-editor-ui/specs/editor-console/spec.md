# Spec Delta

## ADDED Requirements

### Requirement: Consistent list page header and create action

Each section list page SHALL use the shared editor page header. The header SHALL show the section name as its title, MAY show a one-line description of the section, and SHALL place the section's create action as its primary action. The create action SHALL occupy the same position on every list page, and activating it SHALL open the create form directly below the header. The Articles listing SHALL keep its existing article list and draft-creation flow.

#### Scenario: List page header
- **WHEN** the author opens any section list page
- **THEN** the header SHALL show that section's name
- **THEN** the header SHALL NOT render a second heading for the same section

#### Scenario: Create action position
- **WHEN** a list page offers creation
- **THEN** the create action SHALL appear in the header action area
- **THEN** it SHALL appear in the same position in every section

#### Scenario: Opening the create form
- **WHEN** the author activates the create action
- **THEN** the create form SHALL open directly below the page header
- **THEN** the author SHALL remain on the list route

### Requirement: Consistent edit page navigation

Each item edit page SHALL use the shared editor page header with a visible back link to its section's list page. The back link SHALL use the same label pattern and position in every section. When the route identifies an item that does not exist, the page SHALL show the shared error state with a link back to the section's list page and SHALL NOT create a record.

#### Scenario: Returning to the list
- **WHEN** the author activates the back link on an edit page
- **THEN** the system SHALL navigate to that section's list page
- **THEN** the link SHALL name the destination

#### Scenario: Missing item
- **WHEN** the author opens an edit route for an unknown item
- **THEN** the page SHALL show the shared error state
- **THEN** it SHALL offer a link back to the section's list page
- **THEN** it SHALL NOT create a record
