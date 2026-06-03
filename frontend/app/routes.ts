import { index, type RouteConfig, route } from "@react-router/dev/routes";

export default [
	index("routes/home.tsx"),
	route("create", "routes/create.tsx"),
	route("g/:code", "routes/game.$code.tsx"),
] satisfies RouteConfig;
