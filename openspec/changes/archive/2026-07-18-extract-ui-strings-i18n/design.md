## Context

The vue-web frontend has ~150-180 hardcoded English strings across 34 Vue components and 12 TypeScript files. No i18n infrastructure exists. The heaviest concentration is in the coffee subsystem (TastingForm, TastingCard, BeanForm, EquipmentForm, helpers/coffee.ts) with ~70+ strings.

Current state:
- All strings hardcoded in templates and TypeScript
- `helpers/coffee.ts` exports string arrays used as both display labels and data values
- Router `meta.title` hardcoded for document titles
- Date formatting uses hardcoded `en-GB` locale in TastingCard.vue
- `alert()` and `confirm()` used for user notifications in editor views

Constraints:
- Backend API calls are out of scope
- Only English for now; translation work is future scope
- Language preference will be server-side config in the future
- Must not break existing functionality or data contracts

## Goals / Non-Goals

**Goals:**
- Extract all UI strings into centralized YAML files with gettext conventions
- Set up vue-i18n infrastructure for future multi-language support
- Provide translator context (msgctxt, notes, source references) in YAML files
- Separate display labels from data values in coffee constants
- Make date formatting locale-aware

**Non-Goals:**
- Actual translation to other languages (future work)
- Server-side language preference configuration (future work)
- Backend API changes or data structure modifications
- Replacing `alert()`/`confirm()` with custom UI dialogs
- Changing routing structure (no path-based locale like `/en/gallery`)

## Decisions

**1. Use vue-i18n (Intlify) for i18n framework**

Rationale: De facto standard for Vue 3, native composition API support (`useI18n()`), handles pluralization, interpolation, and date/number formatting. Alternatives considered:
- `@tolgee/vue`: Adds SaaS dependency, overkill for current needs
- Custom solution: Reinventing wheel, missing edge cases

**2. Use YAML format with gettext conventions**

Rationale: YAML supports comments for translator context, human-readable. Gettext conventions (`msgid`, `msgstr`, `msgctxt`, `#.`, `#:`) provide standardized translator workflow.

Alternatives considered:
- JSON: No comments, write-only for humans
- Fluent (.ftl): Purpose-built but less familiar, needs custom loader
- TypeScript with JSDoc: Type-safe but loses i18n tooling

**3. Custom Vite plugin for YAML transformation**

Rationale: The gettext-style YAML format (with nested `msgid`/`msgstr`/`msgctxt` objects) is not directly compatible with vue-i18n's expected flat key-value format. A custom Vite plugin (`vite-plugin-gettext-yaml-i18n.ts`) transforms the YAML at build time, extracting `msgstr` values into a flat object keyed by the full path (e.g., `nav.tagline`). This runs at build time with zero runtime overhead.

The plugin:
- Parses YAML using the `yaml` npm package
- Traverses the nested structure recursively
- Extracts `msgid` as the key and `msgstr` as the value
- Returns a flat `Record<string, string>` object
- Only applies to `.yaml` files in the `locales/` directory

This allows translators to organize YAML entries in any nested structure while the code references use the stable `msgid` keys.

Alternatives considered:
- `@intlify/vite-plugin-vue-i18n`: Deprecated and incompatible with Vite 8
- `@intlify/unplugin-vue-i18n`: Cannot handle custom gettext format, expects flat YAML
- Runtime transformation: Adds unnecessary runtime overhead

**4. Nested YAML structure by feature/component**

Rationale: Mirrors codebase organization, makes strings easy to locate, scales well as more strings are added. Structure: `coffee.tastingForm.titleNew`, `nav.gallery`, etc.

**5. Refactor coffee constants to use translation keys**

Rationale: `helpers/coffee.ts` currently exports `['Extra Fine', 'Fine', ...]` which are used as both display labels and data values. Translating these would break backend contracts. Refactor to:
```typescript
{ value: 'extra-fine', labelKey: 'coffee.grindSizes.extraFine' }
```
Components use `$t(grindSize.labelKey)` for display, `grindSize.value` for data.

**6. Router meta.title with translation keys**

Rationale: Current `meta.title` is hardcoded English. Change `meta.title` to contain translation keys instead. Document title composition translates the key. This avoids adding a new meta field and keeps the API simple.

**7. Use dayjs with dynamic locale imports for date formatting**

Rationale: TastingCard.vue currently uses `toLocaleDateString('en-GB', ...)`. dayjs is already a dependency, supports dynamic locale imports, integrates with vue-i18n's datetime formatting.

## Risks / Trade-offs

**[Risk] YAML parsing overhead** → Mitigation: Custom Vite plugin (`vite-plugin-gettext-yaml-i18n.ts`) transforms YAML to optimized JS at build time, no runtime overhead.

**[Risk] Missing strings during extraction** → Mitigation: Systematic search across all .vue and .ts files, verify with grep for hardcoded strings post-extraction.

**[Risk] Breaking coffee data contracts** → Mitigation: Refactor constants to separate `value` (data) from `labelKey` (display). Backend still receives original enum values.

**[Risk] Translator confusion from context** → Mitigation: Use gettext conventions consistently, include source file references (`#:`) and translator notes (`#.`).

**[Risk] Large YAML file becomes unwieldy** → Mitigation: Nested structure by feature keeps it organized. Can split into multiple files later if needed (vue-i18n supports lazy loading).

**[Trade-off] YAML over JSON**: Less tooling support, but human-readable with comments. Acceptable since translators will work with this directly.

**[Trade-off] No path-based routing**: Simpler implementation, but loses SEO benefits and shareable locale URLs. Acceptable for now since only English is supported.
