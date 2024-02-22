const ratings = {
    5: 4,
    4: 1,
    3: 0,
    2: 0,
    1: 0
};

// Calculate the maximum number of ratings for any star level to set relative widths
const maxRatingCount = Math.max(...Object.values(ratings));

// Set the width of each rating bar relative to the maximum count
Object.keys(ratings).forEach(star => {
    const bar = document.querySelector(`.rating-bar-${star}`);
    const width = (ratings[star] / maxRatingCount) * 100;
    bar.style.width = `${width}%`;
});


function filterReviews(type) {
    // First, remove 'active' class from all buttons
    document.querySelectorAll('.review-filter-button').forEach(button => {
        button.classList.remove('active');
    });

    // Add 'active' class to the current selected button
    document.querySelector(`.review-filter-button[onclick="filterReviews('${type}')"]`).classList.add('active');

    // Hide all reviews
    document.querySelectorAll('.reviews-container .review').forEach(review => {
        review.style.display = 'none';
    });

    // Show only reviews that match the selected type
    if (type === 'all') {
        document.querySelectorAll('.reviews-container .review').forEach(review => {
            review.style.display = '';
        });
    } else {
        document.querySelectorAll(`.reviews-container .review[data-review-type="${type}"]`).forEach(review => {
            review.style.display = '';
        });
    }
}

// Initially show all reviews
filterReviews('all');

document.addEventListener('DOMContentLoaded', function() {
    // Assuming you have a reviewsData object/array with review and rating information
    const reviewsData = [
        { id: 1, rating: 4.5 },
        { id: 2, rating: 3 },
        // ... more reviews
    ];

    // Function to convert numerical rating to stars
    function getStars(rating) {
        // Round to nearest half
        rating = Math.round(rating * 2) / 2;
        let output = [];
        // Append all the full stars
        for (var i = rating; i >= 1; i--)
            output.push('★');
        // If there is a half a star, append it
        if (i == .5) output.push('☆');
        // Fill the empty stars
        for (let i = (5 - rating); i >= 1; i--)
            output.push('☆');
        return output.join('');
    }

    // Find each review and set its rating
    reviewsData.forEach(function(reviewData) {
        const reviewElement = document.querySelector(`#review-${reviewData.id}`);
        const starsElement = reviewElement.querySelector('.review-stars');
        starsElement.textContent = getStars(reviewData.rating);
    });
});

document.getElementById('write-review-form').addEventListener('submit', function(event) {
    event.preventDefault();

    // Get form values
    var name = document.getElementById('reviewer-name').value;
    var reviewText = document.getElementById('review-text').value;
    var rating = document.querySelector('input[name="rating"]:checked').value;
    var photos = document.getElementById('review-photos').files;

    // TODO: Handle the form submission, e.g., send to the server
    console.log('Review Submitted', { name, reviewText, rating, photos });

    // After handling submission, you can clear the form or give feedback to the user
    // this.reset(); // Uncomment to reset the form
    // alert('Thank you for your review!'); // Feedback to the user
});