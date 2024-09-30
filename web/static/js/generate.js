document.querySelector("button[type='submit']").addEventListener("click", function(event) {
    event.preventDefault(); // Prevent form submission

    // Collect values from sliders
    const acousticness = parseFloat(document.getElementById("acousticness").value / 100);
    const danceability = parseFloat(document.getElementById("danceability").value / 100);
    const energy = parseFloat(document.getElementById("energy").value / 100);
    const instrumentalness = parseFloat(document.getElementById("instrumentalness").value / 100);
    const popularity = parseInt(document.getElementById("popularity").value);
    const valence = parseFloat(document.getElementById("valence").value / 100);
    const limit = document.getElementById("max-tracks-slider").value;

    const genreButtons = document.querySelectorAll('button.btn-rm-tag.genre');
    const genreIds = Array.from(genreButtons).map(button => button.getAttribute('data-id'));
    const genreStrings = genreIds.join(',');

    const tracks_list = document.querySelector(".tracks-list");
    const add_button = document.querySelector(".add-button");
    const playlist_header = document.querySelector(".playlist-header");

    // Create a JSON object with the form data
    const formData = {
        target_acousticness: acousticness,
        target_danceability: danceability,
        target_energy: energy,
        target_instrumentalness: instrumentalness,
        target_popularity: popularity,
        target_valence: valence,
        limit: limit,
        seed_genres: genreStrings
    };

    // Send the data to the Go server
    fetch('/generate', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
    })
    .then(response => response.json())
    .then(data => {
        console.log('Playlist generated:', data.tracks);
        // Do something with the server's response
        //for each track in data.tracks

        data.tracks.forEach(track => {
            cache = `<li class="track">
                            <audio controls controls-list="play" src="${track.preview_url}">
                                <img src="./media/play.svg"/>
                            </audio>
	                    <img class="album" src="${track.album.images[0].url}" width="50" height="50" aria-label="${track.name} Album Art">
	                    <span class="name">${track.name}</span>
	                    <span class="artist">${track.artists.map((artist) => artist.name).join(', ')}</span>
                    </li>`;
                    tracks_list.insertAdjacentHTML("afterbegin", cache);
        })

        cache = `<button type="submit">ADD TO SPOTIFY</button>`;
        add_button.insertAdjacentHTML("afterbegin", cache);

        cache = `<div class="playlist-image">
			<img src="${data.tracks[0].album.images[0].url}" width="70" height="70">
			<img src="${data.tracks[1].album.images[0].url}" width="70" height="70">
			<img src="${data.tracks[2].album.images[0].url}" width="70" height="70">
			<img src="${data.tracks[3].album.images[0].url}" width="70" height="70">
		</div>
		<input type="text" class="playlist-title" value="My Generated Playlist">`;
        playlist_header.insertAdjacentHTML("afterbegin", cache);

        //create html ul
        //image, name, artist(s)
    })
    .catch(error => console.error('Error:', error));
});
