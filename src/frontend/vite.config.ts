import { ConfigEnv, defineConfig, loadEnv } from "vite";

import { resolve } from "node:path";
import vue from "@vitejs/plugin-vue";
import { URLSearchParams } from "node:url";

// https://vitejs.dev/config/
export default ({ mode }: ConfigEnv) => {
    const env = loadEnv(mode, ".");
    const devInitParams = new URLSearchParams({
        applicationUuid: env.VITE_APP_UUID || "cndTest",
        userSign: env.VITE_USER_SIGN || "1000005",
        sessionId: "00000"
    });
    // const devInitRewrite = `/api/init?${devInitParams.toString()}`;
    return defineConfig({
        plugins: [ vue() ],
        server: {
            proxy: {
                "/api": {
                    target: env.VITE_PROXY_TARGET || "http://localhost:8080",
                    changeOrigin: true
                    // configure: (proxy) => {
                    //     proxy.on("proxyReq", (proxyReq, req) => {
                    //         const origin = `${proxyReq.protocol}//${req.headers["host"]}`;
                    //         proxyReq.setHeader("Origin", origin);
                    //     });
                    // },
                    // rewrite: (path) => {
                    //     const urlpath = new URL(path);
                    //     devInitParams.forEach((value, name) => {
                    //         urlpath.searchParams.append(name, value);
                    //     });
                    //     return urlpath.toString();
                    // }
                }
            },
            cors: true
            // host: env.VITE_HOST || "0.0.0.0"
            // port: parseInt(env.VITE_PORT) || 3000
        },
        build: {
            target: "esnext",
            outDir: "dist",
            manifest: false,
            minify: "esbuild",
            emptyOutDir: true,
            sourcemap: false,
            cssCodeSplit: false
        },
        resolve: {
            alias:
                {
                    "@": resolve(import.meta.dirname, "src")
                }
        }
    });
};
