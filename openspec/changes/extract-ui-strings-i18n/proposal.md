## Why

All UI strings are hardcoded in English across ~34 Vue components and TypeScript files in vue-web. This blocks future multi-language support and makes translation work difficult. Extracting strings now into a centralized, translator-friendly format enables efficient translation work later without touching component code.

## What Changes

- Install `vue-i18n` and configure i18n plugin in `main.ts`
- Create custom Vite plugin (`vite-plugin-gettext-yaml-i18n.ts`) to transform gettext-style YAML into vue-i18n format at build time
- Extract ~150-180 hardcoded strings from Vue templates, TypeScript files, and router meta
- Create YAML translation files following gettext conventions (`msgid`, `msgstr`, `msgctxt`, translator notes, source references)
- Replace all hardcoded strings with `$t()` or `useI18n()` calls
- Refactor `helpers/coffee.ts` constants (grind sizes, equipment types, taste dimensions) to use translation keys instead of raw strings
- Update router `meta.title` to use translation keys for document title composition
- Make date formatting locale-aware (currently hardcoded `en-GB`)

## Capabilities

### New Capabilities
- `ui-string-extraction`: Centralized YAML-based translation file system with gettext conventions for UI strings, including translator context and source references

### Modified Capabilities
<!-- No requirement changes to existing specs. Implementation details (how strings are stored) change, but behavioral specs remain valid. -->

## Impact

- **vue-web package**: All 34 Vue components and 12 TypeScript files with hardcoded strings require updates
- **New dependencies**: `vue-i18n` (Intlify), `yaml` (for YAML parsing in build plugin)
- **New files**: `src/i18n.ts` (plugin setup), `src/locales/en.yaml` (English strings), `vite-plugin-gettext-yaml-i18n.ts` (build-time YAML transformer)
- **Refactored**: `src/helpers/coffee.ts` (constants need translation keys), router config, date formatting in `TastingCard.vue`
- **No backend changes**: API calls and data structures remain unchanged
- **No breaking changes**: All strings stay English; only storage mechanism changes
