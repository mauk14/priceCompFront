function changeImage(src) {
    document.getElementById('main-product-image').src = src;
}

function sortStores() {
    var storesList = document.querySelector('.stores-list');
    var stores = document.querySelectorAll('.store');

    var sortBy = document.getElementById('sortby').value;
    var sortedStores = Array.from(stores);

    sortedStores.sort(function(a, b) {
        var priceA = parseInt(a.getAttribute('data-price'));
        var priceB = parseInt(b.getAttribute('data-price'));
        var ratingA = parseInt(a.getAttribute('data-rating'));
        var ratingB = parseInt(b.getAttribute('data-rating'));
        var recommendedA = parseInt(a.getAttribute('data-recommended'));
        var recommendedB = parseInt(b.getAttribute('data-recommended'));

        if (sortBy === 'cheapest') {
            return priceA - priceB;
        } else if (sortBy === 'rating') {
            return ratingB - ratingA; // Note that we do b - a to get descending order
        } else if (sortBy === 'recommended') {
            return recommendedB - recommendedA; // Highest recommended first
        }
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

document.getElementById('applyFilters').addEventListener('click', function() {
    // Collect filter values
    // Note: This is a simplified example. You would collect all filter values here.
    const manufacturer = document.querySelector('input[name="manufacturer"]:checked').value;
    const headphoneType = document.querySelector('input[name="headphoneType"]:checked').value;

    // TODO: Send these values to the backend or use them to filter the products displayed
    console.log('Applying filters:', manufacturer, headphoneType);

    // You would likely have a function like this:
    // applyFilters(manufacturer, headphoneType);
});






