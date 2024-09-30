function TogglePassword(inputId, iconId) {
    const input = document.getElementById(inputId);
    const icon = document.getElementById(iconId);

    if (input.type === "password") {
        input.type = "text";
        icon.classList.remove("fa-eye");
        icon.classList.add("fa-eye-slash");
    } else {
        input.type = "password";
        icon.classList.add("fa-eye");
        icon.classList.remove("fa-eye-slash");
    }
}

function MatchPassword(inputIdA, inputIdB, promptId) {
    const inputA = document.getElementById(inputIdA);
    const inputB = document.getElementById(inputIdB);
    const prompt = document.getElementById(promptId);

    if (inputA.value != inputB.value && inputB.value != "") {
        prompt.innerText = "Password must match";
    } else {
        prompt.innerText = "";
    }
}
