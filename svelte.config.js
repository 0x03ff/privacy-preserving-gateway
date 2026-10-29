import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
    // 1. Preprocessors (like Tailwind/TypeScript support)
    preprocess: vitePreprocess(),

    // 2. Optional: Force runes mode globally across the app
    compilerOptions: {
        runes: true
    },

    // 3. SvelteKit specific options (like static adapter)
    kit: {
        adapter: adapter({
            pages: 'build',
            assets: 'build',
            fallback: 'index.html',
            precompress: false,
            strict: true
        })
    }
};

export default config;