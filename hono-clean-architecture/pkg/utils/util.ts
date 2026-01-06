export function randomJitter(ttlMs: number): number {
  return ttlMs + Math.floor(Math.random() * 60000);
}
