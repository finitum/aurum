import { SingleUserClient } from "aurum-client";
import router from "@/router";

// TODO: Change baseurl
export const client = new SingleUserClient("http://localhost:8042");

client.AddUnauthorizedHandler(async () => {
    await router.push("/login")
})