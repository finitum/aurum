import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";

// Only add one icon to the library
import { library } from "@fortawesome/fontawesome-svg-core";
import { faPen } from "@fortawesome/free-solid-svg-icons/faPen";
import { faTimes } from "@fortawesome/free-solid-svg-icons/faTimes";
import { faPlus } from "@fortawesome/free-solid-svg-icons/faPlus";
library.add(faPen);
library.add(faTimes);
library.add(faPlus);

// eslint-disable-next-line
createApp(App as any)
  .use(router)
  .mount("#app");
