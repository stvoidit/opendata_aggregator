import App from "@/App.vue";
import { createApp } from "vue";
import router from "@/router";
const app = createApp(App);
app.use(router);
fetch("/api/init").then(() => {
    app.mount("#app");
}).catch(alert);
