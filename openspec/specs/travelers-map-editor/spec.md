# travelers-map-editor Specification

## Purpose

Lets an authenticated author create, place, publish, and maintain travelers' map locations visually, without hand-editing JSON or trusting unvalidated coordinates.

## Requirements

### Requirement: Map locations listing

The travelers' map editor SHALL provide a listing at `/editor/travel` that shows every location with its title, subtitle when present, slug, and publish state, and SHALL allow deleting a location with confirmation.

#### Scenario: List locations
- **WHEN** the author opens the map locations listing
- **THEN** every stored location SHALL be listed with its title, optional subtitle, and slug

#### Scenario: Publish state shown
- **WHEN** the author views the listing
- **THEN** each location SHALL show whether it is published or a draft

#### Scenario: Open a location
- **WHEN** the author selects a location
- **THEN** the editor SHALL navigate to that location's workspace route

#### Scenario: Empty listing
- **WHEN** no locations exist
- **THEN** the listing SHALL show an empty state

#### Scenario: Delete a location
- **WHEN** the author deletes a location and confirms
- **THEN** the location SHALL be removed from the listing and from the public map

#### Scenario: Cancelled deletion
- **WHEN** the author dismisses the deletion confirmation
- **THEN** the location SHALL remain unchanged

#### Scenario: Failed load
- **WHEN** locations cannot be loaded
- **THEN** the editor SHALL report the failure without blocking the rest of the page

### Requirement: Location creation

The listing SHALL create locations through a title-only inline form. Saving SHALL generate a kebab-case slug from the title, create the location as a draft, and open its workspace. Creation MUST NOT overwrite an existing location.

#### Scenario: Create a location
- **WHEN** the author enters a title and confirms
- **THEN** the location SHALL be created with a slug derived from the title in kebab case
- **THEN** the editor SHALL navigate to the new location's workspace

#### Scenario: Created as a draft
- **WHEN** the author creates a location
- **THEN** the location SHALL start unpublished
- **THEN** it SHALL NOT appear on the public map until the author publishes it

#### Scenario: Title required
- **WHEN** the title is blank
- **THEN** the create action SHALL be unavailable

#### Scenario: Slug collision
- **WHEN** the generated slug already belongs to another location
- **THEN** the editor SHALL reject the creation and report the collision
- **THEN** the existing location SHALL remain unchanged

#### Scenario: Creation fails
- **WHEN** the location cannot be created
- **THEN** the editor SHALL report the failure and keep the entered title

### Requirement: Location workspace

Each location SHALL have a workspace at `/editor/travel/:slug` that edits its title, subtitle, icon, colour, text colour, publish state, visibility thresholds, and coordinates. Opening the route directly SHALL load the location without visiting the listing first.

#### Scenario: Direct link
- **WHEN** the author opens a location workspace route directly
- **THEN** the location SHALL load with its current values

#### Scenario: Unknown location
- **WHEN** the route names a slug that does not exist
- **THEN** the editor SHALL show a not-found message and a link back to the listing

#### Scenario: Edit details
- **WHEN** the author changes the title, subtitle, icon, colour, or text colour and saves
- **THEN** the new values SHALL be persisted
- **THEN** the public map SHALL render the updated label

#### Scenario: Text colour default
- **WHEN** no text colour is chosen
- **THEN** the public map SHALL use its default text colour for the location

#### Scenario: Title required
- **WHEN** the title is blank
- **THEN** the save action SHALL be unavailable

### Requirement: Slug lifecycle

A location's slug SHALL be generated at creation and SHALL NOT be editable afterwards. The workspace SHALL display the slug and its public map link but MUST NOT offer a rename control.

#### Scenario: Slug shown
- **WHEN** the author opens a location workspace
- **THEN** the slug SHALL be displayed as read-only information

#### Scenario: No rename path
- **WHEN** the author edits a location
- **THEN** the editor SHALL NOT provide a way to change its slug

### Requirement: Publish state

Each location SHALL be either published or a draft. The workspace SHALL show the current state and SHALL allow toggling it, with the change taking effect when the location is saved. The public map SHALL show published locations only. Existing locations SHALL remain published after the feature is deployed.

#### Scenario: Draft is hidden
- **WHEN** a location is a draft
- **THEN** the public map SHALL NOT show it

#### Scenario: Publish a location
- **WHEN** the author publishes a location and saves
- **THEN** the public map SHALL show it

#### Scenario: Unpublish a location
- **WHEN** the author unpublishes a location and saves
- **THEN** the public map SHALL stop showing it

#### Scenario: Pending toggle
- **WHEN** the author toggles the publish state without saving
- **THEN** the stored state SHALL remain unchanged

#### Scenario: Status visible
- **WHEN** the author views a location in the listing or workspace
- **THEN** its published or draft state SHALL be visible

#### Scenario: Existing locations
- **WHEN** the feature is deployed over existing data
- **THEN** every existing location SHALL remain published

### Requirement: Icon selection

The workspace SHALL provide a searchable icon selector. The selector SHALL search the full icon set, make every icon reachable, and render icons and searches from data served by the application. The editor MUST NOT request icon data or search results from an external icon service.

#### Scenario: Search icons
- **WHEN** the author enters a search query
- **THEN** the selector SHALL list matching icons from the full icon set
- **THEN** the match count SHALL reflect the full set rather than only the rendered page

