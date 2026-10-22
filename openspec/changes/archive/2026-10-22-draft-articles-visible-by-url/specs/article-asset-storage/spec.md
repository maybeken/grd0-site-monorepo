# Spec Delta

## MODIFIED Requirements

### Requirement: Scene access endpoints

Scenes SHALL be addressed through the article URI and the scene identifier. Writing and deleting MUST require authentication; reading SHALL be public for any existing article regardless of publish state.

#### Scenario: Write a scene
- **WHEN** an authenticated author writes a scene for their article
- **THEN** the system SHALL store the JSON at the article's scene path
- **THEN** the system SHALL respond with success

#### Scenario: Read a scene
- **WHEN** a client requests a scene referenced by a published article
- **THEN** the system SHALL return the stored JSON

#### Scenario: Read a scene of an unlisted article
- **WHEN** a client requests a scene referenced by a draft or scheduled article
- **THEN** the system SHALL return the stored JSON

#### Scenario: Read a scene of an unknown article
- **WHEN** a client requests a scene using a URI that no article uses
- **THEN** the system SHALL respond with `404 Not Found`

#### Scenario: Unauthenticated write
- **WHEN** an unauthenticated client attempts to write or delete a scene
- **THEN** the system SHALL reject the request and SHALL NOT modify storage
