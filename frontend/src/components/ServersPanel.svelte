<script lang="ts">
    import { onMount } from "svelte";
    import type { ServerInfo, ServerConfig } from "../lib/types";
    import { statusColor } from "../lib/types";
    import { servers, collections, requestLogs, addLog, selectedServerId } from "../stores/index";
    import * as api from "../lib/api";
    import { Button, CheckboxTree, MethodBadge, Modal, TabStrip } from "../lib/ui";

    let editForm = $state<ServerConfig | null>(null);
    let dirty = $state(false);
    let detailTab = $state<"serve" | "logs">("serve");
    let logSearch = $state("");

    let restartConfirmOpen = $state(false);
    let pendingSaveCfg = $state<ServerConfig | null>(null);

    let selectedServer = $derived($servers.find((s) => s.id === $selectedServerId) ?? null);

    let filteredLogs = $derived(
        ($requestLogs[$selectedServerId ?? ""] ?? []).filter(
            (e) =>
                !logSearch ||
                e.path.toLowerCase().includes(logSearch.toLowerCase()) ||
                e.method.toLowerCase().includes(logSearch.toLowerCase()) ||
                String(e.statusCode).includes(logSearch),
        ),
    );

    let prevSelectedId = $state<string | null>(null);

    $effect(() => {
        const id = $selectedServerId;
        if (id === prevSelectedId) return;
        prevSelectedId = id;
        const srv = $servers.find((s) => s.id === id);
        editForm = srv
            ? {
                  id: srv.id,
                  name: srv.name,
                  port: srv.port,
                  collectionIds: [...(srv.collectionIds ?? [])],
                  https: srv.https,
                  disabledEndpointIds: [...(srv.disabledEndpointIds ?? [])],
              }
            : null;
        dirty = false;
        detailTab = "serve";
        logSearch = "";
        if (id) {
            const sid = id;
            api.getRequestLog(sid).then((log) => {
                requestLogs.update((r) => ({ ...r, [sid]: log }));
            });
        }
    });

    function markDirty() {
        dirty = true;
    }

    async function persistConfig(cfg: ServerConfig) {
        const updated = await api.updateServer(cfg);
        servers.update((list) => list.map((s) => (s.id === updated.id ? updated : s)));
        dirty = false;
    }

    async function saveServer() {
        if (!editForm) return;
        if (!editForm.port || editForm.port < 1 || editForm.port > 65535) {
            alert("Invalid port");
            return;
        }
        if (selectedServer?.status === "running") {
            pendingSaveCfg = $state.snapshot(editForm) as ServerConfig;
            restartConfirmOpen = true;
            return;
        }
        await persistConfig(editForm);
    }

    async function saveWithoutRestart() {
        if (!pendingSaveCfg) return;
        await persistConfig(pendingSaveCfg);
        pendingSaveCfg = null;
        restartConfirmOpen = false;
    }

    async function saveAndRestart() {
        if (!pendingSaveCfg) return;
        await persistConfig(pendingSaveCfg);
        await api.restartServer(pendingSaveCfg.id);
        pendingSaveCfg = null;
        restartConfirmOpen = false;
    }

    function cancelSave() {
        pendingSaveCfg = null;
        restartConfirmOpen = false;
    }

    async function toggleServer(s: ServerInfo) {
        if (s.status === "running") await api.stopServer(s.id);
        else await api.startServer(s.id);
    }

    async function deleteServer(id: string) {
        if (!confirm("Delete this server?")) return;
        await api.deleteServer(id);
        servers.update((list) => list.filter((s) => s.id !== id));
        if ($selectedServerId === id) selectedServerId.set(null);
    }

    function isCollectionEnabled(colId: string): boolean {
        return (editForm?.collectionIds ?? []).includes(colId);
    }

    function isEndpointChecked(epId: string, colId: string): boolean {
        if (!isCollectionEnabled(colId)) return false;
        return !(editForm?.disabledEndpointIds ?? []).includes(epId);
    }

    function toggleCollectionInForm(colId: string) {
        if (!editForm) return;
        const ids = editForm.collectionIds;
        if (ids.includes(colId)) {
            const col = $collections.find((c) => c.id === colId);
            const colEpIds = col ? col.folders.flatMap((f) => f.endpoints.map((e) => e.id)) : [];
            editForm = {
                ...editForm,
                collectionIds: ids.filter((id) => id !== colId),
                disabledEndpointIds: editForm.disabledEndpointIds.filter(
                    (id) => !colEpIds.includes(id),
                ),
            };
        } else {
            editForm = { ...editForm, collectionIds: [...ids, colId] };
        }
        markDirty();
    }

    function toggleEndpoint(epId: string, colId: string) {
        if (!editForm) return;
        if (!isCollectionEnabled(colId)) {
            editForm = { ...editForm, collectionIds: [...editForm.collectionIds, colId] };
        }
        const disabled = editForm.disabledEndpointIds;
        if (disabled.includes(epId)) {
            editForm = { ...editForm, disabledEndpointIds: disabled.filter((id) => id !== epId) };
        } else {
            editForm = { ...editForm, disabledEndpointIds: [...disabled, epId] };
        }
        markDirty();
    }

    async function clearLog(id: string) {
        await api.clearRequestLog(id);
        requestLogs.update((r) => ({ ...r, [id]: [] }));
    }

    function formatTime(ts: string): string {
        return new Date(ts).toLocaleTimeString();
    }

    onMount(() => {
        api.onRequestLogged((entry) => addLog(entry));
    });
