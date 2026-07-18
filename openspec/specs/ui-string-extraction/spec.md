# ui-string-extraction

## Purpose

Extract hardcoded UI strings from Vue templates and TypeScript files into gettext-style YAML translation files, enabling internationalization via vue-i18n.

## Requirements

### Requirement: YAML translation file format
The system SHALL store all UI strings in YAML files following gettext conventions. Each string entry MUST include `msgid` (translation key), `msgstr` (English text), and `msgctxt` (context description). Translator notes SHALL use `#.` comments and source file references SHALL use `#:` comments.

#### Scenario: Translation file structure
- **WHEN** a developer or translator opens `src/locales/en.yaml`
- **THEN** the file SHALL contain nested YAML objects by feature (e.g., `coffee.tastingForm.titleNew`)
- **THEN** each string entry SHALL have `msgid`, `msgstr`, and `msgctxt` fields
- **THEN** translator notes SHALL appear as `#.` comments above the entry
- **THEN** source file references SHALL appear as `#:` comments with file path and line number

#### Scenario: Gettext field usage
- **WHEN** code references a translation key via `$t('coffee.tastingForm.titleNew')`
- **THEN** the YAML file SHALL contain an entry where `msgid` equals `'coffee.tastingForm.titleNew'`
- **THEN** the `msgstr` field SHALL contain the English text `'New Tasting'`
- **THEN** the `msgctxt` field SHALL describe where and when this string appears

### Requirement: YAML transformation plugin
The system SHALL include a custom Vite plugin (`vite-plugin-gettext-yaml-i18n.ts`) that transforms gettext-style YAML files into vue-i18n compatible format at build time. The plugin MUST use `msgid` as the translation key and `msgstr` as the value, producing a flat key-value mapping. The plugin SHALL only apply to `.yaml` files in the `locales/` directory.

#### Scenario: Build-time transformation
- **WHEN** Vite builds the application
- **THEN** the plugin SHALL parse all `.yaml` files in `src/locales/`
- **THEN** the plugin SHALL traverse the nested YAML structure
- **THEN** the plugin SHALL extract `msgid` values as keys and `msgstr` values as values
- **THEN** the transformed output SHALL be importable as a standard ES module

#### Scenario: Key extraction from msgid
- **WHEN** a YAML entry contains `msgid: "coffee.view.title"` and `msgstr: "My Coffee Journey"`
- **THEN** the output SHALL contain the key-value pair `"coffee.view.title": "My Coffee Journey"`
- **THEN** the YAML nesting structure SHALL NOT affect the key name

#### Scenario: Plugin scope
- **WHEN** a `.yaml` file is imported
- **THEN** the plugin SHALL only transform files with paths containing `/locales/`
- **THEN** other `.yaml` files SHALL not be affected by the plugin

### Requirement: vue-i18n plugin integration
The system SHALL initialize vue-i18n as a Vue plugin in `main.ts` and load the English YAML translation file. The i18n instance SHALL be configured to use the composition API.

#### Scenario: Plugin initialization
- **WHEN** the Vue application starts
- **THEN** vue-i18n SHALL be installed as a plugin via `app.use(i18n)`
- **THEN** the English translation file SHALL be loaded from `src/locales/en.yaml`
- **THEN** the default locale SHALL be set to `'en'`

#### Scenario: Composition API usage
- **WHEN** a Vue component needs to access translations
- **THEN** the component SHALL use `const { t } = useI18n()` from `vue-i18n`
- **THEN** the component SHALL call `t('key.path')` to retrieve translated strings

### Requirement: Template string replacement
All hardcoded strings in Vue templates SHALL be replaced with `$t()` calls or composition API `t()` calls. This includes text content, placeholder attributes, title attributes, and aria-labels.

#### Scenario: Template text replacement
- **WHEN** a Vue template contains hardcoded text like `<h1>My Coffee Journey</h1>`
- **THEN** it SHALL be replaced with `<h1>{{ $t('coffee.view.title') }}</h1>`
- **THEN** the translation key SHALL exist in the YAML file with the original English text

#### Scenario: Attribute replacement
- **WHEN** a Vue template contains a hardcoded attribute like `placeholder="Search..."`
- **THEN** it SHALL be replaced with `:placeholder="$t('common.search')"`
- **THEN** the translation key SHALL exist in the YAML file

### Requirement: Coffee constants refactoring
The `helpers/coffee.ts` constants (grind sizes, equipment types, taste dimensions) SHALL be refactored to separate data values from display labels. Each constant entry SHALL have a `value` field (for backend data) and a `labelKey` field (translation key for display).

#### Scenario: Constant structure
- **WHEN** code imports `GRIND_SIZES` from `helpers/coffee.ts`
- **THEN** each entry SHALL be an object with `value` and `labelKey` fields
- **THEN** the `value` field SHALL contain the original enum string (e.g., `'Extra Fine'`)
- **THEN** the `labelKey` field SHALL contain a translation key (e.g., `'coffee.grindSizes.extraFine'`)

#### Scenario: Display in components
- **WHEN** a component displays a grind size label
- **THEN** it SHALL use `$t(grindSize.labelKey)` to retrieve the translated text
- **THEN** the `value` field SHALL remain unchanged for backend data submission

### Requirement: Router meta translation
Router configuration SHALL use `meta.title` with translation keys instead of hardcoded English text. The document title composition logic SHALL translate the key to the current locale.

#### Scenario: Route definition
- **WHEN** a route is defined in `router/index.ts`
- **THEN** it SHALL use `meta: { title: 'nav.gallery' }` instead of `meta: { title: 'Gallery' }`
- **THEN** the translation key SHALL exist in the YAML file

#### Scenario: Document title composition
- **WHEN** the router navigates to a page with `meta.title` containing a translation key
- **THEN** the document title SHALL be composed as `${VITE_DEFAULT_TITLE} | ${t(meta.title)}`
- **THEN** the title SHALL update when the locale changes

### Requirement: Locale-aware date formatting
Date formatting SHALL use the current locale instead of hardcoded `en-GB`. The system SHALL use dayjs with dynamic locale imports or vue-i18n's datetime formatting.

#### Scenario: Date display in TastingCard
- **WHEN** TastingCard displays a tasting date
- **THEN** it SHALL use locale-aware formatting instead of `toLocaleDateString('en-GB', ...)`
- **THEN** the date format SHALL adapt to the current locale

### Requirement: Interpolation for dynamic values
Strings containing dynamic values (names, dates, counts) SHALL use vue-i18n interpolation syntax. The YAML file SHALL define placeholders like `{name}` or `{date}`.

#### Scenario: String interpolation
- **WHEN** a string contains a dynamic value like `Delete "Bean Name"?`
- **THEN** the YAML entry SHALL use `msgstr: 'Delete "{name}"?'`
- **THEN** the code SHALL call `$t('coffee.editor.confirmDeleteBean', { name: bean.name })`
- **THEN** the interpolated value SHALL be inserted at runtime

### Requirement: No hardcoded strings remain
After extraction, a search for hardcoded English strings in Vue templates and TypeScript files SHALL find zero matches (excluding comments, brand names, and technical constants).

#### Scenario: Verification search
- **WHEN** a developer searches for hardcoded strings in `src/**/*.vue` and `src/**/*.ts`
- **THEN** the search SHALL find no hardcoded UI strings in templates or user-facing code
- **THEN** only comments, brand names (e.g., "IDEA GROUND ZERO"), and technical constants (e.g., API endpoints) SHALL remain
