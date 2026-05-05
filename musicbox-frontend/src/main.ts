import { createApp } from "vue";
import { createPinia } from "pinia";
import router from "./router";
// 全局样式
import "./assets/styles/reset.css";
// UnoCSS
import "uno.css";
import App from "./App.vue";

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);
app.mount("#app");
