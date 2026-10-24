# editor-ui Specification

## Purpose

Defines the shared visual and verbal language of the authenticated editor console, so every section presents the same page structure, controls, states, terminology, and feedback to authors.

## Requirements

### Requirement: Shared editor page skeleton

Every editor list and edit page SHALL render through one page skeleton: a vertical stack with consistent spacing, a page header region, and a content region. Pages SHALL NOT hand-roll their own top-level spacing, heading block, or action row. Every editor page SHALL use the full available width.

#### Scenario: List page structure
- **WHEN** an author opens any section list page
- **THEN** the page SHALL render the shared page header followed by the content region
- **THEN** the outer spacing SHALL match every other list page

#### Scenario: Edit page structure
- **WHEN** an author opens any item edit page
- **THEN** the page SHALL render the shared page header followed by the edit form
- **THEN** the content SHALL use the full available width

#### Scenario: Article workspace layout
- **WHEN** an author opens an article edit route
- **THEN** the split workspace SHALL keep its editor and preview panes
- **THEN** its header and actions SHALL still follow the shared header composition

### Requirement: Page header composition

Each page header SHALL contain, in a fixed order: an optional back link, the page title as the only `h1`, an optional one-line description, and an optional action area aligned to the trailing edge. A page SHALL NOT render a second competing heading for the same content.

#### Scenario: List header
- **WHEN** a list page renders
- **THEN** it SHALL show the section name as the `h1`
- **THEN** it MAY show a one-line description
- **THEN** it SHALL NOT repeat the section name in a second heading

#### Scenario: Edit header
- **WHEN** an edit page renders
- **THEN** it SHALL show the item or section name as the `h1`
- **THEN** the back link SHALL appear in the header
- **THEN** the item's save action SHALL appear in the header action area

#### Scenario: Back navigation
- **WHEN** an edit page renders a back link
- **THEN** the link SHALL carry a visible label naming the destination list
- **THEN** the link SHALL appear in the same position on every edit page

### Requirement: Shared form controls

Editor forms SHALL use one set of form controls. Labels, text inputs, selects, and textareas SHALL share the same typography, border, background, padding, and focus treatment. Required fields SHALL mark the label with `*` and expose the requirement to assistive technology. Field-level validation messages SHALL render adjacent to the field in the shared error style.

#### Scenario: Text field
- **WHEN** a form renders a text input with its label
- **THEN** the label and input SHALL use the shared field styles
- **THEN** the input SHALL use the shared focus treatment

#### Scenario: Select and textarea
- **WHEN** a form renders a select or textarea
- **THEN** its visual treatment SHALL match the shared text input

#### Scenario: Required field
- **WHEN** a field is required
- **THEN** its label SHALL end with `*`
- **THEN** the control SHALL expose the required state to assistive technology

#### Scenario: Field validation
- **WHEN** a field fails validation
- **THEN** the message SHALL render next to that field
- **THEN** the message SHALL use the shared error style

### Requirement: Button conventions

The editor SHALL use four button variants: primary, secondary, destructive, and icon-only. Each page region SHALL contain at most one primary button. Every button SHALL use shared padding, radius, typography, and disabled styling. Icon-only buttons SHALL provide an accessible name and a tooltip.

#### Scenario: Primary action
- **WHEN** a page renders its main action
- **THEN** exactly one primary button SHALL represent it
- **THEN** secondary actions SHALL use the secondary variant

#### Scenario: Destructive action
- **WHEN** an action deletes or removes data
- **THEN** its button SHALL use the destructive variant

#### Scenario: Icon-only button
- **WHEN** a button displays only an icon
- **THEN** it SHALL have an accessible name
- **THEN** it SHALL expose a tooltip

#### Scenario: Disabled action
- **WHEN** an action is unavailable
- **THEN** its control SHALL use the shared disabled treatment

### Requirement: Loading, empty, and error states

Every page that depends on remote data SHALL render a loading state before content resolves. An empty collection SHALL render the shared empty state. A failed operation SHALL be reported exactly once, using the shared error presentation and a single wording family.

#### Scenario: Loading content
- **WHEN** a page is fetching its data
- **THEN** it SHALL show the shared loading state
- **THEN** it SHALL NOT show an empty state while loading

#### Scenario: Empty collection
- **WHEN** a collection contains no items
- **THEN** the page SHALL show the shared empty state
- **THEN** the message SHALL name the collection and use consistent punctuation

#### Scenario: Failed operation
- **WHEN** a load, save, or delete operation fails
- **THEN** the page SHALL report the failure once
- **THEN** the report SHALL use the shared error presentation and wording family
- **THEN** it SHALL NOT simultaneously show a toast and an inline error for the same failure

### Requirement: Create action terminology

Editor create actions SHALL be labelled `New <Item>`, using the same pattern in every section.

#### Scenario: Create action
- **WHEN** a page offers a create action
- **THEN** the label SHALL read `New <Item>`
- **THEN** the same pattern SHALL be used in every section

### Requirement: Persistence terminology

Editor persistence actions SHALL be labelled `Save`. `Saved` SHALL only describe state and SHALL NOT be used as a button label.

#### Scenario: Save action
- **WHEN** a page offers to persist changes
- **THEN** the button SHALL read `Save`
- **THEN** a saved state SHALL be shown as text rather than as a button label

