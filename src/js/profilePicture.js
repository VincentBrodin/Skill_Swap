function UpdateProfilePicture(inputId, imageId) {
    const input = document.getElementById(inputId);
    const image = document.getElementById(imageId);

    if (input.files && input.files[0]) {
        const reader = new FileReader();

        reader.onload = function (e) {
            image.src = e.target.result;
        };

        reader.readAsDataURL(input.files[0]);
    }
}
