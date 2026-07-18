import fs from 'node:fs'
import YAML from 'yaml'

export default function gettextYamlI18nPlugin() {
  return {
    name: 'gettext-yaml-i18n',
    transform(code: string, id: string) {
      if (!id.endsWith('.yaml') || !id.includes('/locales/')) return null

      const content = fs.readFileSync(id, 'utf-8')
      const parsed = YAML.parse(content)

      const result: Record<string, string> = {}

      function traverse(obj: Record<string, unknown>) {
        for (const [, value] of Object.entries(obj)) {
          if (typeof value === 'object' && value !== null) {
            const objValue = value as Record<string, unknown>

            if ('msgstr' in objValue && 'msgid' in objValue) {
              const msgid = objValue.msgid as string
              const msgstr = objValue.msgstr as string
              result[msgid] = msgstr
            } else {
              traverse(objValue)
            }
          }
        }
      }

      traverse(parsed)

      return {
        code: `export default ${JSON.stringify(result)}`,
        map: null
      }
    }
  }
}
