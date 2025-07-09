const form = document.getElementById("movieForm")

form.addEventListener("submit", async (e) => {
	e.preventDefault()

	const name = form.name.value

	try {
		const response = await fetch("/api/movies", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({ name }),
		})

		if (!response.ok) {
			throw new Error("Request failed")
		}

		form.reset()
	} catch (error) {
		console.error(error)
	}
})
