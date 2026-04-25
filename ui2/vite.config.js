import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
    plugins: [react({ include: /\.(jsx|js)$/ })],
    esbuild: {
        loader: 'jsx',
        include: /src\/.*\.js$/,
        exclude: []
    },
    optimizeDeps: {
        esbuildOptions: {
            loader: { '.js': 'jsx' }
        }
    },
    build: {
        outDir: 'build'
    },
    server: {
        proxy: {
            '/RunAnimations': 'http://192.168.0.68:8080'
        }
    }
});
