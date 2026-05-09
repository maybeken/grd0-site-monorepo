## ADDED Requirements

### Requirement: Modal dialog component
The system SHALL provide a reusable modal dialog component that can be invoked programmatically or via template, supporting confirmation, input, and custom content use cases.

#### Scenario: Confirmation modal
- **WHEN** a component calls the modal utility requesting user confirmation (e.g., "Are you sure you want to unpublish this post?")
- **THEN** a centered modal overlay appears with the confirmation message and "Confirm" / "Cancel" buttons

#### Scenario: User confirms action
- **WHEN** the user clicks "Confirm" in a confirmation modal
- **THEN** the modal closes and the calling code receives a resolved promise or callback with `true`

#### Scenario: User cancels action
- **WHEN** the user clicks "Cancel" or presses Escape in a confirmation modal
- **THEN** the modal closes and the calling code receives a resolved promise or callback with `false`

#### Scenario: Modal focus trap
- **WHEN** a modal is open
- **THEN** keyboard focus is trapped within the modal and cannot reach elements behind it

#### Scenario: Modal overlay click
- **WHEN** the user clicks the backdrop overlay outside the modal content
- **THEN** the modal does NOT close (clicking the overlay has no effect)

### Requirement: Toast notification component
The system SHALL provide a reusable toast notification system that displays brief, auto-dismissing messages for success, error, and informational states.

#### Scenario: Success toast
- **WHEN** an operation completes successfully (e.g., post saved)
- **THEN** a green success toast appears at the bottom-right corner with the success message and auto-dismisses after 4 seconds

#### Scenario: Error toast
- **WHEN** an operation fails (e.g., save error)
- **THEN** a red error toast appears with the error message and auto-dismisses after 6 seconds

#### Scenario: Manual dismiss
- **WHEN** the user clicks the close button on a toast
- **THEN** the toast is immediately dismissed

#### Scenario: Multiple toasts stack
- **WHEN** multiple toasts are triggered in quick succession
- **THEN** each toast appears below the previous one, forming a vertical stack

### Requirement: Replace browser dialogs in blog editor
All `confirm()`, `alert()`, and `prompt()` calls in the blog editor SHALL be replaced with the modal and toast utilities.

#### Scenario: Replace confirm dialog for post save
- **WHEN** the user clicks "Save" in the blog editor
- **THEN** a styled confirmation modal appears instead of a browser `confirm()` dialog

#### Scenario: Replace alert for save success
- **WHEN** a post is saved successfully
- **THEN** a success toast notification appears instead of a browser `alert()` dialog

#### Scenario: Replace alert for save failure
- **WHEN** a post save fails
- **THEN** an error toast notification appears instead of a browser `alert()` dialog

#### Scenario: Replace confirm dialog for post unpublish
- **WHEN** the user clicks "Delete" on a post in the PostTable
- **THEN** a styled confirmation modal appears asking "Are you sure you want to unpublish this post?"

#### Scenario: Replace prompt for new post URI
- **WHEN** the user clicks "New" to create a new post
- **THEN** a modal dialog with a URI input field appears instead of a browser `prompt()` dialog
