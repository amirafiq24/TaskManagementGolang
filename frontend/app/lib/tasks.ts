const API_URL = "http://localhost:8000"
import { Task } from "@/app/types/task"

export async function checkHealth() {
    const url = API_URL + "/health"

    const response = await fetch(url)
    const data = await response.json()
    return data
}

export async function getTasks() {
    const url = API_URL + "/tasks"
    const response = await fetch(url)
    const data = await response.json()

    return data
}

export async function getTask(id: string) {
    const url = API_URL + `/tasks/${id}`
    const response = await fetch(url)

    if (!response.ok) {
        return null
    }
    const data = await response.json()

    return data
} 

export async function createTask(title:string, status:string) {
    const createTaskObj = {
        title: title,
        status: status,
    }
    
    const url = API_URL + "/tasks"
    const response = await fetch(url, {
        method: `POST`,
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(createTaskObj),
    })

    const data = await response.json()
    return data
}

type UpdateTask = {
    title?: string
    status?: string
}

export async function updateTask(id:string, body:UpdateTask) {
    const url = API_URL + `/tasks/${id}`
    const response = await fetch(url, {
        method: `PATCH`,
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(body),
    })

    const data: Task = await response.json()
    return data
}

export async function deleteTask(id:string) {
    const url = API_URL + `/tasks/${id}`
    const response = await fetch(url, {
      method: `DELETE`,
    })

    console.log(response)

    return response.ok
}