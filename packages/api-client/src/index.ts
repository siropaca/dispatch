import type { components } from './schema'

export type Health = components['schemas']['Health']

// 公開 REST(openapi.yaml)の生成型を使った最小クライアント。
// baseUrl 省略時は同一オリジン(dev は Vite プロキシ、prod は api 同居)。
export async function getHealth(baseUrl = ''): Promise<Health> {
  const res = await fetch(`${baseUrl}/healthz`)
  if (!res.ok) {
    throw new Error(`healthz failed: ${res.status}`)
  }
  return (await res.json()) as Health
}
