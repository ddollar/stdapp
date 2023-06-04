import { DefaultApolloClient } from "@vue/apollo-composable";
import { apolloClient } from "@/lib/apollo";

import { createApp, provide, h } from "vue";
import App from "./App.vue";
const app = createApp({
	setup() {
		provide(DefaultApolloClient, apolloClient);
	},
	render: () => h(App),
});

import "bootstrap";

import { createPinia } from "pinia";
app.use(createPinia());

import router from "./router";
app.use(router);

import timeago from "vue-timeago3";
app.use(timeago, {
	converterOptions: {
		includeSeconds: false,
		useStrict: true,
	},
});

import vSelect from "vue-select";
app.component("v-select", vSelect);

app.mount("#app");
