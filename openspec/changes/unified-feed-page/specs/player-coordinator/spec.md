## ADDED Requirements

### Requirement: Player coordinator store
The system SHALL provide a Pinia store `usePlayerCoordinator` that manages which music player is currently active across the application.

#### Scenario: Store initialization
- **WHEN** the store is first accessed
- **THEN** the store SHALL initialize with activePlayerId set to null

#### Scenario: Set active player
- **WHEN** a music player calls `play(id)` with its unique identifier
- **THEN** the store SHALL set activePlayerId to the provided id

#### Scenario: Clear active player
- **WHEN** a music player calls `pause()`
- **THEN** the store SHALL set activePlayerId to null

### Requirement: Player mutual exclusion
The system SHALL ensure that only one music player can be active at a time. When a player becomes active, all other players SHALL pause.

#### Scenario: Play new player while another is active
- **WHEN** a music player calls `play(myId)` while another player is active (activePlayerId !== myId)
- **THEN** the previously active player SHALL detect the change and pause itself, and the new player SHALL become active

#### Scenario: Play same player again
- **WHEN** a music player calls `play(myId)` while it is already active (activePlayerId === myId)
- **THEN** the player SHALL continue playing without interruption

#### Scenario: Pause active player
- **WHEN** the active music player calls `pause()`
- **THEN** the store SHALL set activePlayerId to null and the player SHALL pause

### Requirement: Player watches active state
Each music player component SHALL watch the activePlayerId in the store and respond to changes.

#### Scenario: Player detects it is no longer active
- **WHEN** activePlayerId changes to a different value than the player's id
- **THEN** the player SHALL pause playback

#### Scenario: Player detects it is now active
- **WHEN** activePlayerId changes to the player's id
- **THEN** the player SHALL start or resume playback

### Requirement: Unique player identification
Each music player instance SHALL have a unique identifier that distinguishes it from other players on the page.

#### Scenario: Music block player ID
- **WHEN** a music block is rendered on the feed listing or detail page
- **THEN** the music player SHALL use the feed block's id as its unique identifier

#### Scenario: Multiple music blocks have different IDs
- **WHEN** multiple music blocks are rendered on the same page
- **THEN** each player SHALL have a unique id, allowing the coordinator to distinguish between them
