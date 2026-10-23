export const BACKGROUND_STYLES: Record<string, string> = {
  default: 'bg-background/75',
  red: 'bg-red-700/75',
  green: 'bg-emerald-600/75',
  blue: 'bg-sky-600/75',
  darkblue: 'bg-blue-900/75',
  orange: 'bg-orange-500/75',
  sweden: 'bg-[#005293]/75',
  white: 'bg-slate-500/75'
}

export const TEXT_STYLES: Record<string, string> = {
  sweden: 'text-[#FFCD00]'
}

export const COLOR_KEYS = Object.keys(BACKGROUND_STYLES)
export const TEXT_COLOR_KEYS = Object.keys(TEXT_STYLES)

export function parseTailwindColor(type: string, color: string = 'default'): string {
  if (type === 'text') {
    return TEXT_STYLES[color] ?? ''
  }

  return BACKGROUND_STYLES[color] ?? BACKGROUND_STYLES['default']!
}
