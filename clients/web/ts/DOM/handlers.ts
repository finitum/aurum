import {DOMState} from "./DOMStateManager";
import {clearFormFields, SeverityLevel, showMessage} from "./DOMFunctions";
import {domstate, tablemanager} from "../globals";
import User, {Role} from "../API/User";
import Client, {ErrorState} from "../API/Client";
import {AdminTableManager} from "./AdminTableManager";
import {verifyPassword} from "../API/passwords";
import zxcvbn from "zxcvbn";

const changeToUserOrAdmin = async (user: User): Promise<void> => {
    const uEl = document.getElementById("username-display");

    if(uEl == null) {
        console.warn("Could not find username-display");
        return;
    }

    uEl.innerText = `User: ${user.username}`;

    if (user.role == Role.Admin) {
        domstate.change(DOMState.Admin);

        // eslint-disable-next-line @typescript-eslint/ban-ts-ignore
        // @ts-ignore
        tablemanager = new AdminTableManager("admin-table");
        const err = await tablemanager.fill();

        if (err != ErrorState.Ok) {
            showMessage("Couldn't get user info", SeverityLevel.ERROR);
        }

    } else {
        domstate.change(DOMState.User);
    }
};

export const checkLogin = async (): Promise<boolean> => {
    const [user, err] = await Client.getInstance().checkLogin();

    if (user != null && err === ErrorState.Ok) {
        await changeToUserOrAdmin(user);
        return true;
    }
    return false;
};


export const login = async (): Promise<void> => {
    if (domstate.state == DOMState.Signup) {
        domstate.change(DOMState.Login);
        const el = document.getElementById("login-button");
        if(el) el.innerText = "Login";
        else console.warn("Login button not found");
        return;
    }

    const username = (document.getElementById("username") as HTMLInputElement).value;
    const password = (document.getElementById("password") as HTMLInputElement).value;

    if (username.length == 0) {
        showMessage("Please provide a Username", SeverityLevel.WARNING);
        return;
    }

    if (username.length == 0) {
        showMessage("Please provide a Password", SeverityLevel.WARNING);
        return;
    }

    const [res, err] = await Client.getInstance().login(username, password);
    if(err == ErrorState.InvalidCredentials) {
        showMessage("Invalid Credentials", SeverityLevel.WARNING);
    } else if(err == ErrorState.ServerError){
        showMessage("Server Error. please retry.", SeverityLevel.ERROR);
    } else if(err == ErrorState.Other || res == null){
        showMessage("An unknown error occurred. please retry.", SeverityLevel.ERROR);
    } else {
        clearFormFields();
        await changeToUserOrAdmin(res);
    }
};

const setPasswordDots = (score: zxcvbn.ZXCVBNScore): void => {
    const dots = document.getElementsByClassName("password-input-dot");
    console.assert(dots.length == 4, "Assert that there are 4 password dots");

    Array.from(dots).forEach(dot => {
        dot.classList.remove("password-input-dot-selected");
    });

    switch (score) {
        case 4:
            dots[0].classList.add("password-input-dot-selected");
            /* fallthrough */
        case 3:
            dots[1].classList.add("password-input-dot-selected");
            /* fallthrough */
        case 2:
            dots[2].classList.add("password-input-dot-selected");
            /* fallthrough */
        case 1:
            dots[3].classList.add("password-input-dot-selected");
            /* fallthrough */
        case 0:
            break;
    }
};

export const onPasswordFieldChange = async (): Promise<void> => {
    if (domstate.state == DOMState.Signup || domstate.state == DOMState.ChangePassword) {
        const username = (document.getElementById("username") as HTMLInputElement);
        const password = (document.getElementById("password") as HTMLInputElement);
        const email = (document.getElementById("email") as HTMLInputElement);

        const el = document.getElementById("password-suggestion");
        if (el == null) {
            console.warn("password-suggestion element could not be found");
            return;
        }

        const passwordStatus = await verifyPassword(password.value, [username.value, email.value]);
        if (typeof passwordStatus == "string") {
            setPasswordDots(0);
            el.innerText = passwordStatus;
            el.style.visibility = "visible";
        } else {
            setPasswordDots(passwordStatus.score);

            if (passwordStatus.score > 2) {
                el.style.visibility = "hidden";
            } else {
                el.style.visibility = "visible";
                if (passwordStatus.feedback.warning == "") {
                    el.innerText = "Password strength too low.";
                } else {
                    el.innerText = passwordStatus.feedback.suggestions.join(" ");
                }
            }
        }
    }
};

export const signup = async (): Promise<void> => {
    if (domstate.state == DOMState.Login) {
        domstate.change(DOMState.Signup);
        onPasswordFieldChange();

        const button = document.getElementById("login-button");
        if (button == null) return;
        button.innerText = "Back";
        return;
    }

    const username = (document.getElementById("username") as HTMLInputElement).value;
    const password = (document.getElementById("password") as HTMLInputElement).value;
    const email = (document.getElementById("email") as HTMLInputElement).value;

    if (username.length === 0) {
        showMessage("Please enter a username", SeverityLevel.WARNING);
        return;
    }

    if (email.length === 0) {
        showMessage("Please enter an email address", SeverityLevel.WARNING);
        return;
    }

    const passwordStatus = await verifyPassword(password, [username, email]);
    if(typeof passwordStatus === "string") {
        showMessage(passwordStatus, SeverityLevel.WARNING);
    } else if (passwordStatus.score < 2) {
        showMessage("Please enter a valid password", SeverityLevel.WARNING);
        return;
    }

    const [res, err] = await Client.getInstance().signup(username, password, email);
    if(err === ErrorState.InvalidCredentials) {
        showMessage("Signup unsuccessful", SeverityLevel.WARNING);
    } else if(err === ErrorState.InvalidPasswordError) {
        showMessage("Password invalid", SeverityLevel.WARNING);
    }else if(err === ErrorState.UserExists) {
        showMessage("Username exists", SeverityLevel.WARNING);
    } else if(err === ErrorState.ServerError) {
        showMessage("Server Error. please retry.", SeverityLevel.ERROR);
    } else if(err === ErrorState.Other || res == null) {
        showMessage("An unknown error occurred. please retry.", SeverityLevel.ERROR);
    } else {
        clearFormFields();
        await changeToUserOrAdmin(res);
    }

};

export const logout = (): void => {
    clearFormFields();

    Client.getInstance().logout();
    domstate.change(DOMState.Login);
};

export const changePassword = async (): Promise<void> => {
    const password = (document.getElementById("password") as HTMLInputElement).value;
    const [res, err] = await Client.getInstance().changePassword(password);

    if (err === ErrorState.InvalidCredentials) {
        showMessage("Please log in again", SeverityLevel.WARNING);
        domstate.change(DOMState.Login);
    } else if (err === ErrorState.ServerError) {
        showMessage("Server Error. please retry.", SeverityLevel.ERROR);
    } else if (err === ErrorState.Other || res == null) {
        showMessage("An unknown error occurred. please retry.", SeverityLevel.ERROR);
    } else {
        await changeToUserOrAdmin(res);
    }
};

export const changePasswordInit = (): void => {
    domstate.change(DOMState.ChangePassword);
};