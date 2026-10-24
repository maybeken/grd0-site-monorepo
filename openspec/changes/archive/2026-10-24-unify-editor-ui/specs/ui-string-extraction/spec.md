# Spec Delta

## ADDED Requirements

### Requirement: Plural forms for counted labels

Editor labels that include a count SHALL use pluralization with distinct singular and plural forms. A count of one SHALL use the singular form; any other count SHALL use the plural form. Counted labels SHALL NOT render a plural noun for a count of one.

#### Scenario: Singular count
- **WHEN** a label reports a count of one asset, song, playlist, or pending change
- **THEN** the label SHALL use the singular noun

#### Scenario: Plural count
- **WHEN** a label reports a count other than one
- **THEN** the label SHALL use the plural noun

### Requirement: Canonical editor terminology

The translation catalogue SHALL define one canonical key per shared editor concept. Section-specific keys SHALL NOT duplicate identical text for the same concept; they SHALL reference the canonical key instead. Unused editor keys SHALL be removed when their behaviour no longer exists.

#### Scenario: Shared concept
- **WHEN** two editor sections display the same concept, such as a delete confirmation or a discard prompt
- **THEN** both SHALL use the same canonical key

#### Scenario: Duplicate keys
- **WHEN** the catalogue contains two keys with identical text for the same concept
- **THEN** the duplicate SHALL be removed or the usages SHALL be consolidated onto one key

### Requirement: Editor string coverage

All editor user-facing text SHALL come from the translation catalogue, including abbreviations, units and format hints, status labels, separators, and decorative glyphs attached to text. Only comments, brand names, and technical constants SHALL remain hardcoded.

#### Scenario: Abbreviation
- **WHEN** the editor displays an abbreviation such as the tasting "IDK" control
- **THEN** the text SHALL come from a translation key

#### Scenario: Format hint
- **WHEN** the editor displays an input format hint such as a time placeholder
- **THEN** the hint SHALL come from a translation key

#### Scenario: Separators and glyphs
- **WHEN** the editor renders separators or decorative glyphs alongside labels
- **THEN** those characters SHALL come from the translation catalogue when they accompany user-facing text

#### Scenario: Verification
- **WHEN** a developer searches editor templates and scripts for hardcoded UI strings
- **THEN** only comments, brand names, and technical constants SHALL remain
