document.addEventListener("DOMContentLoaded", () => {
	// Form submissions
	setupForms()

	// Movie actions
	setupDeleteButtonsForMovies()
	setupEditAndCancelButtonsForMovieEditForm()
	handleEditMovieForms()

	// Manhwa and manga actions

	// TV show actions
})

function setupForms() {
	handleFormSubmit("movieForm", "/api/movies", { name: "string" })
	handleFormSubmit("tvShowForm", "/api/tv-shows", {
		name: "string",
		status: "string",
		season: "number",
		episode: "number",
	})
	handleFormSubmit("manhwaForm", "/api/manhwa-and-manga", {
		name: "string",
		status: "string",
		chapter: "number",
	})
}

function setupDeleteButtonsForMovies() {
	document
		.querySelectorAll(".movie-action-btn[data-method='DELETE']")
		.forEach((btn) => {
			btn.addEventListener("click", (e) => {
				e.preventDefault()
				const li = btn.closest("li")
				const id = li.getAttribute("data-id")
				deleteMovie(id)
			})
		})
}

function setupEditAndCancelButtonsForMovieEditForm() {
	document
		.querySelectorAll(".movie-action-btn[data-method='PUT']")
		.forEach((btn) => {
			btn.addEventListener("click", () => {
				const li = btn.closest("li[data-id]")
				const id = li.getAttribute("data-id")
				const editForm = document.querySelector(
					`.edit-movie-form[data-id="${id}"]`
				)
				if (editForm) {
					editForm.style.display = "flex"
				}
				if (li) {
					const movieActionBtnDiv = li.querySelector("div")
					if (movieActionBtnDiv) {
						movieActionBtnDiv.style.display = "none"
					}
				}
			})
		})

	document.querySelectorAll(".cancel-edit-movie-btn").forEach((btn) => {
		btn.addEventListener("click", () => {
			const editForm = btn.closest(".edit-movie-form")
			if (!editForm) return
			const id = editForm.getAttribute("data-id")
			const li = document.querySelector(`li[data-id="${id}"]`)
			if (editForm) {
				editForm.style.display = "none"
			}
			if (li) {
				const movieActionBtnDiv = li.querySelector("div")
				if (movieActionBtnDiv) {
					movieActionBtnDiv.style.display = "flex"
				}
			}
		})
	})
}

function handleFormSubmit(formId, apiPath, fields, method = "POST") {
	const form = document.getElementById(formId)
	if (!form) return

	form.addEventListener("submit", async (e) => {
		e.preventDefault()
		const body = {}

		for (const [field, type] of Object.entries(fields)) {
			let value = form[field]?.value
			if (value === undefined) continue

			if (type === "number") {
				value = Number(value)
			}

			body[field] = value
		}

		try {
			const response = await fetch(apiPath, {
				method,
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(body),
			})
			if (!response.ok) throw new Error("Request failed")
			form.reset()
			location.reload()
		} catch (error) {
			console.error(error)
		}
	})
}

function deleteMovie(id) {
	if (confirm("Are you sure you want to delete this movie?")) {
		fetch(`api/movies/${id}`, {
			method: "DELETE",
		})
			.then(() => location.reload())
			.catch((error) => console.error(error))
	}
}

function handleEditMovieForms() {
	document.querySelectorAll(".edit-movie-form").forEach((form) => {
		form.addEventListener("submit", async (e) => {
			e.preventDefault()

			const id = form.id.value
			const name = form.name.value
			const date = form.date.value

			try {
				const response = await fetch(`/api/movies/${id}`, {
					method: "PUT",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ name, date }),
				})
				if (!response.ok) throw new Error("Request failed")
				location.reload()
			} catch (error) {
				console.error(error)
			}
		})
	})
}
