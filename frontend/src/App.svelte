<script lang="ts">
    import { onMount } from "svelte";
    import NavRail from "./components/NavRail.svelte";
    import WorkspacePage from "./components/WorkspacePage.svelte";
    import SettingsPage from "./components/SettingsPage.svelte";
    import { activePage, collections, servers, pendingBreakpoints } from "./stores/index";
    import * as api from "./lib/api";

    const ZOOM_KEY = "cavernfish:zoom";

    function loadZoom(): number {
        const raw = localStorage.getItem(ZOOM_KEY);
        if (!raw) return 1;
        const n = parseFloat(raw);
        return Number.isFinite(n) && n >= 0.5 && n <= 2 ? n : 1;
    }

    let zoom = $state(loadZoom());
    let bootState = $state<"loading" | "ready" | "error">("loading");
    let bootError = $state<string>("");
    let mockMode = $state(false);

    function handleKeydown(e: KeyboardEvent) {
        if (!e.ctrlKey) return;
        if (e.key === "=" || e.key === "+") {
            e.preventDefault();
            zoom = Math.min(+(zoom + 0.1).toFixed(1), 2);
        } else if (e.key === "-") {
            e.preventDefault();
            zoom = Math.max(+(zoom - 0.1).toFixed(1), 0.5);
        } else if (e.key === "0") {
            e.preventDefault();
            zoom = 1;
        }
    }

    $effect(() => {
        document.documentElement.style.removeProperty("zoom");
        if (zoom === 1) {
            document.documentElement.style.removeProperty("font-size");
        } else {
            document.documentElement.style.fontSize = `${16 * zoom}px`;
        }
        localStorage.setItem(ZOOM_KEY, String(zoom));
    });

    onMount(async () => {
        await api.waitForWails();
        mockMode = api.isMockMode();
        try {
            const [cols, srvs] = await Promise.all([api.getCollections(), api.getServers()]);
            collections.set(cols ?? []);
            servers.set(srvs ?? []);
            api.onServerStatus((updated) => servers.set(updated));
            api.onBreakpointHit((hit) => {
                pendingBreakpoints.update((list) => {
                    if (list.some((b) => b.id === hit.id)) return list;
                    return [...list, hit];
                });
            });
            api.onBreakpointReleased((id) => {
                pendingBreakpoints.update((list) => list.filter((b) => b.id !== id));
            });
            bootState = "ready";
        } catch (e) {
            bootError = e instanceof Error ? e.message : String(e);
            bootState = "error";
        }
    });
</script>

<svelte:window onkeydown={handleKeydown} />

{#if bootState === "ready"}
    <div class="flex h-full flex-col overflow-hidden">
        {#if mockMode}
            <div
                class="flex shrink-0 items-center justify-center gap-2 border-b border-warn/30 bg-warn/10 px-3 py-1 text-[0.6875rem] text-warn"
            >
                <span>⚠</span>
                <span
                    >Browser dev mode — data is stored in localStorage and HTTP servers are
                    inactive. Use <code class="font-mono">wails dev</code> for the full experience.</span
                >
            </div>
        {/if}
        <div class="flex flex-1 overflow-hidden">
            <NavRail />
            <div class="flex min-w-0 flex-1 overflow-hidden">
                {#if $activePage === "workspace"}
                    <WorkspacePage />
                {:else if $activePage === "settings"}
                    <SettingsPage />
                {/if}
            </div>
        </div>
    </div>
{:else if bootState === "loading"}
    <div class="flex h-full items-center justify-center bg-cave-base">
        <div class="text-sm text-ink-muted">Initializing…</div>
    </div>
{:else}
    <div class="flex h-full items-center justify-center bg-cave-base p-8">
        <div class="max-w-lg rounded-lg border border-err/40 bg-cave-surface p-6 text-center">
            <h2 class="mb-2 text-lg font-semibold text-err">Failed to start</h2>
            <p class="text-sm text-ink-mid">{bootError}</p>
        </div>
    </div>
{/if}
