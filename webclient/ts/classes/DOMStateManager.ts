// manual SPA
export enum DOMState {
    Login = "login",
    Admin = "admin",
    User = "user",
    Signup = "signup",
    ChangePassword = "changepassword",
}

export default class DOMStateManager {
    private currentState: DOMState;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    private static readonly stateStrings: Array<string> = Object.keys(DOMState).map(k => (DOMState as any)[k as any]);

    constructor(startstate: DOMState) {
        this.currentState = startstate;
        this.change(startstate);
    }

    private static hideAll(): void {
        DOMStateManager.stateStrings.forEach(key => {
            for (const i of Array.from(document.getElementsByClassName(key))) {
                (i as HTMLElement).style.display = "none";
            }
        });
    }

    change(state: DOMState): void {
        DOMStateManager.hideAll();

        // Show all matching
        for (const i of Array.from(document.getElementsByClassName(state))) {
            (i as HTMLElement).style.display = "inherit";
        }

        this.currentState = state;
    }

    get state(): DOMState {
        return this.currentState;
    }
}