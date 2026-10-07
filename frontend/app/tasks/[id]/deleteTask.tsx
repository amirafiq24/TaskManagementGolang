"use client"

import { useRouter } from "next/navigation"
import { updateTask } from "@/app/lib/tasks"
import { useState, useEffect } from "react"
import { Task } from "@/app/types/task"
import { deleteTask } from "@/app/lib/tasks"
const API_URL = "http://localhost:8000"

type Props = {
    id: string
}

export default function DeleteTask({id}: Props) {
    const [deleteTaskStatus, setDeleteTaskStatus] = useState("")
    const router = useRouter()

    async function handlerDeleteTask() {
      setDeleteTaskStatus("")
      const data = await deleteTask(id)
      if (data) {
        setDeleteTaskStatus("Deleted")
        router.push("/tasks")
      }
      else {
        setDeleteTaskStatus("Failed to delete task")
      }
    }

    return (
      <div>
        <h2>Delete Task</h2>

        <button onClick={handlerDeleteTask} >
            Delete Task
        </button>

        <p>{deleteTaskStatus}</p>
      </div>
    )
    
  }