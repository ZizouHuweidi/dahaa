import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";
import { VitePWA } from "vite-plugin-pwa";

export default defineConfig({
	plugins: [
		tailwindcss(),
		reactRouter(),
		VitePWA({
			outDir: "build/client",
			registerType: "autoUpdate",
			includeAssets: ["favicon.ico", "apple-touch-icon.svg"],
			manifest: {
				name: "Dahaa",
				short_name: "Dahaa",
				description: "A fast invite-link trivia party game for the web.",
				theme_color: "#111827",
				background_color: "#fff7ed",
				display: "standalone",
				orientation: "portrait-primary",
				start_url: "/",
				scope: "/",
				icons: [
					{
						src: "/pwa-icon.svg",
						sizes: "any",
						type: "image/svg+xml",
						purpose: "any maskable",
					},
				],
			},
			workbox: {
				navigateFallback: "/offline.html",
				runtimeCaching: [
					{
						urlPattern: ({ url }) => url.pathname.startsWith("/api/"),
						handler: "NetworkOnly",
						method: "GET",
					},
				],
			},
		}),
	],
	resolve: {
		tsconfigPaths: true,
	},
});
