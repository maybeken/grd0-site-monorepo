# Taste
- Requires British English spelling and phrasing (e.g. behaviour, colour, initialise, licence) in code, comments, documentation, and commit messages. Confidence: 0.9
- Wants non-trivial changes (e.g. dependency upgrades) done on a dedicated new git branch rather than directly on the current/main branch. Confidence: 0.6
- When creating a new branch, expects the branch to start from up-to-date remote main: fetch origin and rebase onto origin/main before beginning work. Confidence: 0.8
- Favors moving dependencies to the latest versions, including major bumps, accepting that code fixes/breakage follow. Confidence: 0.5
- When asking to "stage the changes made," expects staging scoped to the current session's work — unrelated pre-existing working-tree edits and untracked build artifacts/binary outputs left out. Confidence: 0.5
- Prefers staging changes for review without immediately committing. Confidence: 0.4
- Works in explicit increments — stage, then commit, then push — and expects work committed locally to be pushed to origin on the feature branch with its upstream set, not left local. Confidence: 0.5
- Prefers Vue over React for frontend work; questions React-based solutions and asks whether they can be done as a Vue project instead. Confidence: 0.5
- Expects a single, non-duplicated taste/memory store and will question redundant or oddly nested files in the agent's config directory. Confidence: 0.3
