import { mkdir, rm, writeFile } from 'node:fs/promises'
import { zodToJsonSchema } from 'zod-to-json-schema'
import { ConfigSchema } from '../src/config/schema/config'

const OUTPUT = './schema'

await rm(OUTPUT, { recursive: true, force: true })
await mkdir(OUTPUT, { recursive: true })

const Schemas = {
  config: ConfigSchema,
}

for (const [name, schema] of Object.entries(Schemas)) {
  const jsonSchema = zodToJsonSchema(schema, 'mySchema')
  await writeFile(`${OUTPUT}/${name}.json`, JSON.stringify(jsonSchema, null, 2), { encoding: 'utf-8' })
}
