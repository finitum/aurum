import DOMStateManager, {DOMState} from "./DOMStateManager";

describe("DOMStateManager", () => {
    it("Should show initial state and hide all others", () => {
        // Setup the DOM
        document.body.innerHTML = `
            <div id="div1" class="login"></div>
            <div id="div2" class="admin"></div>
            <div id="div3" class="user"></div>`;

        // The actual call
        new DOMStateManager(DOMState.Login);

        expect(document.getElementById("div1").classList.contains("hidden")).toBeFalsy();
        expect(document.getElementById("div2").classList.contains("hidden")).toBeTruthy();
        expect(document.getElementById("div3").classList.contains("hidden")).toBeTruthy();
    });

    it("Should be able to change to a new state", () => {
        // Setup the DOM
        document.body.innerHTML = `
            <div id="div1" class="login"></div>
            <div id="div2" class="admin"></div>
            <div id="div3" class="user"></div>`;

        // The actual call
        const m = new DOMStateManager(DOMState.Login);


        expect(document.getElementById("div1").classList.contains("hidden")).toBeFalsy();
        expect(document.getElementById("div2").classList.contains("hidden")).toBeTruthy();
        expect(document.getElementById("div3").classList.contains("hidden")).toBeTruthy();

        m.change(DOMState.User);

        expect(m.state).toBe(DOMState.User);
        expect(document.getElementById("div1").classList.contains("hidden")).toBeTruthy();
        expect(document.getElementById("div2").classList.contains("hidden")).toBeTruthy();
        expect(document.getElementById("div3").classList.contains("hidden")).toBeFalsy();
    });

    it("Should be able to handle multi classed divs", () => {
        // Setup the DOM
        document.body.innerHTML = `
            <div id="div1" class="login admin"></div>
            <div id="div2" class="admin user"></div>
            <div id="div3" class="user login"></div>`;

        // The actual call
        const m = new DOMStateManager(DOMState.Login);

        expect(document.getElementById("div1").classList.contains("hidden")).toBeFalsy();
        expect(document.getElementById("div2").classList.contains("hidden")).toBeTruthy();
        expect(document.getElementById("div3").classList.contains("hidden")).toBeFalsy();

        m.change(DOMState.Login);

        expect(document.getElementById("div1").classList.contains("hidden")).toBeFalsy();
        expect(document.getElementById("div2").classList.contains("hidden")).toBeTruthy();
        expect(document.getElementById("div3").classList.contains("hidden")).toBeFalsy();

        m.change(DOMState.User);

    });
});