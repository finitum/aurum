import {client, tablemanager} from "../globals";
import {ErrorState} from "./Client";
import {Role} from "./User";


export class AdminTableManager {
    private tableElement: HTMLTableElement;

    constructor(id: string) {
        this.tableElement = document.getElementById(id) as HTMLTableElement;
    }

    private addHeader(): void {
        const head = this.tableElement.createTHead();
        const row = head.insertRow(0);

        row.appendChild(document.createElement("th")).innerText = "Username";
        row.appendChild(document.createElement("th")).innerText = "Email";
        row.appendChild(document.createElement("th")).innerText = "Role";
        row.appendChild(document.createElement("th")).innerText = "Blocked";
    }

    private addRow(username: string, email: string, role: Role, blocked: boolean): void {
        const body = this.tableElement.createTBody();
        const row = body.appendChild(document.createElement("tr")) as HTMLTableRowElement;

        row.appendChild(document.createElement("td")).innerText = username;
        row.appendChild(document.createElement("td")).innerText = email;
        row.appendChild(document.createElement("td")).innerText = Role[role];
        row.appendChild(document.createElement("td")).innerText = String(blocked);
    }


    clear(): void {
        while (this.tableElement.firstChild) {
            this.tableElement.removeChild(this.tableElement.firstChild);
        }

        this.addHeader();
    }

    async fill(): Promise<ErrorState> {
        tablemanager.clear();
        const [users, err] = await client.getUsers(0, 100);
        if (err !== ErrorState.Ok) {
            return err;
        }

        for (const user of users) {
            this.addRow(user.username, user.email, user.role, user.blocked);
        }

        return ErrorState.Ok;
    }
}