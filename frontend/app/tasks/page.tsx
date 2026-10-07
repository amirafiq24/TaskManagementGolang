"use client"

import TasksTable from "./tasksTable" 
import CreateTask  from "./createTask"
import { useEffect, useState } from "react"
import { Task } from "@/app/types/task"
import { getTasks } from "@/app/lib/tasks"
import { useRouter } from "next/navigation"
const API_URL = "http://localhost:8000"

export default function TasksPage() {
    const [tasks, setTasks] = useState<Task[]>([])
    const router = useRouter()

    async function loadTasks() {
        let data = await getTasks()
        if (!data) {
            data = []
        }
        setTasks(data)
    }   

    useEffect(() => {loadTasks()}, [])

    return (
        <div>
            <TasksTable tasks={tasks}/>  
            <CreateTask onCreated={loadTasks}/>
            <button onClick={() => router.push('/')}>Back</button>
        </div>
    )
}