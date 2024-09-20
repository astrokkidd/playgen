const input_tag = document.querySelector(".input-tag")

let tags_list = []

input_tag.addEventListener("keyup", (e) => {
	const val = input_tag.value
	if (e.key == "Enter") {
		if (tags_list.length >= 5) return;
		if (tags_list.some(e => e.text == val)) return;
		if (val == "") return;
		
		var flag = 0;
		
		availableSeedGenres.genres.forEach(genre => {
			if (val == genre) {
				flag = 1;	
			}
		})

		if (flag != 1) return;



		const tags = val.split(',').map(e => e.trim()).filter(e => e !== "")

		for (let i of tags) {
			tags_list.push({
				id: i,
				text: i
			})
		}
		input_tag.value = ""
		RenderTags()
	}
})

window.addEventListener("load", (e) => {
	const available_seed_genres = document.querySelector(".available-seed-genres")
	let cache = ""

	availableSeedGenres.genres.forEach(e => {
		cache = `<option value = "${e}">`;
		available_seed_genres.insertAdjacentHTML("afterbegin", cache)
	})

})


function RenderTags() {
	const genre_tags = document.querySelector(".tags")
	let cache = ""

	document.querySelectorAll(".tag").forEach(e => e.remove())
	cache = ""

	tags_list.forEach(e => {
		cache = `<div class = "tag">
				<span>${e.text}</span>
				<button type="button" data-id="${e.id}"class="btn-rm-tag">
					<img src="./x.svg"/>
				</button>
			</div>`;
		genre_tags.insertAdjacentHTML("afterbegin", cache)
		HandleRmTags()
	})
}

function HandleRmTags() {
	const btns = document.querySelectorAll('.btn-rm-tag')
	if (btns.length > 0) {
		btns.forEach(e => {
			e.onclick = function () {
				const data_id = e.getAttribute('data-id')
				tags_list = tags_list.filter(x => x.id != data_id)
				RenderTags();
			}
		})
	}

}
