# Tasks

## 1. Shared player support

- [x] 1.1 Add a `variant?: 'full' | 'compact'` prop to `MusicPlayer`, defaulting to `full`, leaving the existing full template behaviour unchanged — verify `pnpm type-check` passes and `MusicView` still renders the stacked player.
- [x] 1.2 Remove the responsive width caps from `MusicProgressBar`'s root and apply `w-full md:w-1/2 lg:w-1/3` in the full layout's wrapper — verify the music page progress bar keeps its current width and `pnpm type-check` passes.
- [x] 1.3 Add a compact elapsed-time formatter that omits the hour field below one hour while the full layout keeps `HH:MM:SS` — verify a 225-second duration renders `3:45` in compact and `00:03:45` in full.

## 2. Compact article layout

- [x] 2.1 Build the compact row inside `MusicPlayer`: square cover, truncated title and artist, inline elapsed/total time, and play control, with 88px cover and 48px control on desktop and 64px/40px with the time hidden on narrow screens — verify `pnpm type-check` passes and the row stays horizontal at a 375px viewport.
- [x] 2.2 Place the progress bar beneath the row so it spans the full width of the block — verify by inspection that the bar reaches both edges of the card.
- [x] 2.3 Render `MusicBlock` with `variant="compact"` and tighten its card padding — verify the article page shows the compact player and the editor live preview matches it.
- [x] 2.4 Render `MusicBlockEditor` with `variant="compact"` while keeping the change-song action below — verify the editor control and preview show the same compact player.

## 3. Verification

- [x] 3.1 Run `pnpm type-check` and `pnpm build` in `vue-web` and confirm both succeed.
- [x] 3.2 Check play, pause, and seek on an article music block and on the editor control, and confirm the music page player is visually unchanged.

## 4. YouTube link authoring

- [x] 4.1 Add a YouTube link form above the catalogue list in `MusicPicker`, reusing the existing `select` event so the editor needs no new state — verify the picker still closes and updates the block after a link selection.
- [x] 4.2 Parse a bare video identifier and the watch, `youtu.be`, `embed`, `shorts`, and `live` forms, including `music.youtube.com` and `youtube-nocookie.com` hosts — verify the identifier is extracted for each form and non-YouTube hosts are rejected.
- [x] 4.3 Resolve the video's title and artist through YouTube oEmbed before emitting the selection, falling back to the identifier when the title is empty — verify a link selection stores the same identifier, title, and artist shape as a catalogue selection.
- [x] 4.4 Report an invalid or unreadable link inline in the picker and leave the current block selection unchanged — verify the error clears on a subsequent valid submission.
- [x] 4.5 Add the picker's placeholder, action, loading, and error strings to `vue-web/src/locales/en.yaml` — verify no hard-coded copy remains in the component.
- [x] 4.6 Run `pnpm type-check` and `pnpm build` in `vue-web` and confirm both succeed.

## Workflow follow-up

- Archive the change after verification is satisfactory.
