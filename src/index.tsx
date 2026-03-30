import { definePlugin } from "@kleff/plugin-sdk";
import { ExamplePage } from "./pages/ExamplePage";
import { ExampleWidget } from "./components/ExampleWidget";

export default definePlugin({
  // Must match the plugin ID in your Kleff registry manifest
  id: "my-plugin",

  // Full pages — accessible at /p/my-plugin/example
  pages: [{ path: "/p/my-plugin/example", component: ExamplePage }],

  // Sidebar navigation items
  navItems: [
    {
      label: "My Plugin",
      icon: "Puzzle",
      path: "/p/my-plugin/example",
    },
  ],

  // Inject components into named slots in existing panel pages
  slots: {
    "dashboard.widgets": ExampleWidget,
  },
});
