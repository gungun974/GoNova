import { defineConfig } from "vite";
import { createServer } from "http";

import { resolve } from "path";

const r = (p: string) => resolve(__dirname, p);

export const alias: Record<string, string> = {
  "@": r("./resources"),
  "@/": r("./resources/"),
};

export default defineConfig({
  server: {
    origin: "http://localhost:5173",
    watch: {
      ignored: ["**/mails/**/*.html"],
    },
  },
  base: "/assets",
  build: {
    copyPublicDir: false,
    outDir: "build/public/assets",
    assetsDir: "",
    manifest: true,
    rollupOptions: {
      input: {
        main: "resources/main.ts",
      },
    },
  },
  plugins: [
    {
      name: "disable hot update templ",
      handleHotUpdate({ file }) {
        if (file.endsWith(".templ")) {
          return [];
        }
        if (file.endsWith(".go")) {
          return [];
        }
      },
    },

    {
      name: "manual-refresh",
      configureServer(server) {
        const refreshServer = createServer((req, res) => {
          if (req.url === "/trigger-refresh") {
            server.ws.send({ type: "full-reload" });
            res.end("Refresh triggered!");
          } else {
            res.end("Unknown request");
          }
        });

        refreshServer.listen(5174);
      },
    },
  ],
  resolve: {
    alias,
  },
  clearScreen: false,
});
