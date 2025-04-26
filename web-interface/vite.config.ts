import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import svgr from 'vite-plugin-svgr';

export default defineConfig({
  plugins: [react(), svgr()],
  server: {
    port: 3001,
    allowedHosts: ["api.homeshare.pro"],
    strictPort: true,
    host: true
  },
});
