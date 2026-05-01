# Cavernfish

Cavernfish is a cross-platform desktop & web app built using the Wails
framework (Go & Svelte TypeScript). It is an interactive web server manager that
you can use to mock endpoints in real-time.

## Color Theme

Dark blue-black palette. Update `tailwind.config.js` when changing.

| Token | Hex | Use |
|---|---|---|
| `cave-base` | `#1a1f30` | App background |
| `cave-deep` | `#222740` | Sidebar, nav rail |
| `cave-surface` | `#2c3250` | Panel surfaces, cards |
| `cave-raised` | `#363d5f` | Inputs, raised elements |
| `cave-elevated` | `#424a6e` | Modals, dropdowns |
| `wire` | `#525a82` | Borders |
| `wire-hi` | `#6671a0` | Highlighted borders, focuses |
| `ink` | `#f3f6fc` | Primary text |
| `ink-mid` | `#c9d2e8` | Secondary / label text |
| `ink-muted` | `#969fbe` | Placeholder, disabled text |
| `accent` | `#5483ff` | Primary accent (blue) |
| `accent-hi` | `#729cff` | Accent hover state |
| `ok` | `#22d3a0` | Success / GET method |
| `warn` | `#f59e0b` | Warning / PUT method / unsaved |
| `err` | `#f43f5e` | Error / DELETE method |

## Code Style

- No code comments, no section divider code comments
- No backwards compatibility, it's okay to rewrite or refactor existing features
- TypeScript strict mode
