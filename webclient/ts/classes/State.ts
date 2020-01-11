import {TokenPair} from "./jwt";


class Storage {
    public tokenPair: TokenPair;

    public json(): string {
        return JSON.stringify(this);
    }

    public static fromJSON(json: string): Storage {
        return JSON.parse(json);
    }
}

export default class State {
    private static readonly storage_key: string = "storage";
    private readonly stored: Storage;

    constructor() {
        const local = localStorage.getItem(State.storage_key);
        if(local){
            this.stored = Storage.fromJSON(local);
        } else {
            this.stored = new Storage();
        }
    }

    private store(): void {
        localStorage.setItem(State.storage_key, JSON.stringify(this.stored));
    }

    get tokenPair(): TokenPair {
        if (this.stored.tokenPair === null) {
            return null;
        }
        return new TokenPair(this.stored.tokenPair.loginToken, this.stored.tokenPair.refreshToken);
    }

    set tokenPair(tokenPair: TokenPair) {
        this.stored.tokenPair = tokenPair;
        this.store();
    }
}
