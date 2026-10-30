# vehicle-catalogue Specification

## Purpose

Defines the vehicle records behind the Grand Tour — what a vehicle entry contains, how it is read publicly, how authenticated authors mutate it, and how its cover assets are referenced.

## Requirements

### Requirement: Vehicle record

The system SHALL store one record per vehicle the author has driven. A record SHALL include make, model and year as required values, and MAY include body type, country code, country label, city, drive date, engine, power, transmission, drivetrain, fuel, notes, rating, sort order, cutout asset key and silhouette asset key.

#### Scenario: Required fields persist
- **WHEN** a vehicle is saved with make, model and year
- **THEN** the record is stored and a stable identifier is returned

#### Scenario: Missing required field is rejected
- **WHEN** a save omits the model
- **THEN** the request is rejected and no record is written

#### Scenario: Country label and code travel together
- **WHEN** a vehicle is saved with a selected country
- **THEN** the record carries both the country label and the matching upper-case ISO code

### Requirement: Public vehicle listing

The system SHALL expose the vehicle collection over a public read endpoint returning every stored vehicle with its technical fields and cover asset keys. The endpoint MUST NOT require authentication.

#### Scenario: Anonymous read succeeds
- **WHEN** an anonymous client requests the vehicle collection
- **THEN** the response contains the stored vehicles and no credentials are required

#### Scenario: Empty collection
- **WHEN** no vehicles are stored
- **THEN** the endpoint returns an empty collection with a success status

### Requirement: Public listing order

The system SHALL order the public listing by explicit sort order ascending when present, and by drive date ascending otherwise, so the collection reads as a chronology.

#### Scenario: Explicit order wins
- **WHEN** two vehicles have different sort orders
- **THEN** the lower sort order appears first regardless of drive date

#### Scenario: Fallback to chronology
- **WHEN** neither vehicle defines a sort order
- **THEN** the earlier drive date appears first

### Requirement: Authenticated vehicle mutation

The system SHALL restrict vehicle creation, update and deletion to authenticated authors. Requests without valid authentication MUST be rejected without modifying data.

#### Scenario: Authenticated update
- **WHEN** an authenticated author submits a changed vehicle
- **THEN** the stored record is updated and the result is returned

#### Scenario: Unauthenticated write is rejected
- **WHEN** a request creates or deletes a vehicle without valid authentication
- **THEN** the system rejects it and the stored data is unchanged

#### Scenario: Delete removes from listing
- **WHEN** an authenticated author deletes a vehicle
- **THEN** subsequent public listings no longer include it

### Requirement: Cover asset reference

A vehicle SHALL reference two optional cover assets: a photographic cutout and a flat silhouette. Both SHALL be optional, and a vehicle without a silhouette SHALL remain publishable.

#### Scenario: Both covers absent
- **WHEN** a vehicle is saved with no cover keys
- **THEN** the record is accepted and the public view renders the "Unknown" placeholder

#### Scenario: Only one cover present
- **WHEN** a vehicle has a cutout but no silhouette, or a silhouette but no cutout
- **THEN** the public view renders whichever cover exists

### Requirement: Cover upload authorization

The system SHALL issue short-lived presigned upload URLs for vehicle cover assets to authenticated authors only. Uploaded objects SHALL live under a dedicated vehicle prefix, separate from blog and article attachments.

#### Scenario: Author requests upload URL
- **WHEN** an authenticated author requests an upload URL for a cover file
- **THEN** the response contains a presigned URL scoped to the vehicle prefix and the object key to store on the record

#### Scenario: Storage unavailable
- **WHEN** storage is not configured, as in local development
- **THEN** the endpoint reports that uploads are unavailable rather than failing silently
