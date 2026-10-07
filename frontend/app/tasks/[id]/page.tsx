"use client"

import { useRouter } from "next/navigation"
import { useState, useEffect } from "react"
import { getTask, getTasks } from "@/app/lib/tasks"
import { Task } from "@/app/types/task"
import { useParams } from "next/navigation"
import UpdateTask from "./updateTask"
import DeleteTask from "./deleteTask"

const API_URL = "http://localhost:8000"

export default function TasksPage() {
    const [task, setTask] = useState<Task | null>(null)
    const [getTaskStatus, setGetTaskStatus] = useState("")
    const router = useRouter()

    const params = useParams<{ id: string }>()  
    const id = params.id

    async function loadTask() {
        const data: Task = await getTask(id)
        setTask(data)

        if (!data) {
            setGetTaskStatus("Failed to fetch task")
        }
        else {
            setGetTaskStatus("success")
        }
    }

    useEffect(() => {loadTask()}, [])

    return (
        <div>
            {getTaskStatus == "success" && (
                <div>
                    <p>Title: {task?.title}</p>
                    <p>Status: {task?.status}</p>
                </div>
            )}

            {getTaskStatus != "success" && (
                <div>
                    <p>{getTaskStatus}</p>
                </div>
            )}

            <UpdateTask id={id} onUpdated={loadTask}/>
            <DeleteTask id={id}/>

            <button onClick={() => router.push('/tasks')}>Back</button>
        </div>
    )
}