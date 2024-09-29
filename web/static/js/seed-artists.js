document.addEventListener('DOMContentLoaded', function() {
    const input_tag = document.querySelector('.input-tag');
    const tagContainer = document.querySelector('.tags');
    const dataList = document.getElementById('available-seed-artists');

    let selectedArtists = [];

    // Spotify search API call
    async function searchSpotifyArtists(query) {
        const response = await fetch(`/search-artists?q=${query}`);
        const artists = await response.json();
        return artists;
    }
    

});
