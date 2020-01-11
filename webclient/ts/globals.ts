import Client from "./classes/Client";
import config from "./config";
import DOMStateManager from "./classes/DOMStateManager";
import {AdminTableManager} from "./classes/AdminTableManager";

export const client = new Client(config.API_URL);
export let domstate: DOMStateManager;
export let tablemanager: AdminTableManager;
