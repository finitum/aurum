import {DOMState} from "./classes/DOMStateManager";
import {clearFormFields, SeverityLevel, showMessage} from "./DOM";
import {client, domstate, tablemanager} from "./globals";
import User, {Role} from "./classes/User";
import {ErrorState} from "./classes/Client";
import {AdminTableManager} from "./classes/AdminTableManager";
import {verifyPassword} from "./passwords";
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
    const [user, err] = await client.checkLogin();

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

    const [res, err] = await client.login(username, password);
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

    for(const dot of Array.from(dots)) {
        dot.className = "password-input-dot";
    }

    switch (score) {
        case 4:
            dots[0].className = "password-input-dot password-input-dot-selected";
            /* fallthrough */
        case 3:
            dots[1].className = "password-input-dot password-input-dot-selected";
            /* fallthrough */
        case 2:
            dots[2].className = "password-input-dot password-input-dot-selected";
            /* fallthrough */
        case 1:
            dots[3].className = "password-input-dot password-input-dot-selected";
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

        const passwordStatus = await verifyPassword(password.value, [username.value, email.value]);

        setPasswordDots(passwordStatus.score);

        const el = document.getElementById("password-suggestion");
        if(el == null) {
            console.warn("password-suggestion element could not be found");
            return;
        }

        if (passwordStatus.score >= 2) {
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
};

export const signup = async (): Promise<void> => {
    if (domstate.state == DOMState.Login) {
        domstate.change(DOMState.Signup);
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
    if (passwordStatus.score < 2) {
        showMessage("Please enter a valid password", SeverityLevel.WARNING);
        return;
    }

    const [res, err] = await client.signup(username, password, email);
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

    client.logout();
    domstate.change(DOMState.Login);
};

export const changePassword = async (): Promise<void> => {
    const password = (document.getElementById("password") as HTMLInputElement).value;
    const [res, err] = await client.changePassword(password);

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