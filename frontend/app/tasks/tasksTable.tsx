"use client"

import { useRouter } from "next/navigation"
import { Task } from "@/app/types/task"
const API_URL = "http://localhost:8000"

type Props = {
    tasks: Task[]
}

export default function TasksTable({tasks}: Props) {
    const router = useRouter()

    return (
      <div>
        <h1 className="text-xl font-semibold text-left mb-4">List of all Tasks</h1>

        <table className="w-full table-fixed border-collapse">
            <thead>
                <tr>
                    <th className="text-left px-4 py-2">ID</th>
                    <th className="text-left px-4 py-2">Title</th>
                    <th className="text-left px-4 py-2">Status</th>
                    <th></th>
                </tr>
            </thead>

            <tbody>
                {tasks.map((task) => (
                    <tr key={task.id} className="border-t">
                        <td className="px-4 py-2">{task.id}</td>
                        <td className="px-4 py-2">{task.title}</td>
                        <td className="px-4 py-2">{task.status}</td>
                        <td className="align-middle">
                            <button onClick={() => router.push(`/tasks/${task.id}`)}>Open</button>
                        </td>
                    </tr>
                ))}
            </tbody>
        </table>

      </div>
    )
  
}