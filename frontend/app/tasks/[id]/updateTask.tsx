"use client"

import { useRouter } from "next/navigation"
import { updateTask } from "@/app/lib/tasks"
import { useState, useEffect } from "react"
import { Task } from "@/app/types/task"
const API_URL = "http://localhost:8000"

type Props = {
    id: string
    onUpdated: () => void
}

export default function UpdateTask({id, onUpdated}: Props) {
    const [updateTitle, setUpdateTitle] = useState("")
    const [updateStatus, setUpdateStatus] = useState("")
    const [updatedTask, setUpdatedTask] = useState<Task | null>(null)
    const [updateTaskStatus, setUpdateTaskStatus] = useState("")

    type updateTaskStr = {
      title?: string
      status?: string
    }

    const body: updateTaskStr = {}
    if (updateTitle) body.title = updateTitle
    if (updateStatus) body.status = updateStatus

    async function handlerUpdateTask() {
        setUpdatedTask(null)
        setUpdateTaskStatus("")
        if (updateTitle == "" && updateStatus == "") {
            setUpdateTaskStatus("Please enter updated title or status")
            return
        }
        const data: Task = await updateTask(id, body)
        setUpdatedTask(data)
        setUpdateTaskStatus("success")
        onUpdated()
    }

    return (
      <div>
        <h2>Update Task</h2>

        <input
          type="text"
          value={updateTitle}
          onChange={(e) => setUpdateTitle(e.target.value)}
          placeholder="Task Title"
        />
        <select
            value={updateStatus}
            onChange={(e) => setUpdateStatus(e.target.value)}
            className="border border-gray-300 rounded-md px-3 py-2 bg-white
            focus:outline-none focus:ring-2 focus:ring-green-500"
        >
            <option value="">Choose Status</option>
            <option value={"todo"}>To do</option>
            <option value="inprogress">In Progress</option>
            <option value="blocked">Blocked</option>
            <option value={"onhold"}>On Hold</option>
            <option value={"done"}>Done</option>
            <option value={"cancelled"}>Cancelled</option>
        </select>

        <button onClick={handlerUpdateTask}>
            Updated Task
        </button>

        {updatedTask && (
          <div>
            <h2>Updated Task</h2>
            <p>{updatedTask.id}</p>
            <p>{updatedTask.title}</p>
            <p>{updatedTask.status}</p>
          </div>
        )}

        {updateTaskStatus && (
            <p>{updateTaskStatus}</p>
        )}
      </div>
    )
  }