### Requirement: Removal terminology

Destructive removal of an item SHALL be labelled `Delete`. Removing an item from a playlist without deleting its catalogue entry SHALL be labelled `Remove from playlist` and SHALL NOT be used for permanent deletion.

#### Scenario: Playlist removal
- **WHEN** an action removes a song from a playlist without deleting the catalogue entry
- **THEN** the action SHALL read `Remove from playlist`
- **WHEN** an action permanently deletes a catalogue entry
- **THEN** the action SHALL read `Delete` and its confirmation SHALL state that the entry leaves the catalogue

### Requirement: Editor text style

Editor UI text SHALL use sentence case. Labels and buttons SHALL omit trailing periods, full-sentence messages SHALL end with a period, and ellipses SHALL use the `…` character.

#### Scenario: Punctuation
- **WHEN** UI text is a label or button
- **THEN** it SHALL NOT end with a period
- **WHEN** UI text is a full sentence
- **THEN** it SHALL end with a period
- **THEN** ellipses SHALL use the `…` character

### Requirement: Text contrast

Editor text SHALL meet a contrast ratio of at least 4.5:1 against its background for normal text and 3:1 for large text. Muted text, placeholder text, and status labels SHALL meet the normal-text threshold. Status SHALL NOT be conveyed by colour alone.

#### Scenario: Body and muted text
- **WHEN** the editor renders body text, hints, or metadata
- **THEN** the text SHALL meet the normal-text contrast threshold against its background

#### Scenario: Placeholder text
- **WHEN** a control shows placeholder text
- **THEN** the placeholder SHALL meet the normal-text contrast threshold

#### Scenario: Status
- **WHEN** the editor shows a status such as published or draft
- **THEN** the status SHALL include a text label
- **THEN** the status SHALL NOT be conveyed by colour alone

### Requirement: Readable type scale

Editor body and interactive text SHALL render at no less than 14px, metadata and status labels at no less than 12px, and no editor text below 12px. Multi-line text SHALL use a line height of at least 1.5. Text inputs and textareas SHALL span the full width of their form container.

#### Scenario: Smallest text
- **WHEN** the editor renders body, control, or metadata text
- **THEN** interactive text SHALL be at least 14px
- **THEN** metadata text SHALL be at least 12px

#### Scenario: Multi-line text
- **WHEN** the editor renders a description, notes field, or hint spanning multiple lines
- **THEN** it SHALL use a line height of at least 1.5

#### Scenario: Field width
- **WHEN** the editor renders a text input or textarea
- **THEN** it SHALL span the full width of its form container

### Requirement: Visible focus and control boundaries

Every interactive editor control SHALL show a visible focus indicator with at least 3:1 contrast against adjacent colours. Text input, select, and textarea boundaries SHALL meet at least 3:1 contrast against the page background.

#### Scenario: Keyboard focus
- **WHEN** the author focuses a control with the keyboard
- **THEN** the control SHALL show a visible focus indicator
- **THEN** the indicator SHALL meet the non-text contrast threshold

#### Scenario: Control boundary
- **WHEN** a text input, select, or textarea renders
- **THEN** its boundary SHALL meet the non-text contrast threshold against the page background

### Requirement: Long-value readability

Values that exceed their container in list rows or page headers SHALL wrap or truncate. A truncated value SHALL expose its full text to assistive technology and through a tooltip.

#### Scenario: Long value in a row
- **WHEN** a title or description is longer than its container
- **THEN** the row SHALL wrap or truncate without overlapping adjacent content

#### Scenario: Truncated value
- **WHEN** a value is truncated
- **THEN** the full value SHALL be available to assistive technology
- **THEN** the full value SHALL be available through a tooltip

### Requirement: Feedback and discard confirmation

A mutating action SHALL report success with the shared toast when the page does not otherwise show an immediate visible result. A destructive action SHALL require confirmation whose text names the affected item and states whether it can be undone. Cancelling an editor form with unsaved edits SHALL require confirmation before the edits are discarded.

#### Scenario: Successful mutation
- **WHEN** a save, create, or delete succeeds
- **THEN** the author SHALL receive success feedback
- **THEN** the feedback SHALL use the shared toast unless the page visibly reflects the result

#### Scenario: Destructive confirmation
- **WHEN** the author triggers a destructive action
- **THEN** the system SHALL ask for confirmation
- **THEN** the confirmation SHALL name the item and state whether the action can be undone

#### Scenario: Discarding a dirty form
- **WHEN** the author cancels a form with unsaved edits
- **THEN** the system SHALL ask for confirmation
- **THEN** cancelling the confirmation SHALL keep the edits

### Requirement: Keyboard save shortcut

Every editable editor form SHALL support Ctrl+S on Windows and Linux and ⌘+S on macOS for its save action. The shortcut SHALL NOT bypass validation, the disabled state, or unsaved-change protection.

#### Scenario: Shortcut saves
- **WHEN** the author presses the save shortcut in an editable form
- **THEN** the form SHALL perform its save action just as the save button does

#### Scenario: Shortcut with nothing to save
- **WHEN** the form has no changes or invalid required fields
- **THEN** the shortcut SHALL NOT submit the form