#### Scenario: Browse the full set
- **WHEN** the author opens the selector without a query
- **THEN** it SHALL offer an initial set of icons
- **THEN** the author SHALL be able to reach every remaining icon by scrolling or requesting more

#### Scenario: Select an icon
- **WHEN** the author selects an icon
- **THEN** that icon SHALL become the location's icon
- **THEN** the label preview and the preview beside the selector SHALL update

#### Scenario: Clear an icon
- **WHEN** the author clears the icon
- **THEN** the location SHALL have no icon

#### Scenario: Local data only
- **WHEN** the author opens or searches the selector, or the app renders an icon
- **THEN** the app SHALL NOT make a request to an external icon service

#### Scenario: Icon data fails to load
- **WHEN** the local icon data cannot be loaded
- **THEN** the selector SHALL report the failure without blocking the rest of the workspace

### Requirement: Coordinate editing

The workspace SHALL allow a location's latitude and longitude to be edited numerically or by selecting a point on a map preview. Coordinates outside the valid latitude or longitude range SHALL be rejected.

#### Scenario: Numeric entry
- **WHEN** the author enters a latitude and longitude and saves
- **THEN** the coordinates SHALL be persisted as the location's position

#### Scenario: Pick on the map
- **WHEN** the author activates a point on the map preview
- **THEN** the latitude and longitude fields SHALL update to that point
- **THEN** the marker SHALL move to that point

#### Scenario: Preview follows the form
- **WHEN** the author changes the coordinates numerically
- **THEN** the map preview marker SHALL move without requiring a save

#### Scenario: Invalid coordinate
- **WHEN** a latitude is outside -90 to 90 or a longitude is outside -180 to 180
- **THEN** the save action SHALL be unavailable and the invalid field SHALL be identified

#### Scenario: Existing invalid data
- **WHEN** a stored location has an out-of-range coordinate
- **THEN** the workspace SHALL load its values so they can be corrected

### Requirement: Visibility thresholds

The workspace SHALL edit `display_at` and `hide_at` as zoom levels while storing the existing resolution values. An empty threshold SHALL mean the location is not constrained on that side.

#### Scenario: Show from a zoom level
- **WHEN** the author sets a visible-from zoom and saves
- **THEN** a display threshold SHALL be stored so the location appears only when zoomed in to at least that level

#### Scenario: Hide beyond a zoom level
- **WHEN** the author sets a visible-until zoom and saves
- **THEN** a hide threshold SHALL be stored so the location disappears when zoomed in past that level

#### Scenario: Unconstrained visibility
- **WHEN** neither threshold is set
- **THEN** the location SHALL be eligible at every zoom level, as before

#### Scenario: Existing thresholds preserved
- **WHEN** a location with stored thresholds is opened
- **THEN** the workspace SHALL show the equivalent zoom levels
- **THEN** saving without changes SHALL NOT alter the stored values

### Requirement: Palette parity

The colours and text colours offered by the editor SHALL exactly match the palette the public map renders. Every stored colour value SHALL render with its intended background, and the editor MUST NOT offer a value the map does not render.

#### Scenario: Editor options
- **WHEN** the author opens the colour selector
- **THEN** it SHALL offer every palette entry the public map supports

#### Scenario: Stored colour renders
- **WHEN** a location stores a palette colour
- **THEN** the public map SHALL render that colour rather than the default

#### Scenario: Dark blue palette entry
- **WHEN** the palette includes `darkblue`
- **THEN** the editor SHALL offer it
- **THEN** the public map SHALL render it instead of the default background

### Requirement: Save and discard behaviour

The workspace SHALL persist changes only when the author saves, and SHALL warn before discarding unsaved changes through in-application navigation, reload, or closing the tab.

#### Scenario: Pending changes
- **WHEN** the author edits fields without saving
- **THEN** the stored location SHALL remain unchanged

#### Scenario: Save pending changes
- **WHEN** the author saves with valid pending changes
- **THEN** every changed field SHALL be persisted
- **THEN** the workspace SHALL report success and remain on the route

#### Scenario: Navigate away with unsaved changes
- **WHEN** the author navigates away with unsaved changes
- **THEN** the editor SHALL ask for confirmation
- **THEN** cancelling SHALL keep the workspace and the changes

#### Scenario: Save fails
- **WHEN** saving fails
- **THEN** the workspace SHALL report the failure and keep the pending changes

#### Scenario: Keyboard save
- **WHEN** the author presses Ctrl+S or Cmd+S with valid pending changes
- **THEN** the pending changes SHALL be saved as if the Save action were used
- **THEN** the browser's save dialog SHALL NOT open

#### Scenario: Keyboard save with an invalid form
- **WHEN** the author presses Ctrl+S or Cmd+S while the title is blank or a coordinate is invalid
- **THEN** nothing SHALL be saved

### Requirement: Map preview isolation

The coordinate picker SHALL be a dedicated editor component. It SHALL reuse the shared coordinate conversion and colour helpers that the public map uses, and MUST NOT fetch public map data or emit navigation.

#### Scenario: Independent preview
- **WHEN** the picker renders in the workspace
- **THEN** it SHALL display only the location being edited
- **THEN** it SHALL NOT request the public location list

#### Scenario: Shared conversion
- **WHEN** the picker converts a clicked point to coordinates
- **THEN** it SHALL use the same conversion as the public map
- **THEN** the same stored position SHALL render at the same place in both maps

#### Scenario: No navigation
- **WHEN** the author interacts with the picker
- **THEN** the page SHALL NOT navigate or open a context menu
