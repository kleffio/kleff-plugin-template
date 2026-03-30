import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

/**
 * Kleff plugin build config.
 *
 * Builds to a single IIFE bundle (dist/plugin.js).
 * React and all @kleff/* packages are externalized — they are provided
 * by the panel at runtime via window.__kleff.
 */
export default defineConfig({
  plugins: [react()],
  build: {
    lib: {
      entry: "src/index.tsx",
      formats: ["iife"],
      name: "__kleff_plugin_bundle",
      fileName: () => "plugin.js",
    },
    rollupOptions: {
      external: ["react", "react-dom", "@kleff/ui", "@kleff/plugin-sdk"],
      output: {
        globals: {
          react: "window.__kleff.React",
          "react-dom": "window.__kleff.ReactDOM",
          "@kleff/ui": "window.__kleff.ui",
          "@kleff/plugin-sdk": "window.__kleff.sdk",
        },
      },
    },
    outDir: "dist",
    emptyOutDir: true,
  },
});
