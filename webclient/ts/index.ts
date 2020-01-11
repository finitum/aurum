import DOMStateManager, {DOMState} from "./classes/DOMStateManager";
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import {domstate} from "./globals";
import {changePassword, changePasswordInit, checkLogin, login, logout, signup} from "./handlers";



window.onload = (): void => {
    // @ts-ignore this works :tm:
    domstate = new DOMStateManager(DOMState.Login);

    checkLogin().then();

    document.getElementById("login_button").onclick = login;
    document.getElementById("signup_button").onclick = signup;
    document.getElementById("logout_button").onclick = logout;
    document.getElementById("changepassword_button").onclick = changePassword;
    document.getElementById("changepassword_init_button").onclick = changePasswordInit;
};


