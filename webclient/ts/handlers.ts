import {DOMState} from "./classes/DOMStateManager";
import {clearFormFields, SeverityLevel, showMessage} from "./DOM";
import {client, domstate, tablemanager} from "./globals";
import User, {Role} from "./classes/User";
import {ErrorState} from "./classes/Client";
import {AdminTableManager} from "./classes/AdminTableManager";

const changeToUserOrAdmin = async (user: User): Promise<void> => {
    document.getElementById("username_display").innerText = `User: ${user.username}`;

    if (user.role == Role.Admin) {
        domstate.change(DOMState.Admin);

        // eslint-disable-next-line @typescript-eslint/ban-ts-ignore
        // @ts-ignore
        tablemanager = new AdminTableManager("admin_table");
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
        return;
    }

    const username = (document.getElementById("username") as HTMLInputElement).value;
    const password = (document.getElementById("password") as HTMLInputElement).value;

    const [res, err] = await client.login(username, password);
    if(err == ErrorState.InvalidCredentials) {
        showMessage("Invalid Credentials", SeverityLevel.ERROR);
    } else if(err == ErrorState.ServerError){
        showMessage("Server Error. please retry.", SeverityLevel.ERROR);
    } else if(err == ErrorState.Other){
        showMessage("An unknown error occurred. please retry.", SeverityLevel.ERROR);
    } else {
        clearFormFields();
        await changeToUserOrAdmin(res);
    }

};

export const signup = async (): Promise<void> => {
    if (domstate.state == DOMState.Login) {
        domstate.change(DOMState.Signup);
        return;
    }

    const username = (document.getElementById("username") as HTMLInputElement).value;
    const password = (document.getElementById("password") as HTMLInputElement).value;
    const email = (document.getElementById("email") as HTMLInputElement).value;

    const [res, err] = await client.signup(username, password, email);
    if(err === ErrorState.InvalidCredentials) {
        showMessage("Signup unsuccessful", SeverityLevel.ERROR);
    } else if(err === ErrorState.UserExists) {
        showMessage("Username exists", SeverityLevel.ERROR);
    } else if(err === ErrorState.ServerError) {
        showMessage("Server Error. please retry.", SeverityLevel.ERROR);
    } else if(err === ErrorState.Other) {
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
        showMessage("Please log in again", SeverityLevel.ERROR);
        domstate.change(DOMState.Login);
    } else if (err === ErrorState.ServerError) {
        showMessage("Server Error. please retry.", SeverityLevel.ERROR);
    } else if (err === ErrorState.Other) {
        showMessage("An unknown error occurred. please retry.", SeverityLevel.ERROR);
    } else {
        await changeToUserOrAdmin(res);
    }
};

export const changePasswordInit = (): void => {
    domstate.change(DOMState.ChangePassword);
};