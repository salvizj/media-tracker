document.addEventListener("DOMContentLoaded", () => {
	// Form submissions
	setupForms()

	setupFilterByStatus()

	// Movie actions
	handleEditMovieForms()
	setupDeleteButtons("movies")

	// Manhwa and manga actions
	setupDeleteButtons("manhwa-and-manga")
	setupEditAndCancelButtons("manhwa-and-manga")

	// TV show actions
	setupDeleteButtons("tv-shows")
	setupEditAndCancelButtons("tv-shows")
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

function setupDeleteButtons(type) {
	document
		.querySelectorAll(`.${type}-action-btn[data-method='DELETE']`)
		.forEach((btn) => {
			btn.addEventListener("click", (e) => {
				e.preventDefault()
				const id = btn.closest("li[data-id]").getAttribute("data-id")
				console.log(id)
				let confirmationMessage =
					"Are you sure you want to delete this item?"
				if (type === "movies")
					confirmationMessage =
						"Are you sure you want to delete this movie?"
				else if (type === "tv-shows")
					confirmationMessage =
						"Are you sure you want to delete this TV show?"
				else if (type === "manhwa-and-manga")
					confirmationMessage =
						"Are you sure you want to delete this manhwa?"
				deleteItem(type, id, confirmationMessage)
			})
		})
}

function setupEditAndCancelButtons(type) {
	document
		.querySelectorAll(`.${type}-action-btn[data-method='PUT']`)
		.forEach((btn) => {
			btn.addEventListener("click", () => {
				const li = btn.closest(`li[data-id]`)
				const id = li.getAttribute("data-id")
				const editForm = document.querySelector(
					`.edit-${type}-form[data-id="${id}"]`
				)
				if (editForm) editForm.style.display = "flex"
				const btnDiv = li.querySelector("div")
				if (btnDiv) btnDiv.style.display = "none"
			})
		})

	document.querySelectorAll(`.cancel-edit-${type}-btn`).forEach((btn) => {
		btn.addEventListener("click", () => {
			const editForm = btn.closest(`.edit-${type}-form`)
			if (!editForm) return
			const id = editForm.getAttribute("data-id")
			const li = document.querySelector(`li[data-id="${id}"]`)
			if (editForm) editForm.style.display = "none"
			const btnDiv = li.querySelector("div")
			if (btnDiv) btnDiv.style.display = "flex"
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

function deleteItem(type, id, confirmationMessage) {
	if (confirm(confirmationMessage)) {
		fetch(`api/${type}/${id}`, {
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
function setupFilterByStatus() {
	document.querySelectorAll(".filter-btn").forEach((btn) => {
		btn.addEventListener("click", function () {
			const status = this.getAttribute("data-status")
			document.querySelectorAll("ul li[data-status]").forEach((li) => {
				if (
					status === "All" ||
					li.getAttribute("data-status") === status
				) {
					li.style.display = ""
				} else {
					li.style.display = "none"
				}
			})
		})
	})
}
