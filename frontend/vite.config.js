import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

export default defineConfig({
    plugins: [react()],
    css: {
        // semantic-ui-css 2.5.0 contains a technically invalid selector
        // ([data-tooltip]:after .header) that lightningcss rejects without this
        lightningcss: { errorRecovery: true }
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
