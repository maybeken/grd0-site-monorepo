import { cp, mkdir, rm } from 'node:fs/promises'
import { createRequire } from 'node:module'
import { dirname, join } from 'node:path'

const require = createRequire(import.meta.url)
const entry = require.resolve('@excalidraw/excalidraw')
const source = join(dirname(entry), 'fonts')
const target = join(process.cwd(), 'public', 'fonts')

await rm(target, { recursive: true, force: true })
await mkdir(target, { recursive: true })
await cp(source, target, { recursive: true })
