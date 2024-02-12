function changeImage(src) {
    document.getElementById('main-product-image').src = src;
}

function sortStores() {
    var storesList = document.querySelector('.stores-list');
    var stores = document.querySelectorAll('.store');

    var sortBy = document.getElementById('sortby').value;
    var sortedStores = Array.from(stores);

    sortedStores.sort(function(a, b) {
        if (sortBy === 'cheapest') {
            return parseInt(a.getAttribute('data-price')) - parseInt(b.getAttribute('data-price'));
        }
        // Extend with other sorting criteria
    });

    // Clear the list and re-insert sorted items
    while (storesList.firstChild) {
        storesList.removeChild(storesList.firstChild);
    }

    sortedStores.forEach(function(node) {
        storesList.appendChild(node);
    });
}

// Add event listener to the select element
document.getElementById('sortby').addEventListener('change', sortStores);