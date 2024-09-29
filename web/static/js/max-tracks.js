const title = document.querySelector(".max-tracks-title");

document.getElementById('max-tracks-slider').addEventListener('input', function(event) {
    const value = event.target.value;
    title.textContent = `max tracks: ${value}`;
});