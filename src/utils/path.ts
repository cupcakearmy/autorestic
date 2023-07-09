function parseKey(key: any) {
  const asNumber = parseInt(key)
  const isString = isNaN(asNumber)
  return [isString ? key : asNumber, isString]
}

export function setByPath(source: object, path: string, value: unknown) {
  const segments = path.split('.')
  const last = segments.length - 1
  let node: any = source
  for (const [i, segment] of segments.entries()) {
    const [key, isString] = parseKey(segment)
    if (Array.isArray(node) && isString) throw new Error(`array require a numeric index`)
    if (typeof node !== 'object') throw new Error(`could not set path "${segment}" on ${node}.`)
    if (i === last) {
      node[key] = value
    } else {
      const [_, isNextString] = parseKey(segments[i + 1])
      if (node[key] === undefined) node[key] = isNextString ? {} : []
      node = node[key]
    }
  }
  return source
}
