// manual SPA
export enum DOMState {
    Login = "login",
    Admin = "admin",
    User = "user",
    Signup = "signup",
    ChangePassword = "changepassword",
}

//TODO: also singleton?
// export const DOMStateManagerInstance = new DOMStateManager(DOMState.Login);
// All methods static?
export default class DOMStateManager {
    private currentState: DOMState;

    constructor(startstate: DOMState) {
        this.currentState = startstate;
        this.registerHistoryHandler();
        this.change(startstate);
    }

    private popStateHandler(event: PopStateEvent): void {
        // setTimeout of 0 to make sure DOM is loaded
        setTimeout(() => {
            if(Object.values(DOMState).includes(event.state)) {
                this.change(event.state, false);
            }
        }, 0);
    }

    private registerHistoryHandler(): void {
        // Wrap in a lambda to preserve `this`
        window.addEventListener("popstate", (event: PopStateEvent) => {
            this.popStateHandler(event);
        });
    }

    private pushStateHandler(state: DOMState = this.currentState): void {
        history.pushState(state, "");
    }

    private static hideAll(): void {
        Object.values(DOMState).forEach(key => {
            for (const i of Array.from(document.getElementsByClassName(key))) {
                (i as HTMLElement).classList.add("hidden");
            }
        });
    }

    change(state: DOMState, setHistory= true): void {
        DOMStateManager.hideAll();

        // Show all matching
        for (const i of Array.from(document.getElementsByClassName(state))) {
            (i as HTMLElement).classList.remove("hidden");
        }

        this.currentState = state;

        if (setHistory) {
            this.pushStateHandler(state);
        }
    }

    back(): void {
        history.back();
    }

    get state(): DOMState {
        return this.currentState;
    }
}

