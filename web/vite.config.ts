import path from "path";
import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
    plugins: [
        react({
            babel: {
                plugins: [["babel-plugin-react-compiler"]],
            },
        }),
        tailwindcss(),
    ],
    server: {
        proxy: {
            "/api": "http://127.0.0.1:3000",
        },
        watch: {
            usePolling: true,
        },
        host: '0.0.0.0',
        port: 5173,
        hmr: {
            host: 'localhost',
            protocol: 'ws',
            clientPort: 5173,
        }
    },
    resolve: {
        alias: {
            "@": path.resolve(__dirname, "./src"),
        },
    },
});
