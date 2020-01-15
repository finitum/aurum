import DOMStateManager, {DOMState} from "./classes/DOMStateManager";
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import {domstate} from "./globals";
import {changePassword, changePasswordInit, checkLogin, login, logout, onPasswordFieldChange, signup} from "./handlers";

export function setProperty(element: string, property: string, value: any): void {
    const el = document.getElementById(element);

    if (el === null || el === undefined) {
        console.warn(`Couldn't assign ${value} to ${element}.${property}`);
        return;
    }
    // @ts-ignore
    el[property] = value;
}

window.onload = (): void => {
    // @ts-ignore this works :tm:
    domstate = new DOMStateManager(DOMState.Login);

    checkLogin().then();

    setProperty("login-button", "onclick", login);
    setProperty("signup-button", "onclick",signup);
    setProperty("logout-button", "onclick",logout);
    setProperty("changepassword-button", "onclick",changePassword);
    setProperty("changepassword-init-button", "onclick",changePasswordInit);
    setProperty("password", "onkeyup",onPasswordFieldChange);
};



