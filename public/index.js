document.addEventListener("DOMContentLoaded", () => {
	const movieForm = document.getElementById("movieForm")
	if (movieForm) {
		movieForm.addEventListener("submit", async (e) => {
			e.preventDefault()
			const name = movieForm.name.value
			try {
				const response = await fetch("/api/movies", {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ name }),
				})
				if (!response.ok) throw new Error("Request failed")
				movieForm.reset()
				location.reload()
			} catch (error) {
				console.error(error)
			}
		})
	}

	const tvShowForm = document.getElementById("tvShowForm")
	if (tvShowForm) {
		tvShowForm.addEventListener("submit", async (e) => {
			e.preventDefault()
			const name = tvShowForm.name.value
			const status = tvShowForm.status.value
			const season = Number(tvShowForm.season.value)
			const episode = Number(tvShowForm.episode.value)
			try {
				const response = await fetch("/api/tv-shows", {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ name, status, season, episode }),
				})
				if (!response.ok) throw new Error("Request failed")
				tvShowForm.reset()
				location.reload()
			} catch (error) {
				console.error(error)
			}
		})
	}

	const manhwaForm = document.getElementById("manhwaForm")
	if (manhwaForm) {
		manhwaForm.addEventListener("submit", async (e) => {
			e.preventDefault()
			const name = manhwaForm.name.value
			const status = manhwaForm.status.value
			const chapter = Number(manhwaForm.chapter.value)
			try {
				const response = await fetch("/api/manhwa-and-manga", {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ name, status, chapter }),
				})
				if (!response.ok) throw new Error("Request failed")
				manhwaForm.reset()
				location.reload()
			} catch (error) {
				console.error(error)
			}
		})
	}
})
