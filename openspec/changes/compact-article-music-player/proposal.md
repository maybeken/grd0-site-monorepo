# Proposal

## Why

A music block in an article stacks a full-width square cover, title, artist, progress bar, time and a large round control vertically, so a single block takes roughly half a screen. The article content spans the full container, which makes the cover alone around 400px tall. The full presentation belongs on the music page, where the player is the content; in an article it is one line of flavour and should not interrupt reading.

## What Changes

- Article music blocks render a compact, full-width horizontal player: a small square cover, title and artist on one line each, elapsed and total time, a play/pause control, and a progress bar spanning the block.
- The compact presentation is additive: the music page keeps the existing full stacked player, while the editor's live preview and inline block control both use the compact player.
- The music picker accepts a YouTube link as well as a catalogue entry. The link is parsed for its video identifier and the video's title and artist are resolved at authoring time, so the stored block shape is unchanged and link-authored blocks render exactly like catalogue entries.
- Link input accepts a bare video identifier or a link in the watch, `youtu.be`, `embed`, `shorts`, or `live` form, including `music.youtube.com` and `youtube-nocookie.com` hosts.
- An invalid or unreadable link reports an inline error in the picker and leaves the block unchanged.
- The compact variant drops the leading hour from the time display when the song is under an hour, so it reads `3:45` rather than `00:03:45`.
- The progress bar's width caps move out of the shared bar into the full layout, letting the compact bar span the whole block without changing the music page.
- The editor's live preview renders through the article block renderer and therefore shows the compact player, matching what a reader sees.
- Long titles and artist names truncate rather than wrap, and the block stays horizontal on narrow screens.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `article-view`: The **Music block rendering** requirement changes so article music blocks use the compact full-width presentation, including a span-width progress bar and inline time display.
- `article-editor`: The **Music block authoring** requirement changes so a block can also be authored from a YouTube link, and so the block's inline control renders the compact player, matching the article and the live preview.

## Impact

- `vue-web/src/components/music/MusicPlayer.vue`: an opt-in compact variant alongside the existing full layout.
- `vue-web/src/components/music/MusicProgressBar.vue`: width constraints move to the callers.
- `vue-web/src/components/article/blocks/MusicBlock.vue`: requests the compact variant and hosts the card shell.
- `vue-web/src/components/editor/article/blocks/MusicBlockEditor.vue`: requests the compact variant for its control, with the change-song action kept below.
- `vue-web/src/components/editor/article/MusicPicker.vue`: a YouTube link form that parses the video identifier, resolves title and artist through YouTube oEmbed, and leaves the block unchanged when the link is invalid or unreadable.
- No API, worker, storage, or database changes.
