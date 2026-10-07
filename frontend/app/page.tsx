"use client"

import { useRouter } from "next/navigation"
import { useState } from "react"
import "./app.css"
import { checkHealth, getTask, getTasks, createTask, updateTask, deleteTask } from "./lib/tasks"
import { HealthStatus } from "@/app/types/task"

const API_URL = "http://localhost:8000"

export default function Home() {
  const router = useRouter()

  function CheckHealth() {
    const [health, setHealth] = useState<HealthStatus | null>(null)

    async function handlerCheckHealth() {
      const data = await checkHealth()
      setHealth(data)
    }

    return (
      <div>
        <h3>Backend Health Check</h3>
        <button onClick={handlerCheckHealth}>Chech Health</button>
        <p>Status: {health?.status}</p>
        <p>Uptime: {health?.uptime}</p>
        <p>Dependencies</p>
        {health?.dependencies && (
          Object.entries(health.dependencies).map(([name, dep]) => (
            <div key={name}>
              <p>{name}</p>
              <p>Status: {dep.status}</p>
              {dep.error && (
                <p>Error: {dep.error}</p>
              )}
            </div>
          ))
        )}
      </div>

    )
  }


  return (
    <main>

      <h1>Task Management Systems</h1>

      <button onClick={() => router.push('/tasks')}>See all Tasks</button>
      <CheckHealth />
      
    </main>
  )
}