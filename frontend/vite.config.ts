import { defineConfig } from "vite";
import reactRefresh from "@vitejs/plugin-react-refresh";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [reactRefresh()],
  server: {
    proxy: {
      // string shorthand
      // "/ws": "ws://localhost:8080/ws",
      // with options
      "/ws": {
        target: "ws://localhost:8080",
        // changeOrigin: true,
        ws: true,
      },
    },
  },
});
