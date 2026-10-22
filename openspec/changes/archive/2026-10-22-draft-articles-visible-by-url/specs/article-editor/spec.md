# Spec Delta

## MODIFIED Requirements

### Requirement: Metadata sidebar

The editing page SHALL show the article's URI, publication state, and publishing controls in a sidebar.

#### Scenario: URI
- **WHEN** the workspace is open
- **THEN** the sidebar SHALL show the article's current URI
- **THEN** the author SHALL be able to edit the URI

#### Scenario: Publication state
- **WHEN** the workspace is open
- **THEN** the sidebar SHALL show the publish timestamp and whether the article is a draft, scheduled, or published

#### Scenario: Public link
- **WHEN** the workspace is open
- **THEN** the sidebar SHALL offer a link to the article's public page regardless of its publication state

#### Scenario: Copy link
- **WHEN** the article has a publish date
- **THEN** the sidebar SHALL offer a copy control beside the URI
- **THEN** activating it SHALL copy the article's public URL

#### Scenario: Publication actions
- **WHEN** the article is published
- **THEN** Publish now SHALL be hidden
- **WHEN** the article is a draft
- **THEN** Unpublish SHALL be hidden
- **WHEN** the article is scheduled
- **THEN** both Publish now and Unpublish SHALL be available

#### Scenario: Delete
- **WHEN** the author chooses to delete the article from the sidebar
- **THEN** the editor SHALL require confirmation before permanently deleting it
