import {setProperty} from "./index";

export enum SeverityLevel {
    INFO,
    WARNING,
    ERROR
}

/**
 * Shows a toast notification a message.
 * @param message The message to display
 * @param level The Severity of the message, this determines the background color of the message.
 */
export const showMessage = (message: string, level: SeverityLevel): void => {
    const messageElement = document.getElementById("message");
    if (messageElement ==  null) return;
    messageElement.innerText = message;

    switch (level) {
        case SeverityLevel.INFO:
            messageElement.style.background = "#1e92f4";
            break;
        case SeverityLevel.WARNING:
            messageElement.style.background = "#ffc107";
            break;
        case SeverityLevel.ERROR:
            messageElement.style.background = "#cc3300";
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
    setProperty("username", "value", "");
    setProperty("password", "value", "");
    setProperty("email", "value", "");
};