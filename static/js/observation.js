var canvas = document.getElementById('canvas');
var context = canvas.getContext('2d');
var video = document.getElementById('video');
var button = document.getElementById("snap");
var videoRunning = true;

if (navigator.mediaDevices.getUserMedia) {
    navigator.mediaDevices.getUserMedia({video: true})
        .then(function (stream) {
            video.srcObject = stream;
        })
        .catch(function (err) {
            console.log("Something went wrong!");
        });
}

// Trigger photo take
button.addEventListener("click", function () {
    if (videoRunning) {
        video.style.display = "none";
        canvas.style.display = "block";
        context.drawImage(video, 0, 0, 640, 480);

        this.innerText = "Reset";

        // canvas to img
        document.getElementById('image').value = canvas.toDataURL();

    } else {
        video.style.display = "block";
        canvas.style.display = "none";
        context.drawImage(video, 0, 0, 640, 480);

        this.innerText = "Foto nemen";
    }
    // reminder for switching from video to snap and vice versa
    videoRunning = !videoRunning;
});