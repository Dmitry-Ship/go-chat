import { defineConfig } from "vite";
import reactRefresh from "@vitejs/plugin-react-refresh";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [reactRefresh()],
  server: {
    proxy: {
      "/ws": {
        target: "wss://api:8080",
        changeOrigin: true,
        ws: true,
        secure: false,
      },
    },
  },
});
