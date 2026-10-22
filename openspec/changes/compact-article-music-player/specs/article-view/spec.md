# Spec Delta

## MODIFIED Requirements

### Requirement: Music block rendering

A music block SHALL render a playable player through the shared player component using the compact, full-width horizontal presentation. The compact presentation MUST place the cover, title, artist, elapsed and total time, and the play control in one row, with a progress bar spanning the block. The full stacked presentation remains in use where the player is the page content.

#### Scenario: Play a song
- **WHEN** a reader presses play on a music block
- **THEN** the referenced song SHALL play through the shared player
- **THEN** the displayed title and artist SHALL come from the block's stored values

#### Scenario: Several music blocks
- **WHEN** an article contains more than one music block
- **THEN** each block SHALL control only its own playback

#### Scenario: Player readiness
- **WHEN** a music block's embedded player signals readiness
- **THEN** the play control SHALL become enabled
- **THEN** play, pause, and seek commands SHALL be delivered to that player

#### Scenario: Compact horizontal presentation
- **WHEN** an article renders a music block
- **THEN** the cover, title, artist, time, and play control SHALL be arranged in a single horizontal row
- **THEN** the block SHALL occupy substantially less vertical space than the full stacked player

#### Scenario: Progress bar spans the block
- **WHEN** an article renders a music block
- **THEN** the progress bar SHALL span the full width of the block

#### Scenario: Compact time display
- **WHEN** an article renders a music block
- **THEN** the elapsed and total time SHALL be shown together
- **THEN** a duration under one hour SHALL be shown without a leading hour field

#### Scenario: Long title on a narrow viewport
- **WHEN** the viewport is narrow or the title or artist is long
- **THEN** the title and artist SHALL truncate rather than wrap
- **THEN** the block SHALL keep its horizontal presentation

#### Scenario: Full presentation on the music page
- **WHEN** the shared player renders on the music page
- **THEN** it SHALL keep the full stacked presentation
