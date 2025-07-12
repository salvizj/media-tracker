document.addEventListener("DOMContentLoaded", () => {
	initializeForms()
	initializeSearchForms()
	initializeStatusFilter()
	initializeDeleteButtons()
	initializeEditAndCancelButtons()
	initializeAllEditForms()
	initializeIncrementDecrementButtons(
		"tv-shows-increment-btn",
		"tv-shows-decrement-btn",
		"tv-shows-increment-and-decrement-form",
		"tv-shows",
		"episode"
	)
	initializeIncrementDecrementButtons(
		"manhwa-and-manga-increment-btn",
		"manhwa-and-manga-decrement-btn",
		"manhwa-and-manga-increment-and-decrement-form",
		"manhwa-and-manga",
		"chapter"
	)
})

function initializeEditAndCancelButtons() {
	setupEditAndCancelButtonsForType("movies")
	setupEditAndCancelButtonsForType("manhwa-and-manga")
	setupEditAndCancelButtonsForType("tv-shows")
}

function initializeDeleteButtons() {
	setupDeleteButtonsForType("movies")
	setupDeleteButtonsForType("manhwa-and-manga")
	setupDeleteButtonsForType("tv-shows")
}

function initializeForms() {
	handleFormSubmission("movieForm", "/api/movies", "POST")
	handleFormSubmission("tvShowForm", "/api/tv-shows", "POST")
	handleFormSubmission("manhwaAndMangaForm", "/api/manhwa-and-manga", "POST")
	handleFormSubmission("loginForm", "/api/login", "POST")
	handleFormSubmission("registerForm", "/api/register", "POST")
}

function initializeSearchForms() {
	setupSearchFormForType("movieSearchForm", "movie-item", "movie-item-name")
	setupSearchFormForType(
		"tvShowSearchForm",
		"tv-show-item",
		"tv-show-item-name"
	)
	setupSearchFormForType(
		"manhwaAndMangaSearchForm",
		"manhwa-and-manga-item",
		"manhwa-and-manga-item-name"
	)
}

function setupDeleteButtonsForType(type) {
	document
		.querySelectorAll(`.${type}-action-btn[data-method='DELETE']`)
		.forEach((btn) => {
			btn.addEventListener("click", (e) => {
				e.preventDefault()
				const id = btn.closest("li[data-id]").getAttribute("data-id")
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

function setupEditAndCancelButtonsForType(type) {
	document
		.querySelectorAll(`.${type}-action-btn[data-method='PUT']`)
		.forEach((btn) => {
			btn.addEventListener("click", () => {
				const li = btn.closest(`li[data-id]`)
				const editForm = li.querySelector(`.edit-${type}-form`)
				if (editForm) editForm.classList.remove("hidden")
				const btnDiv = li.querySelector("div")
				if (btnDiv) btnDiv.style.display = "none"
			})
		})

	document.querySelectorAll(`.${type}-cancel-edit-btn`).forEach((btn) => {
		btn.addEventListener("click", () => {
			const editForm = btn.closest(`.edit-${type}-form`)
			if (!editForm) return
			editForm.classList.add("hidden")
			const li = editForm.closest("li[data-id]")
			const btnDiv = li.querySelector("div")
			if (btnDiv) btnDiv.style.display = "flex"
		})
	})
}

function handleFormSubmission(formOrId, apiPath, method = "POST") {
	const form =
		typeof formOrId === "string"
			? document.getElementById(formOrId)
			: formOrId
	if (!form) return

	form.addEventListener("submit", async (e) => {
		e.preventDefault()
		const body = {}
		new FormData(form).forEach((value, key) => {
			body[key] = isNaN(value) ? value : Number(value)
		})
		let url = apiPath
		if (method === "PUT" && form.id && form.id.value) {
			url += `/${form.id.value}`
		}
		try {
			const response = await fetch(url, {
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

function initializeAllEditForms() {
	document.querySelectorAll(".edit-movies-form").forEach((form) => {
		handleFormSubmission(form, "/api/movies", { name: "string" }, "PUT")
	})
	document.querySelectorAll(".edit-tv-shows-form").forEach((form) => {
		handleFormSubmission(form, "/api/tv-shows", "PUT")
	})
	document.querySelectorAll(".edit-manhwa-and-manga-form").forEach((form) => {
		handleFormSubmission(form, "/api/manhwa-and-manga", "PUT")
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

function initializeStatusFilter() {
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
function setupSearchFormForType(formId, itemClass, itemNameClass) {
	const form = document.getElementById(formId)
	if (form) {
		const itemsWhole = document.querySelectorAll(`.${itemClass}`)
		const searchInput = form.querySelector('input[name="name"]')
		searchInput.addEventListener("input", function () {
			const searchVal = searchInput.value.trim().toLowerCase()

			itemsWhole.forEach((item) => {
				const itemName = item.querySelector(`.${itemNameClass}`)
				const itemNameText = itemName
					? itemName.textContent.toLowerCase()
					: ""

				if (itemNameText.includes(searchVal)) {
					item.style.display = ""
				} else {
					item.style.display = "none"
				}
			})
		})
	}
}

function initializeIncrementDecrementButtons(
	incrementBtnClass,
	decrementBtnClass,
	incrementAndDecrementFormClass,
	type,
	inputName
) {
	document
		.querySelectorAll(`.${incrementBtnClass}, .${decrementBtnClass}`)
		.forEach((btn) => {
			btn.addEventListener("click", async (e) => {
				e.preventDefault()
				const li = btn.closest("li[data-id]")
				const form = li.querySelector(
					`.${incrementAndDecrementFormClass}`
				)
				const input = form.querySelector(`input[name="${inputName}"]`)
				let currentValue = Number(input.value)

				if (btn.classList.contains(incrementBtnClass)) {
					currentValue += 1
				} else {
					currentValue = Math.max(1, currentValue - 1)
				}
				input.value = currentValue

				const body = {}
				new FormData(form).forEach((value, key) => {
					body[key] = isNaN(value) ? value : Number(value)
				})
				body[inputName] = currentValue

				const id = form.querySelector('input[name="id"]').value
				try {
					const response = await fetch(`/api/${type}/${id}`, {
						method: "PUT",
						headers: { "Content-Type": "application/json" },
						body: JSON.stringify(body),
					})
					if (!response.ok) throw new Error("Request failed")
					location.reload()
				} catch (error) {
					console.error(error)
				}
			})
		})
}
