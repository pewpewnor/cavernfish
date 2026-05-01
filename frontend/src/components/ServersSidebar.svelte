<script lang="ts">
    import type { ServerConfig } from "../lib/types";
    import { servers, selectedServerId } from "../stores/index";
    import * as api from "../lib/api";
    import { Button, ListItem } from "../lib/ui";
    import SidebarShell from "./SidebarShell.svelte";

    async function createServer() {
        const cfg: ServerConfig = {
            id: crypto.randomUUID(),
            name: "New Server",
            port: 8000,
            collectionIds: [],
            https: false,
            disabledEndpointIds: [],
        };
        const created = await api.createServer(cfg);
        servers.update((list) => [...list, created]);
        selectedServerId.set(created.id);
    }
</script>

<SidebarShell title="Servers">
    {#snippet actions()}
        <Button variant="icon" onclick={createServer} title="New Server">+</Button>
    {/snippet}

    <div class="flex-1 overflow-y-auto py-1">
        {#each $servers as srv (srv.id)}
            <ListItem
                active={$selectedServerId === srv.id}
                onclick={() => selectedServerId.set(srv.id)}
            >
                {#snippet leading()}
                    <span
                        class="h-2 w-2 shrink-0 rounded-full
                {srv.status === 'running'
                            ? 'bg-ok'
                            : srv.status === 'error'
                              ? 'bg-err'
                              : 'bg-ink-muted'}"
                    ></span>
                {/snippet}
                <div class="truncate text-[0.8125rem] font-medium">{srv.name || "Unnamed"}</div>
                <div class="flex items-center gap-1.5 font-mono text-[0.6875rem] text-ink-muted">
                    <span>:{srv.port}</span>
                    {#if srv.https}<span class="text-accent">HTTPS</span>{/if}
                </div>
            </ListItem>
        {/each}

        {#if $servers.length === 0}
            <div class="flex flex-col items-center gap-2.5 px-4 py-10 text-center">
                <p class="text-xs text-ink-muted">No servers yet</p>
                <Button variant="primary" onclick={createServer}>+ New Server</Button>
            </div>
        {/if}
    </div>
</SidebarShell>
