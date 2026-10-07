"use client"

import { useState } from "react"
import { createTask } from "@/app/lib/tasks"
import { Task } from "@/app/types/task"
const API_URL = "http://localhost:8000"

type Props = {
    onCreated: () => void
}

export default function CreateTask({ onCreated }: Props) {
    const [createTitle, setCreateTitle] = useState("")
    const [createStatus, setCreateStatus] = useState("")
    const [createTaskStatus, setCreateTaskStatus] = useState("")
    const [createdTask, setCreatedTask] = useState<Task | null>(null)

    async function handlerCreateTask() {
        if (createTitle == "" || createStatus == "") {
            setCreateTaskStatus("Please include title and status")
        }
        else {
            const data: Task = await createTask(createTitle, createStatus)
            setCreatedTask(data)
            setCreateTaskStatus("success")
            onCreated()
        }
    }

    return (
    <div>
        <h2>Create Task</h2>

        <input
        type="text"
        value={createTitle}
        onChange={(e) => setCreateTitle(e.target.value)}
        placeholder="Task Title"
        />
        <select
            value={createStatus}
            onChange={(e) => setCreateStatus(e.target.value)}
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

        <button onClick={handlerCreateTask}>
            Create Task
        </button>

        {createTaskStatus == "success" && (
        <div>
            <h3>Task Created Successfully!</h3>
            <p>{createdTask?.id}</p>
            <p>{createdTask?.title}</p>
            <p>{createdTask?.status}</p>
        </div>
        )}

        {createTaskStatus != "success" && (
        <div>
            <h3>{createTaskStatus}</h3>
        </div>
        )}


    </div>
    )
}