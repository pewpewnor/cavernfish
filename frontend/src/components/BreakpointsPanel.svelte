<script lang="ts">
    import type { MockResponse, Endpoint, Folder, Collection } from "../lib/types";
    import { statusColor } from "../lib/types";
    import { collections, pendingBreakpoints, patchCollection } from "../stores/index";
    import * as api from "../lib/api";
    import ResponseEditor from "./ResponseEditor.svelte";
    import JsonEditor from "./JsonEditor.svelte";
    import {
        Button,
        CheckboxTree,
        ListItem,
        MethodBadge,
        SectionHeader,
        TabStrip,
    } from "../lib/ui";

    let panelTab = $state<"endpoints" | "queue">("endpoints");

    async function toggleEpBp(colId: string, folderId: string, ep: Endpoint) {
        const updated = await api.updateEndpoint(colId, folderId, ep.id, {
            ...ep,
            breakpoint: !ep.breakpoint,
        });
        patchCollection(updated);
    }

    async function setFolderBp(col: Collection, folder: Folder, on: boolean) {
        for (const ep of folder.endpoints) {
            if (ep.breakpoint !== on) {
                const updated = await api.updateEndpoint(col.id, folder.id, ep.id, {
                    ...ep,
                    breakpoint: on,
                });
                patchCollection(updated);
            }
        }
    }

    async function setCollectionBp(col: Collection, on: boolean) {
        for (const folder of col.folders) {
            await setFolderBp(col, folder, on);
        }
    }

    let activeQueueIdx = $state(0);
    let activeHit = $derived($pendingBreakpoints[activeQueueIdx] ?? null);

    let editingResponse = $state<MockResponse | null>(null);
    let prevHitId = $state("");
    let respTab = $state<"body" | "headers" | "cookies">("body");
    let reqTab = $state<"body" | "headers">("body");

    $effect(() => {
        if (activeHit?.id !== prevHitId) {
            prevHitId = activeHit?.id ?? "";
            editingResponse = activeHit ? JSON.parse(JSON.stringify(activeHit.response)) : null;
            respTab = "body";
            reqTab = "body";
        }
    });

    async function release() {
        if (!activeHit || !editingResponse) return;
        const id = activeHit.id;
        await api.releaseBreakpoint(id, $state.snapshot(editingResponse) as MockResponse);
        pendingBreakpoints.update((list) => list.filter((b) => b.id !== id));
        activeQueueIdx = Math.min(activeQueueIdx, $pendingBreakpoints.length - 1);
    }

    async function discard() {
        if (!activeHit) return;
        const id = activeHit.id;
        await api.discardBreakpoint(id);
        pendingBreakpoints.update((list) => list.filter((b) => b.id !== id));
        activeQueueIdx = Math.min(activeQueueIdx, $pendingBreakpoints.length - 1);
    }

    function formatHeaders(h: Record<string, string>): string {
        return Object.entries(h)
            .map(([k, v]) => `${k}: ${v}`)
            .join("\n");
    }

    let totalBpEndpoints = $derived(
        $collections
            .flatMap((c) => c.folders.flatMap((f) => f.endpoints))
            .filter((e) => e.breakpoint).length,
    );
</script>

