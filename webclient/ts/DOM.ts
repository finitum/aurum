export enum SeverityLevel {
    INFO,
    WARNING,
    ERROR
}

/**
 *
 * @param message
 * @param level
 */
export const showMessage = (message: string, level: SeverityLevel): void => {
    const messageElement = document.getElementById("message");
 
    messageElement.innerText = message;

    switch (level) {
        case SeverityLevel.INFO:
            messageElement.style.background = "lightblue";
            break;
        case SeverityLevel.WARNING:
            messageElement.style.background = "orange";
            break;
        case SeverityLevel.ERROR:
            messageElement.style.background = "red";
            break;
    }

    messageElement.style.opacity = "1";
    messageElement.style.visibility = "visible";

    messageElement.onclick = (): void => {
        messageElement.style.opacity = "0";
        setTimeout(() => {
                messageElement.style.visibility = "hidden";
        },1000);
    };
};

export const clearFormFields = (): void => {
    (document.getElementById("username") as HTMLInputElement).value = "";
    (document.getElementById("password") as HTMLInputElement).value = "";
    (document.getElementById("email") as HTMLInputElement).value = "";
};