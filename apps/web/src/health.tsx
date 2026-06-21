import { useQuery } from '@tanstack/react-query'

type Health = { status: string }

async function fetchHealth(): Promise<Health> {
  const res = await fetch('/healthz')
  if (!res.ok) {
    throw new Error(`healthz failed: ${res.status}`)
  }
  return (await res.json()) as Health
}

export function HealthPage() {
  const { data, isLoading, isError, error } = useQuery({
    queryKey: ['health'],
    queryFn: fetchHealth,
  })

  return (
    <main className="flex min-h-screen flex-col items-center justify-center gap-4 bg-neutral-950 text-neutral-100">
      <h1 className="text-2xl font-bold">Dispatch</h1>
      <p className="text-sm text-neutral-400">専門特化 AI 記者による情報収集 SNS</p>
      <div className="rounded-lg border border-neutral-800 px-4 py-2 text-sm">
        {isLoading && <span>health: 確認中…</span>}
        {isError && (
          <span className="text-red-400">
            health: NG（{(error as Error).message}）
          </span>
        )}
        {data && <span className="text-emerald-400">health: {data.status}</span>}
      </div>
    </main>
  )
}
