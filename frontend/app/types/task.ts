export type Task = {
  id: number
  title: string
  status: string
}

export type DependencyStatus = {
  status: string
  error: string
}

export type HealthStatus = {
  status: string
  uptime: string
  dependencies: Record<string, DependencyStatus>
}