document.querySelector("button[type='submit']").addEventListener("click", function(event) {
    event.preventDefault(); // Prevent form submission

    // Collect values from sliders
    const acousticness = parseInt(document.getElementById("acousticness").value / 100, 10);
    const danceability = parseInt(document.getElementById("danceability").value / 100, 10);
    const energy = parseInt(document.getElementById("energy").value / 100, 10);
    const instrumentalness = parseInt(document.getElementById("instrumentalness").value / 100, 10);
    const popularity = parseInt(document.getElementById("popularity").value / 100, 10);
    const valence = parseInt(document.getElementById("valence").value / 100, 10);
    const limit = document.getElementById("max-tracks-slider").value;

    const genreButtons = document.querySelectorAll('button.btn-rm-tag.genre');
    const genreIds = Array.from(genreButtons).map(button => button.getAttribute('data-id'));
    const genreStrings = genreIds.join(',');

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
        console.log('Playlist generated:', data);
        // Do something with the server's response
        
        //for each track in data.tracks
        //create html ul
        //image, name, artist(s)
    })
    .catch(error => console.error('Error:', error));
});
