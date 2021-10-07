import { defineConfig, loadEnv } from "vite";
import reactRefresh from "@vitejs/plugin-react-refresh";

// https://vitejs.dev/config/

export default ({ mode }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };
  return defineConfig({
    plugins: [reactRefresh()],
    clearScreen: false,
    server: {
      proxy: {
        "/ws": {
          target: process.env.VITE_WS_DOMAIN,
          changeOrigin: true,
          ws: true,
          secure: false,
        },
        "/api": {
          target: process.env.VITE_DOMAIN,
          changeOrigin: true,
          secure: false,
          rewrite: (path) => path.replace(/^\/api/, ""),
        },
      },
    },
  });
};
