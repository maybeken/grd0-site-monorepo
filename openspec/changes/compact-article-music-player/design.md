# Design

## Context

See `proposal.md` for motivation. The relevant current state:

- `MusicPlayer` owns all playback behaviour: a hidden YouTube iframe, postMessage handling for readiness, time and state, seek translation, and the elapsed-time formatter.
- It is rendered in three places: `MusicView` (the full player, with previous/next slots), `MusicBlockEditor` (an authoring control), and `MusicBlock` (article rendering through `BlockRenderer`).
- `MusicPicker` can only select songs that already exist in the music catalogue; authoring a block from an arbitrary YouTube video is not possible.
- The article API requires a music block's video identifier and a non-empty title, with an optional artist.
- `MusicProgressBar` hard-codes `w-full md:w-1/2 lg:w-1/3` on its root.
- Article content spans the Tailwind `container` (up to 1280px) with no text-width cap, so the full player's `lg:w-1/3` cover is roughly 400px square.

## Goals / Non-Goals

**Goals:**

- Give article music blocks a horizontal, full-width presentation that is roughly a quarter of the current height.
- Let an author use any YouTube video in a music block, not only catalogue entries, without changing the block schema.
- Keep playback logic in one place and leave the music page visually unchanged.
- Keep the editor control and preview consistent with the reader's view.

**Non-Goals:**

- Redesigning the music page, which stays the only full-player surface.
- Changing the YouTube integration, playback state model, or block schema.
- Adding catalogue entries from a link; a link-authored block is part of the article only.
- Keyboard or screen-reader work on the play control.

## Decisions

**A variant prop on `MusicPlayer` rather than a second component.** `MusicPlayer` gains `variant?: 'full' | 'compact'`, defaulting to `full`. Both layouts share the hidden iframe, message handling, state and formatting. A separate `MusicPlayerCompact` would either duplicate that logic or need a new composable; the composable is the cleaner refactor of the two, but it is an abstraction this change does not need. If the layouts diverge significantly later, extracting the shared state into a composable is the natural follow-up.

**Progress bar width caps move to the caller.** `MusicProgressBar` becomes `relative w-full h-2`, and the full layout wraps it in `w-full md:w-1/2 lg:w-1/3`. A `fluid` prop was rejected: the constraint is layout, not behaviour, so it belongs with the layout that wants it. The music page keeps the same rendered width.

**The card shell stays with each consumer.** `MusicBlock` and `MusicBlockEditor` each own their `rounded-xl bg-shade` container and pass `variant="compact"`; the variant only decides the inner arrangement. This keeps page-specific chrome out of the shared player and lets article spacing change without touching other consumers. The editor additionally keeps its change-song action below the compact player.

**Compact time formatting.** The existing formatter always emits `HH:MM:SS`. The compact layout uses a second formatter that omits the hour field when the duration is under an hour. The full layout keeps the existing format so the music page is untouched.

**Dimensions.** Desktop: 88px square cover, 48px play control, inline `elapsed / total` on the same row as the title and artist, progress bar beneath the row, 12px internal gaps. Narrow screens: 64px cover, 40px control, time hidden, title and artist truncate (`min-w-0` plus `truncate`) and the row never stacks.

**YouTube link input lives in the existing picker.** `MusicPicker` gains a link field above the catalogue list rather than a separate route or a second picker. The picker is already the one place a block's song is chosen, so both sources share the same `select` event and the editor needs no extra state. A bare video identifier is accepted in addition to URLs, and parsing normalises the watch, `youtu.be`, `embed`, `shorts`, and `live` forms plus the `music.` and `youtube-nocookie.com` hosts to one identifier.

**Title and artist come from YouTube oEmbed at selection time.** The block schema and the API require a non-empty title, so the picker resolves the video's metadata through `https://www.youtube.com/oembed` before emitting the selection. The resolved identifier, title, and artist are stored exactly as a catalogue selection is, so link-authored blocks need no special rendering path. Falling back to an empty title and letting the player fill it was rejected: the article save would fail before the hidden iframe delivered its metadata. A backend proxy was rejected as unnecessary for a one-off client request.

**A failed link leaves the block untouched.** Parsing and resolution errors are reported inline in the picker, matching the picker's existing empty-state copy. The block is only mutated through the same `select` handler used by catalogue rows, so an invalid link cannot clear or partially update an existing selection.

## Risks / Trade-offs

- [Two layouts in one component accumulate conditionals] → keep the compact template small and self-contained; extract a composable if a third layout appears.
- [Hiding the time below the `sm` breakpoint loses an explicit readout] → the progress bar still shows position, and the full player remains one navigation away.
- [Reducing `MusicBlock` padding changes block rhythm slightly] → it is the intended effect of a shorter block; the outer `gap-6` between blocks is unchanged.
- [Existing behaviour: the hidden iframe loads as soon as a music block renders] → unchanged by this design; noted so it is not mistaken for a regression.
- [The browser-facing oEmbed request depends on YouTube's CORS response] → YouTube reflects the calling origin today; a failure surfaces as an inline picker error rather than a partial block, and the catalogue path is unaffected.
- [A well-formed link may point at a private, removed, or region-blocked video] → the picker reports the resolution failure and leaves the block unchanged; the identifier is never stored without resolved metadata.
