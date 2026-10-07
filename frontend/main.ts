async function chechHealth() {
    const response = await fetch("http://localhost:8000/health")

    const data = await response.json()
    console.log(data)
}

chechHealth()