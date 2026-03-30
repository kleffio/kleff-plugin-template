# Kleff Plugin Template

Starter template for building a Kleff plugin. Plugins run as Docker containers alongside Kleff and can contribute UI (pages, sidebar nav, dashboard widgets) and backend logic (custom API routes, middleware).

## Prerequisites

- Node 22+
- Go 1.23+
- Docker

## Quick start

```sh
# 1. Install JS dependencies
npm install

# 2. Build the JS bundle
npm run build        # → dist/plugin.js
# npm run dev        # watch mode

# 3. Run the gRPC server locally
cd server
go run .
# Plugin listening on :50051
```

Connect a running Kleff instance to it by pointing `PLUGIN_GRPC_ADDR` at `localhost:50051`, or build the Docker image and install through the marketplace.

## Build the Docker image

```sh
docker build -t my-plugin:latest .
docker run --rm -p 50051:50051 my-plugin:latest
```

## Rename the plugin

Find and replace `my-plugin` in these files:

| File | What to change |
|------|----------------|
| `src/index.tsx` | `id` field, route paths |
| `server/main.go` | `PLUGIN_ID` default, nav item path |
| `server/go.mod` | module name |
| `package.json` | `name` field |
| `kleff-plugin.json` | `id`, `name`, `image`, `repo`, `author` |

## Publish to the Kleff marketplace

1. Build and push your Docker image to a public registry:
   ```sh
   docker build -t ghcr.io/your-org/my-plugin:1.0.0 .
   docker push ghcr.io/your-org/my-plugin:1.0.0
   ```

2. Fill in `kleff-plugin.json` with your plugin's details.

3. Open a pull request to [kleffio/plugin-registry](https://github.com/kleffio/plugin-registry) adding your entry to `plugins.json`.

Once merged, your plugin appears in the Kleff marketplace for all self-hosters automatically — the platform fetches the catalog from GitHub on startup.

---

## Plugin capabilities

Declare what your plugin does in `server/main.go` → `GetCapabilities`:

| Capability | What it does |
|------------|-------------|
| `ui.manifest` | Contributes nav items, pages, and widgets to the panel |
| `identity.provider` | Acts as the Kleff authentication provider |
| `api.middleware` | Intercepts every authenticated API request (for RBAC, auditing, etc.) |
| `api.routes` | Owns one or more HTTP routes, handled directly by the plugin |

## UI slots

When you declare `ui.manifest`, the panel injects your components into named slots. Slots are declared in `src/index.tsx`:

```tsx
slots: {
  "dashboard.widgets": MyWidget,
}
```

| Slot ID | Where it renders |
|---------|-----------------|
| `dashboard.widgets` | Dashboard page — metric cards, charts, widgets |
| `dashboard.actions` | Dashboard header — action buttons |
| `sidebar.nav` | Sidebar navigation — links with icons |
| `sidebar.bottom` | Sidebar footer — extra links, status badges |
| `server.detail.tabs` | Server detail page — extra tabs |
| `server.detail.actions` | Server detail header — action buttons |
| `settings.tabs` | Settings page — settings sections |
| `admin.tabs` | Admin page — admin tool panels |
| `topbar.actions` | Top bar — notification icons, quick actions |

Full pages are registered separately and are accessible at `/p/{plugin-id}/...`:

```tsx
pages: [
  { path: "/p/my-plugin/dashboard", component: DashboardPage },
]
```

## Available UI components

All `@kleff/ui` components are available. They resolve to the panel's live component tree at runtime — nothing extra is bundled into your plugin.

```tsx
import {
  Button, Card, CardHeader, CardTitle, CardContent,
  Input, Label, Badge, Alert, Tabs, TabsList, TabsTrigger, TabsContent,
  Table, TableHeader, TableBody, TableRow, TableCell,
  Dialog, DialogContent, DialogHeader, DialogTitle,
  Select, SelectTrigger, SelectValue, SelectContent, SelectItem,
  Separator, Skeleton, Tooltip, TooltipContent, TooltipTrigger,
} from "@kleff/ui";
```

## Config fields in kleff-plugin.json

Config fields are rendered as form inputs in the Kleff install modal and injected as environment variables into your container.

| `type` | Input rendered | Notes |
|--------|---------------|-------|
| `string` | Text input | |
| `secret` | Password input | Stored encrypted at rest |
| `url` | URL input | |
| `number` | Number input | |
| `boolean` | Toggle | |
| `select` | Dropdown | Requires `options` array |

## Environment variables injected by Kleff

The platform always injects these into your container:

| Variable | Value |
|----------|-------|
| `PLUGIN_ID` | Your plugin ID (matches `id` in `kleff-plugin.json`) |
| `PLUGIN_PORT` | gRPC port to listen on (always `50051`) |

Plus any fields you declared in `config` in `kleff-plugin.json`.

## Project structure

```
src/
  index.tsx                   — plugin entry: pages, navItems, slots
  pages/ExamplePage.tsx       — full page at /p/my-plugin/example
  components/
    ExampleWidget.tsx         — widget injected into the dashboard.widgets slot

server/
  main.go                     — gRPC server: Health, GetCapabilities, GetUIManifest
  go.mod

kleff-plugin.json             — registry manifest (submit this to the marketplace)
vite.config.ts                — IIFE build; React and @kleff/* are externalized
Dockerfile                    — multi-stage: JS build → Go build → final image
```

---

## How it works (for the curious)

1. Self-hosters run `docker compose up` — no plugin-related repos needed.
2. The Kleff platform fetches the plugin catalog from `github.com/kleffio/plugin-registry` automatically.
3. When an admin installs a plugin, the platform pulls the Docker image, starts it on the `kleff` Docker network, and dials gRPC at `kleff-{plugin-id}:50051`.
4. The platform calls `GetCapabilities` → `GetUIManifest` to discover what the plugin contributes.
5. The panel fetches `GET /api/v1/plugins/ui-manifests` and renders nav items, pages, and widgets.
