# Cavernfish

Cavernfish is a cross-platform desktop & web app built using the Wails
framework (Go & Svelte TypeScript). It is an interactive web server manager that
you can use to mock endpoints in real-time.

## Color Theme

Dark blue-black palette. Update `tailwind.config.js` when changing.

| Token | Hex | Use |
|---|---|---|
| `cave-base` | `#0d0f17` | App background |
| `cave-deep` | `#111420` | Sidebar, nav rail |
| `cave-surface` | `#191c2a` | Panel surfaces, cards |
| `cave-raised` | `#1e2236` | Inputs, raised elements |
| `cave-elevated` | `#252b40` | Modals, dropdowns |
| `wire` | `#2c3048` | Borders |
| `wire-hi` | `#3d4462` | Highlighted borders, focuses |
| `ink` | `#e2e8f5` | Primary text |
| `ink-mid` | `#9aa5be` | Secondary / label text |
| `ink-muted` | `#636f8f` | Placeholder, disabled text |
| `accent` | `#3b6ef0` | Primary accent (blue) |
| `accent-hi` | `#5080f8` | Accent hover state |
| `ok` | `#22d3a0` | Success / GET method |
| `warn` | `#f59e0b` | Warning / PUT method / unsaved |
| `err` | `#f43f5e` | Error / DELETE method |

## Code Style

- No code comments, no section divider code comments
- No backwards compatibility, it's okay to rewrite or refactor existing features
- TypeScript strict mode