</script>

<div class="flex min-w-0 flex-1 flex-col overflow-hidden">
    {#if editForm && selectedServer}
        <div class="flex shrink-0 flex-col gap-3 border-b border-wire bg-cave-surface px-5 py-4">
            <div class="flex items-center gap-3">
                <input
                    class="flex-1 text-base font-semibold"
                    bind:value={editForm.name}
                    placeholder="Server name"
                    oninput={markDirty}
                />
                <div class="flex shrink-0 items-center gap-2">
                    {#if dirty}
                        <span class="h-2.5 w-2.5 rounded-full bg-warn" title="Unsaved changes"
                        ></span>
                    {/if}
                    <Button onclick={saveServer} disabled={!dirty}>
                        {dirty ? "Save" : "Saved ✓"}
                    </Button>
                    <Button
                        variant={selectedServer.status === "running" ? "ghost" : "primary"}
                        onclick={() => selectedServer && toggleServer(selectedServer)}
                        disabled={selectedServer.status === "error"}
                    >
                        {selectedServer.status === "running" ? "⏹ Stop" : "▶ Start"}
                    </Button>
                    <Button
                        variant="icon"
                        onclick={() => selectedServer && deleteServer(selectedServer.id)}>🗑</Button
                    >
                </div>
            </div>

            <div class="flex flex-wrap items-end gap-5">
                <div class="flex flex-col gap-1">
                    <span class="form-label">Port</span>
                    <input
                        type="number"
                        min="1"
                        max="65535"
                        class="w-24 font-mono"
                        bind:value={editForm.port}
                        oninput={markDirty}
                    />
                </div>
                <label class="flex cursor-pointer items-center gap-2 pb-1">
                    <input type="checkbox" bind:checked={editForm.https} onchange={markDirty} />
                    <span class="text-sm text-ink-mid">HTTPS</span>
                    {#if editForm.https}
                        <span class="font-mono text-xs text-ink-muted"
                            >(HTTP also on :{editForm.port + 1})</span
                        >
                    {/if}
                </label>
                <div class="flex items-center gap-2 pb-1">
                    <span class="status-pill status-{selectedServer.status}"
                        >{selectedServer.status}</span
                    >
                    {#if selectedServer.requestCount > 0}
                        <span class="font-mono text-xs text-ink-muted"
                            >{selectedServer.requestCount} req</span
                        >
                    {/if}
                    {#if selectedServer.errorMsg}
                        <span
                            class="max-w-[240px] truncate text-xs text-err"
                            title={selectedServer.errorMsg}>⚠ {selectedServer.errorMsg}</span
                        >
                    {/if}
                </div>
            </div>
        </div>

        <TabStrip
            bind:value={detailTab}
            tabs={[
                { value: "serve", label: "Serve" },
                {
                    value: "logs",
                    label: "Logs",
                    badge: selectedServer.requestCount || null,
                },
            ]}
            onChange={(v) => {
                if (v === "logs" && $selectedServerId) {
                    const sid = $selectedServerId;
                    api.getRequestLog(sid).then((log) => {
                        requestLogs.update((r) => ({ ...r, [sid]: log }));
                    });
                }
            }}
        />

        <div class="flex min-h-0 flex-1 flex-col overflow-hidden">
            {#if detailTab === "serve"}
                <div class="flex-1 overflow-y-auto p-3">
                    {#if $collections.length === 0}
                        <p class="py-6 text-center text-sm text-ink-muted">
                            No collections. Create one in the Collections panel.
                        </p>
                    {:else}
                        <CheckboxTree
                            collections={$collections}
                            isEndpointChecked={(ep, col) => isEndpointChecked(ep.id, col.id)}
                            onEndpointToggle={(ep, col) => toggleEndpoint(ep.id, col.id)}
                            onCollectionToggle={(col) => toggleCollectionInForm(col.id)}
                            onFolderToggle={(col, folder, state) => {
                                if (!editForm) return;
                                if (!editForm.collectionIds.includes(col.id)) {
                                    editForm = {
                                        ...editForm,
                                        collectionIds: [...editForm.collectionIds, col.id],
                                    };
                                }
                                const epIds = folder.endpoints.map((e) => e.id);
                                let disabled = [...editForm.disabledEndpointIds];
                                if (state === "none") {
                                    disabled = disabled.filter((id) => !epIds.includes(id));
                                } else {
                                    disabled = [...new Set([...disabled, ...epIds])];
                                }
                                editForm = { ...editForm, disabledEndpointIds: disabled };
                                markDirty();
                            }}
                        />
                    {/if}
                </div>
            {:else}
                <div class="flex shrink-0 items-center gap-2 border-b border-wire px-3 py-2">
                    <input
                        class="flex-1"
                        placeholder="Search by path, method, or status…"
                        bind:value={logSearch}
                    />
                    <Button size="sm" onclick={() => selectedServer && clearLog(selectedServer.id)}
                        >Clear</Button
                    >
                </div>
                <div class="flex-1 overflow-y-auto">
                    {#if filteredLogs.length === 0}
                        <div class="py-10 text-center text-sm text-ink-muted">
                            {logSearch ? "No matching entries" : "No requests yet"}
                        </div>
                    {:else}
                        {#each [...filteredLogs].reverse() as entry (entry.id)}
                            <div
                                class="flex items-center gap-2 border-b border-wire/40 px-3 py-2 text-xs
                {!entry.matched ? 'opacity-50' : ''}"
                            >
                                <span class="w-16 shrink-0 font-mono text-ink-muted"
                                    >{formatTime(entry.timestamp)}</span
                                >
                                <MethodBadge method={entry.method} size="sm" />
                                <span class="min-w-0 flex-1 truncate font-mono text-ink-mid"
                                    >{entry.path}{entry.query ? "?" + entry.query : ""}</span
                                >
                                <span
                                    class="shrink-0 font-mono font-bold"
                                    style="color:{statusColor(entry.statusCode)}"
                                    >{entry.statusCode}</span
                                >
                                <span class="w-12 shrink-0 text-right font-mono text-ink-muted"
                                    >{entry.latencyMs}ms</span
                                >
                                {#if entry.isProxy}
                                    <span
                                        class="shrink-0 rounded-[3px] bg-cave-elevated px-1 py-0.5 font-mono text-[0.5625rem] text-ink-muted"
                                        >proxy</span
                                    >
                                {/if}
                                {#if entry.isWs}
                                    <span
                                        class="shrink-0 rounded-[3px] bg-cave-elevated px-1 py-0.5 font-mono text-[0.5625rem] text-ink-muted"
                                        >ws</span
                                    >
                                {/if}
                                {#if entry.breakpointed}
                                    <span
                                        class="shrink-0 rounded-[3px] bg-accent/[.12] px-1 py-0.5 font-mono text-[0.5625rem] text-accent"
                                        >⚡bp</span
                                    >
                                {/if}
                            </div>
                        {/each}
                    {/if}
                </div>
            {/if}
        </div>
    {:else}
        <div class="flex flex-1 flex-col items-center justify-center gap-3 p-10 text-center">
            <div class="text-5xl opacity-[0.12]">◎</div>
            <p class="text-sm text-ink-mid">Select a server or create a new one</p>
        </div>
    {/if}
</div>

{#if restartConfirmOpen}
    <Modal open={restartConfirmOpen} title="Server is running" onClose={cancelSave}>
        <div class="flex flex-col gap-4">
            <p class="text-sm text-ink-mid">
                Collection and endpoint changes apply immediately. Port and HTTPS changes only take
                effect after a restart.
            </p>
            <div class="flex justify-end gap-2">
                <Button variant="primary" onclick={saveAndRestart}>Save & Restart</Button>
                <Button onclick={saveWithoutRestart}>Save without restart</Button>
                <Button onclick={cancelSave}>Cancel</Button>
            </div>
        </div>
    </Modal>
{/if}
