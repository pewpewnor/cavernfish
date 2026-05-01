<script lang="ts">
    import type { Collection, Folder, Endpoint } from "../types";
    import { METHOD_COLORS } from "../types";

    type State = "all" | "some" | "none";

    interface Props {
        collections: Collection[];
        isEndpointChecked: (ep: Endpoint, col: Collection, folder: Folder) => boolean;
        onEndpointToggle: (ep: Endpoint, col: Collection, folder: Folder) => void;
        onFolderToggle?: (col: Collection, folder: Folder, current: State) => void;
        onCollectionToggle?: (col: Collection, current: State) => void;
        showCollectionCheckbox?: boolean;
    }

    let {
        collections,
        isEndpointChecked,
        onEndpointToggle,
        onFolderToggle,
        onCollectionToggle,
        showCollectionCheckbox = true,
    }: Props = $props();

    let expandedCols = $state(new Set<string>());
    let expandedFolders = $state(new Set<string>());

    function toggleCol(id: string) {
        const next = new Set(expandedCols);
        if (next.has(id)) next.delete(id);
        else next.add(id);
        expandedCols = next;
    }

    function toggleFolder(id: string) {
        const next = new Set(expandedFolders);
        if (next.has(id)) next.delete(id);
        else next.add(id);
        expandedFolders = next;
    }

    function folderState(col: Collection, folder: Folder): State {
        if (folder.endpoints.length === 0) return "none";
        const n = folder.endpoints.filter((e) => isEndpointChecked(e, col, folder)).length;
        if (n === 0) return "none";
        if (n === folder.endpoints.length) return "all";
        return "some";
    }

    function collectionState(col: Collection): State {
        const eps = col.folders.flatMap((f) => f.endpoints.map((e) => ({ e, f })));
        if (eps.length === 0) return "none";
        const n = eps.filter(({ e, f }) => isEndpointChecked(e, col, f)).length;
        if (n === 0) return "none";
        if (n === eps.length) return "all";
        return "some";
    }

    function mc(m: string): string {
        return METHOD_COLORS[m.toUpperCase()] ?? "#64748b";
    }
</script>

<div class="flex flex-col gap-2">
    {#each collections as col (col.id)}
        {@const cs = collectionState(col)}
        <div class="overflow-hidden rounded-lg border border-wire">
            <div class="flex items-center gap-2 bg-cave-surface px-3 py-2.5">
                <button
                    class="flex h-4 w-4 shrink-0 items-center justify-center bg-transparent text-[0.625rem] text-ink-muted"
                    onclick={() => toggleCol(col.id)}
                    aria-label={expandedCols.has(col.id) ? "Collapse" : "Expand"}
                >
                    {expandedCols.has(col.id) ? "▾" : "▸"}
                </button>
                <label class="flex flex-1 cursor-pointer items-center gap-2">
                    {#if showCollectionCheckbox}
                        <input
                            type="checkbox"
                            checked={cs !== "none"}
                            indeterminate={cs === "some"}
                            onchange={() => onCollectionToggle?.(col, cs)}
                        />
                    {/if}
                    <span class="text-sm text-accent">⬡</span>
                    <span class="flex-1 font-medium text-ink">{col.name}</span>
                    {#if cs === "some"}
                        <span class="text-[0.6875rem] italic text-ink-muted">partial</span>
                    {:else if cs === "all"}
                        <span class="text-[0.6875rem] text-accent">all on</span>
                    {/if}
                </label>
            </div>

            {#if expandedCols.has(col.id) && col.folders.length > 0}
                <div class="bg-cave-deep pb-1 pt-0.5">
                    {#each col.folders as folder (folder.id)}
                        {@const fs = folderState(col, folder)}
                        <div class="px-3">
                            <div class="flex items-center gap-2 py-1.5">
                                <button
                                    class="flex h-4 w-4 shrink-0 items-center justify-center bg-transparent text-[0.625rem] text-ink-muted"
                                    onclick={() => toggleFolder(folder.id)}
                                    aria-label={expandedFolders.has(folder.id)
                                        ? "Collapse"
                                        : "Expand"}
                                >
                                    {expandedFolders.has(folder.id) ? "▾" : "▸"}
                                </button>
                                <label class="flex flex-1 cursor-pointer items-center gap-2">
                                    <input
                                        type="checkbox"
                                        checked={fs !== "none"}
                                        indeterminate={fs === "some"}
                                        onchange={() => onFolderToggle?.(col, folder, fs)}
                                    />
                                    <span class="text-xs"
                                        >{expandedFolders.has(folder.id) ? "📂" : "📁"}</span
                                    >
                                    <span class="text-sm text-ink-mid">{folder.name}</span>
                                    {#if fs === "some"}
                                        <span class="ml-1 text-[0.6875rem] italic text-ink-muted"
                                            >partial</span
                                        >
                                    {/if}
                                </label>
                            </div>

                            {#if expandedFolders.has(folder.id)}
                                {#each folder.endpoints as ep (ep.id)}
                                    <label
                                        class="flex cursor-pointer items-center gap-2 py-1 pl-7 text-ink-mid"
                                    >
                                        <input
                                            type="checkbox"
                                            checked={isEndpointChecked(ep, col, folder)}
                                            onchange={() => onEndpointToggle(ep, col, folder)}
                                        />
                                        <span
                                            class="inline-block shrink-0 rounded-[3px] px-1.5 py-0.5 font-mono text-[0.625rem] font-bold uppercase"
                                            style="background:{mc(ep.method)};color:#1a1f30"
                                            >{ep.method === "*" ? "ANY" : ep.method}</span
                                        >
                                        <span class="min-w-0 flex-1 truncate font-mono text-xs"
                                            >{ep.path}</span
                                        >
                                        <span class="max-w-[140px] truncate text-xs text-ink-muted"
                                            >{ep.name}</span
                                        >
                                    </label>
                                {/each}
                            {/if}
                        </div>
                    {/each}
                </div>
            {/if}
        </div>
    {/each}
</div>
