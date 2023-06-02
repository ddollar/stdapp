import { fileURLToPath, URL } from "url";
import { defineConfig } from "vite";

// plugins
import vue from "@vitejs/plugin-vue";
import basicSsl from "@vitejs/plugin-basic-ssl";

export default defineConfig({
	base: process.env.VITE_PREFIX,
	plugins: [basicSsl(), vue()],
	resolve: {
		alias: {
			"@": fileURLToPath(new URL("./src", import.meta.url)),
		},
	},
	server: {
		hmr: {
			clientPort: process.env.VITE_CLIENT_PORT,
		},
		https: true,
		port: process.env.PORT,
	},
});
