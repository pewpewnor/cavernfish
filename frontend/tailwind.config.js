/** @type {import('tailwindcss').Config} */
export default {
    content: ["./index.html", "./src/**/*.{svelte,ts,js}"],
    theme: {
        extend: {
            colors: {
                cave: {
                    base: "#1a1f30",
                    deep: "#222740",
                    surface: "#2c3250",
                    raised: "#363d5f",
                    elevated: "#424a6e",
                },
                wire: {
                    DEFAULT: "#525a82",
                    hi: "#6671a0",
                },
                ink: {
                    DEFAULT: "#f3f6fc",
                    mid: "#c9d2e8",
                    muted: "#969fbe",
                    code: "#bcd0ff",
                },
                accent: {
                    DEFAULT: "#5483ff",
                    hi: "#729cff",
                    glow: "rgba(84,131,255,0.22)",
                    dim: "rgba(84,131,255,0.14)",
                },
                ok: "#22d3a0",
                warn: "#f59e0b",
                err: "#f43f5e",
                violet: "#a78bfa",
                // HTTP method colours
                "method-get": "#22d3a0",
                "method-post": "#3b82f6",
                "method-put": "#f59e0b",
                "method-patch": "#a78bfa",
                "method-delete": "#f43f5e",
                "method-head": "#64748b",
                "method-options": "#64748b",
            },
            fontFamily: {
                mono: ["JetBrains Mono", "Fira Code", "Cascadia Code", "ui-monospace", "monospace"],
                sans: ["Inter", "system-ui", "-apple-system", "sans-serif"],
            },
            borderRadius: {
                sm: "4px",
                DEFAULT: "6px",
                lg: "10px",
            },
            boxShadow: {
                modal: "0 24px 64px rgba(0,0,0,0.6), 0 0 0 1px rgba(59,110,240,0.08)",
                bp: "0 0 0 1px rgba(59,110,240,0.18), 0 32px 80px rgba(0,0,0,0.7)",
            },
            keyframes: {
                pulse: {
                    "0%,100%": { opacity: "1", transform: "scale(1)" },
                    "50%": { opacity: "0.5", transform: "scale(0.75)" },
                },
                breathe: {
                    "0%,100%": { transform: "scale(0.9)", opacity: "0.6" },
                    "50%": { transform: "scale(1.1)", opacity: "1" },
                },
            },
            animation: {
                "pulse-slow": "pulse 1.2s ease-in-out infinite",
                breathe: "breathe 3s ease-in-out infinite",
            },
        },
    },
    plugins: [],
};
