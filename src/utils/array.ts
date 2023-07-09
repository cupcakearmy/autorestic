export function asArray<T>(singleOrArray: T | T[]): T[] {
  return Array.isArray(singleOrArray) ? singleOrArray : [singleOrArray]
}

export function isSubset<T>(subset: T[], set: T[]): boolean {
  return subset.every((v) => set.includes(v))
}
