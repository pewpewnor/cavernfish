/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{svelte,ts,js}'],
  theme: {
    extend: {
      colors: {
        cave: {
          base: '#0d0f17',
          deep: '#111420',
          surface: '#191c2a',
          raised: '#1e2236',
          elevated: '#252b40',
        },
        wire: {
          DEFAULT: '#2c3048',
          hi: '#3d4462',
        },
        ink: {
          DEFAULT: '#e2e8f5',
          mid: '#9aa5be',
          muted: '#636f8f',
          code: '#a5c8ff',
        },
        accent: {
          DEFAULT: '#3b6ef0',
          hi: '#5080f8',
          glow: 'rgba(59,110,240,0.18)',
          dim: 'rgba(59,110,240,0.12)',
        },
        ok: '#22d3a0',
        warn: '#f59e0b',
        err: '#f43f5e',
        violet: '#a78bfa',
        // HTTP method colours
        'method-get': '#22d3a0',
        'method-post': '#3b82f6',
        'method-put': '#f59e0b',
        'method-patch': '#a78bfa',
        'method-delete': '#f43f5e',
        'method-head': '#64748b',
        'method-options': '#64748b',
      },
      fontFamily: {
        mono: ['JetBrains Mono', 'Fira Code', 'Cascadia Code', 'ui-monospace', 'monospace'],
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
      borderRadius: {
        sm: '4px',
        DEFAULT: '6px',
        lg: '10px',
      },
      boxShadow: {
        modal: '0 24px 64px rgba(0,0,0,0.6), 0 0 0 1px rgba(59,110,240,0.08)',
        bp: '0 0 0 1px rgba(59,110,240,0.18), 0 32px 80px rgba(0,0,0,0.7)',
      },
      keyframes: {
        pulse: {
          '0%,100%': { opacity: '1', transform: 'scale(1)' },
          '50%': { opacity: '0.5', transform: 'scale(0.75)' },
        },
        breathe: {
          '0%,100%': { transform: 'scale(0.9)', opacity: '0.6' },
          '50%': { transform: 'scale(1.1)', opacity: '1' },
        },
      },
      animation: {
        'pulse-slow': 'pulse 1.2s ease-in-out infinite',
        breathe: 'breathe 3s ease-in-out infinite',
      },
    },
  },
  plugins: [],
}
