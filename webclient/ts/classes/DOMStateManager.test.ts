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

        expect(document.getElementById("div1").style.display).toBe("inherit");
        expect(document.getElementById("div2").style.display).toBe("none");
        expect(document.getElementById("div3").style.display).toBe("none");
    });

    it("Should be able to change to a new state", () => {
        // Setup the DOM
        document.body.innerHTML = `
            <div id="div1" class="login"></div>
            <div id="div2" class="admin"></div>
            <div id="div3" class="user"></div>`;

        // The actual call
        const m = new DOMStateManager(DOMState.Login);

        expect(document.getElementById("div1").style.display).toBe("inherit");
        expect(document.getElementById("div2").style.display).toBe("none");
        expect(document.getElementById("div3").style.display).toBe("none");

        m.change(DOMState.User);

        expect(m.state).toBe(DOMState.User);
        expect(document.getElementById("div1").style.display).toBe("none");
        expect(document.getElementById("div2").style.display).toBe("none");
        expect(document.getElementById("div3").style.display).toBe("inherit");
    });

    it("Should be able to handle multi classed divs", () => {
        // Setup the DOM
        document.body.innerHTML = `
            <div id="div1" class="login admin"></div>
            <div id="div2" class="admin user"></div>
            <div id="div3" class="user login"></div>`;

        // The actual call
        const m = new DOMStateManager(DOMState.Admin);

        expect(document.getElementById("div1").style.display).toBe("inherit");
        expect(document.getElementById("div2").style.display).toBe("inherit");
        expect(document.getElementById("div3").style.display).toBe("none");

        m.change(DOMState.Login);

        expect(document.getElementById("div1").style.display).toBe("inherit");
        expect(document.getElementById("div2").style.display).toBe("none");
        expect(document.getElementById("div3").style.display).toBe("inherit");

        m.change(DOMState.User);

        expect(document.getElementById("div1").style.display).toBe("none");
        expect(document.getElementById("div2").style.display).toBe("inherit");
        expect(document.getElementById("div3").style.display).toBe("inherit");
    });
});