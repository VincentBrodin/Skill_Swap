const fileInput = document.getElementById("fileInput");
const fileList = document.getElementById("fileList");
const form = document.getElementById("form");

const files = [];
let id = 0;


fileInput.addEventListener("change", UpdateFileInput);

function PrettySize(bytes) {
    const units = [
        "B",
        "KB",
        "MB",
        "GB",
        "TB",
        "PB",
        "EB",
        "ZB",
        "YB",
    ];
    const exponent = Math.min(
        Math.floor(Math.log(bytes) / Math.log(1000)),
        units.length - 1,
    );
    const approx = bytes / 1000 ** exponent;
    const output =
        exponent === 0
            ? `${bytes} bytes`
            : `${approx.toFixed(0)} ${units[exponent]
            } (${bytes} bytes)`;
    return output

}

function UpdateFileInput(_) {
    fileList.innerHTML = ""

    if (fileInput.files.length == 0) {
        return
    }

    console.log(fileInput.files)
    for (let i = 0; i < fileInput.files.length; i++) {
        const file = fileInput.files[i];
        const li = document.createElement("li");
        fileList.appendChild(li);
        li.className = 'list-group-item d-flex justify-content-between align-items-center';
        li.innerHTML = `
                <span>${file.name}</span>
                <span class="badge bg-secondary">${PrettySize(file.size)} KB</span>
            `;
    }
}
