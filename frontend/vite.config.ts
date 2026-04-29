import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [
    svelte({
      onwarn(warning, handler) {
        // Suppress a11y warnings — intentional UX decisions in this codebase
        if (warning.code.startsWith('a11y-')) return
        handler(warning)
      },
    }),
  ],
})