<div class="flex min-w-0 flex-1 flex-col overflow-hidden">
    <div class="flex shrink-0 items-center gap-3 border-b border-wire bg-cave-surface px-5 py-3">
        <span class="text-base text-accent">⚡</span>
        <h2 class="font-semibold text-ink">Breakpoints</h2>
        {#if totalBpEndpoints > 0}
            <span class="rounded-full bg-accent/10 px-2 py-0.5 text-xs text-accent"
                >{totalBpEndpoints} active</span
            >
        {/if}
        {#if $pendingBreakpoints.length > 0}
            <span class="rounded-full bg-warn/10 px-2 py-0.5 text-xs text-warn"
                >{$pendingBreakpoints.length} pending</span
            >
        {/if}
    </div>

    <TabStrip
        bind:value={panelTab}
        tabs={[
            { value: "endpoints", label: "Endpoints" },
            { value: "queue", label: "Queue", badge: $pendingBreakpoints.length || null },
        ]}
    />

    {#if panelTab === "endpoints"}
        <div class="flex-1 overflow-y-auto p-4">
            <p class="mb-3 text-xs text-ink-muted">
                Enable breakpoint mode on specific endpoints. Matched requests will pause and appear
                in the Queue.
            </p>
            {#if $collections.length === 0}
                <p class="py-6 text-center text-sm text-ink-muted">
                    No collections. Create one in the Collections panel.
                </p>
            {:else}
                <CheckboxTree
                    collections={$collections}
                    isEndpointChecked={(ep) => ep.breakpoint}
                    onEndpointToggle={(ep, col, folder) => toggleEpBp(col.id, folder.id, ep)}
                    onFolderToggle={(col, folder, state) =>
                        setFolderBp(col, folder, state === "none")}
                    onCollectionToggle={(col, state) => setCollectionBp(col, state === "none")}
                />
            {/if}
        </div>
    {:else if $pendingBreakpoints.length === 0}
        <div class="flex flex-1 flex-col items-center justify-center gap-3 text-center">
            <span class="text-4xl opacity-20">⚡</span>
            <p class="text-sm text-ink-mid">No pending breakpoints</p>
            <p class="max-w-xs text-xs text-ink-muted">
                Enable breakpoints on endpoints and trigger requests to capture them here.
            </p>
        </div>
    {:else}
        <div class="flex min-h-0 flex-1 overflow-hidden">
            <aside class="flex w-52 shrink-0 flex-col border-r border-wire bg-cave-deep">
                <SectionHeader title="Queue ({$pendingBreakpoints.length})" />
                <div class="flex-1 overflow-y-auto py-1">
                    {#each $pendingBreakpoints as hit, i (hit.id)}
                        <ListItem
                            active={activeQueueIdx === i}
                            onclick={() => (activeQueueIdx = i)}
                        >
                            {#snippet leading()}
                                <span
                                    class="h-1.5 w-1.5 shrink-0 animate-pulse-slow rounded-full bg-accent"
                                ></span>
                            {/snippet}
                            <div class="flex items-center gap-1.5">
                                <MethodBadge method={hit.method} size="sm" />
                                <span class="truncate font-mono text-xs">{hit.path}</span>
                            </div>
                            <span class="text-[0.6875rem] text-ink-muted"
                                >{new Date(hit.timestamp).toLocaleTimeString()}</span
                            >
                        </ListItem>
                    {/each}
                </div>
            </aside>

            {#if activeHit}
                <div class="flex min-w-0 flex-1 flex-col overflow-hidden">
                    <div class="grid flex-1 grid-cols-2 overflow-hidden">
                        <div class="flex flex-col overflow-hidden border-r border-wire">
                            <div
                                class="flex shrink-0 flex-col gap-1 border-b border-wire bg-cave-surface px-4 py-2.5"
                            >
                                <span class="form-label">Incoming Request</span>
                                <div class="flex items-center gap-2">
                                    <MethodBadge method={activeHit.method} />
                                    <span class="truncate font-mono text-xs text-ink">
                                        {activeHit.path}{activeHit.query
                                            ? "?" + activeHit.query
                                            : ""}
                                    </span>
                                </div>
                            </div>

                            <TabStrip
                                bind:value={reqTab}
                                tabs={[
                                    { value: "body", label: "Body" },
                                    { value: "headers", label: "Headers" },
                                ]}
                            />

                            <div class="flex flex-1 flex-col overflow-y-auto p-3">
                                {#if reqTab === "body"}
                                    {#if activeHit.reqBody}
                                        <JsonEditor value={activeHit.reqBody} readOnly={true} />
                                    {:else}
                                        <div class="py-8 text-center text-sm text-ink-muted">
                                            No request body
                                        </div>
                                    {/if}
                                {:else}
                                    <pre
                                        class="whitespace-pre-wrap break-all font-mono text-xs leading-relaxed text-ink-mid">{formatHeaders(
                                            activeHit.reqHeaders,
                                        )}</pre>
                                {/if}
                            </div>
                        </div>

                        <div class="flex flex-col overflow-hidden">
                            <div
                                class="flex shrink-0 items-center gap-3 border-b border-wire bg-cave-surface px-4 py-2.5"
                            >
                                <span class="form-label">Response (editable)</span>
                                {#if editingResponse}
                                    <span
                                        class="font-mono text-lg font-bold"
                                        style="color:{statusColor(editingResponse.statusCode)}"
                                        >{editingResponse.statusCode}</span
                                    >
                                {/if}
                            </div>

                            <TabStrip
                                bind:value={respTab}
                                tabs={[
                                    { value: "body", label: "Body" },
                                    { value: "headers", label: "Headers" },
                                    { value: "cookies", label: "Cookies" },
                                ]}
                            />

                            {#if editingResponse}
                                <div class="flex min-h-0 flex-1 flex-col overflow-hidden">
                                    <ResponseEditor
                                        response={editingResponse}
                                        activeTab={respTab}
                                        showStatusBar={false}
                                        onChange={(r) => (editingResponse = r)}
                                    />
                                </div>
                            {/if}
                        </div>
                    </div>

                    <div
                        class="flex shrink-0 items-center justify-between border-t border-wire bg-cave-surface px-4 py-3"
                    >
                        <div class="flex items-center gap-2">
                            <span class="font-mono text-xs text-ink-muted">
                                {activeHit.endpoint.name || activeHit.endpoint.path}
                            </span>
                        </div>
                        <div class="flex gap-2">
                            <Button variant="danger" onclick={discard}>⊘ Discard</Button>
                            <Button variant="primary" onclick={release}>
                                ▶ Release with {editingResponse?.statusCode}
                            </Button>
                        </div>
                    </div>
                </div>
            {/if}
        </div>
    {/if}
</div>
