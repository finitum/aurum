import {Claims, TokenPair} from "../API/jwt";

// ServiceWorker's postMessage
// eslint-disable-next-line @typescript-eslint/no-explicit-any
declare function postMessage(message: any, transferList?: Transferable[]): void;

// Worker code
let localToken: TokenPair | null = null;

// Receiving tokens from the client
onmessage = (event: MessageEvent): void => {
    localToken = event.data as TokenPair;
};

// Checking token expiry
setInterval(() => {

    console.log("hit");

    if(localToken != null && localToken.loginToken != ""  && localToken.refreshToken != ""){
        const claims = Claims.parse(localToken.loginToken);

        const date = new Date();
        date.setMinutes(date.getMinutes() + 5);
        if (claims.expiresAt < date) {
            postMessage(null);
            localToken = null;
        }
    }
}, 1000 * 10);