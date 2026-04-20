import { defineConfig } from "tsup";

export default defineConfig({
  entry: { index: "src/index.tsx" },
  format: ["esm"],
  dts: false,
  splitting: false,
  clean: true,
  external: ["react", "react-dom", "@kleffio/ui", "@kleffio/plugin-sdk"],
});
