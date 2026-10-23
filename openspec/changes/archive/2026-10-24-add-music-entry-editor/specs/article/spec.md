# Spec Delta

## MODIFIED Requirements

### Requirement: Block composition

An article SHALL store an ordered list of blocks. Each block MUST have a stable identifier and a type of `markdown`, `image`, `music`, or `excalidraw`, with a payload appropriate to that type.

#### Scenario: Ordered blocks returned
- **WHEN** an article with multiple blocks is read
- **THEN** the blocks SHALL be returned in their stored order

#### Scenario: Unknown block type rejected
- **WHEN** an update contains a block type outside the supported set
- **THEN** the system SHALL reject the update with `400 Bad Request`
- **THEN** the stored article SHALL remain unchanged

#### Scenario: Block-level references
- **WHEN** an image block is read
- **THEN** it SHALL reference a gallery collection and filename
- **WHEN** a music block is read
- **THEN** it SHALL reference a song identifier and SHALL carry the title and artist captured at authoring time
- **WHEN** an excalidraw block is read
- **THEN** it SHALL reference the stored scene path

#### Scenario: Music block authored from a YouTube link
- **WHEN** a music block is read that was authored from a YouTube link
- **THEN** it SHALL reference the link's video identifier and SHALL carry the title and artist confirmed at authoring time, which MAY be author overrides of the resolved values
- **THEN** it SHALL be indistinguishable from a music block authored from the catalogue
