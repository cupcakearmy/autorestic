export function wait(seconds: number): Promise<never> {
  return new Promise((resolve) => setTimeout(resolve, seconds * 1000))
}
