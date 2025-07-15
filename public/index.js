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
	setupLogout()
	updateNavVisibility()
	setupUserIdValuePlaceholders()
	setupBulkUploadForms()
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
	handleFormSubmission("loginForm", "/login", "POST", null, true)
	handleFormSubmission("registerForm", "/register", "POST", null, true)
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

function handleFormSubmission(
	formOrFormId,
	apiPath,
	method = "POST",
	apiPathID = null,
	authSubmission = false
) {
	const form =
		typeof formOrFormId === "string"
			? document.getElementById(formOrFormId)
			: formOrFormId
	if (!form) return

	form.addEventListener("submit", async (e) => {
		e.preventDefault()
		const body = {}
		new FormData(form).forEach((value, key) => {
			body[key] = isNaN(value) ? value : Number(value)
		})
		const url = apiPathID ? `${apiPath}/${apiPathID}` : apiPath
		try {
			const response = await fetch(url, {
				method,
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(body),
			})

			if (!response.ok) {
				const errorData = await response.json()
				throw new Error(errorData.error || "Request failed")
			}

			const data = await response.json()

			if (authSubmission) {
				showMessage(data.message || "Success!", "message")
				setTimeout(() => {
					window.location.href = "/"
				}, 1000)
			} else {
				form.reset()
				location.reload()
			}
		} catch (error) {
			showMessage(error.message, "error")
		}
	})
}

function showMessage(message, type = "message") {
	const messageBox = document.getElementById("message-box")
	if (!messageBox) return

	messageBox.textContent = message

	messageBox.classList.remove(
		"bg-green-500",
		"bg-red-500",
		"bg-blue-500",
		"text-white"
	)

	if (type === "message") {
		messageBox.classList.add("bg-green-500", "text-white")
	} else if (type === "error") {
		messageBox.classList.add("bg-red-500", "text-white")
	}

	messageBox.classList.remove("hidden")

	setTimeout(() => {
		messageBox.classList.add("hidden")
	}, 3000)
}

function getCurrentUserID() {
	return getCookie("user_id")
}

function setupLogout() {
	const logoutBtn = document.getElementById("logout-btn")
	if (!logoutBtn) return

	logoutBtn.addEventListener("click", async (e) => {
		fetch("/logout", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
		})
			.then((response) => response.json())
			.then((data) => {
				showMessage(data.message, "message")
				setTimeout(() => {
					window.location.href = "/login"
				}, 1000)
			})
			.catch((error) => {
				showMessage(error.message, "error")
			})
	})
}

function getCookie(name) {
	const value = `; ${document.cookie}`
	const parts = value.split(`; ${name}=`)
	if (parts.length === 2) return parts.pop().split(";").shift()
	return null
}
function isLoggedIn() {
	return getCookie("session_id") !== null && getCookie("user_id") !== null
}

function initializeAllEditForms() {
	document.querySelectorAll(".edit-movies-form").forEach((form) => {
		handleFormSubmission(
			form,
			"/api/movies",
			"PUT",
			form.querySelector('input[name="id"]').value
		)
	})
	document.querySelectorAll(".edit-tv-shows-form").forEach((form) => {
		handleFormSubmission(
			form,
			"/api/tv-shows",
			"PUT",
			form.querySelector('input[name="id"]').value
		)
	})
	document.querySelectorAll(".edit-manhwa-and-manga-form").forEach((form) => {
		handleFormSubmission(
			form,
			"/api/manhwa-and-manga",
			"PUT",
			form.querySelector('input[name="id"]').value
		)
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

function updateNavVisibility() {
	document
		.getElementById("nav-logged-in")
		.classList.toggle("hidden", !isLoggedIn())
	document
		.getElementById("nav-logged-out")
		.classList.toggle("hidden", isLoggedIn())
}
function setupUserIdValuePlaceholders() {
	document.querySelectorAll(".user-id-input").forEach((input) => {
		input.value = getCookie("user_id")
	})
}

function setupBulkTvShows() {
	const form = document.getElementById("bulkTvShowsForm")
	if (!form) return
	form.addEventListener("submit", async (e) => {
		e.preventDefault()
		const textarea = document.getElementById("bulk_tv_shows")
		const lines = textarea.value
			.split("\n")
			.map((l) => l.trim())
			.filter((l) => l)
		const items = lines
			.map((line) => {
				const parts = line.split("|").map((p) => p.trim())
				if (parts.length < 5) return null
				return {
					name: parts[0],
					status: parts[1],
					season: Number(parts[2]),
					episode: Number(parts[3]),
					date: parts[4],
				}
			})
			.filter(Boolean)
		try {
			const res = await fetch("/bulk-add/tv-shows", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(items),
			})
			const data = await res.json()
			showMessage(
				data.message || "Bulk add complete!",
				res.ok ? "message" : "error"
			)
			if (res.ok) textarea.value = ""
			setTimeout(() => {
				location.reload()
			}, 1000)
		} catch (err) {
			showMessage(err.message || "Bulk add failed", "error")
		}
	})
}

function setupBulkManhwa() {
	const form = document.getElementById("bulkManhwaAndMangaForm")
	if (!form) return
	form.addEventListener("submit", async (e) => {
		e.preventDefault()
		const textarea = document.getElementById("bulk_manhwa_and_manga")
		const lines = textarea.value
			.split("\n")
			.map((l) => l.trim())
			.filter((l) => l)
		const items = lines
			.map((line) => {
				const parts = line.split("|").map((p) => p.trim())
				if (parts.length < 4) return null
				return {
					name: parts[0],
					status: parts[1],
					chapter: Number(parts[2]),
					date: parts[3],
				}
			})
			.filter(Boolean)
		try {
			const res = await fetch("/bulk-add/manhwa-and-manga", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(items),
			})
			const data = await res.json()
			showMessage(
				data.message || "Bulk add complete!",
				res.ok ? "message" : "error"
			)
			if (res.ok) textarea.value = ""
			setTimeout(() => {
				location.reload()
			}, 1000)
		} catch (err) {
			showMessage(err.message || "Bulk add failed", "error")
		}
	})
}

function setupBulkMovies() {
	const form = document.getElementById("bulkMoviesForm")
	if (!form) return
	form.addEventListener("submit", async (e) => {
		e.preventDefault()
		const textarea = document.getElementById("bulk_movies")
		const lines = textarea.value
			.split("\n")
			.map((l) => l.trim())
			.filter((l) => l)
		const items = lines
			.map((line) => {
				const parts = line.split("|").map((p) => p.trim())
				if (parts.length < 2) return null
				return {
					name: parts[0],
					date: parts[1],
				}
			})
			.filter(Boolean)
		try {
			const res = await fetch("/bulk-add/movies", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify(items),
			})
			const data = await res.json()
			showMessage(
				data.message || "Bulk add complete!",
				res.ok ? "message" : "error"
			)
			if (res.ok) textarea.value = ""
			setTimeout(() => {
				location.reload()
			}, 1000)
		} catch (err) {
			showMessage(err.message || "Bulk add failed", "error")
		}
	})
}

function setupBulkUploadForms() {
	setupBulkTvShows()
	setupBulkManhwa()
	setupBulkMovies()
}
