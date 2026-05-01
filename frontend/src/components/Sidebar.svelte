<script lang="ts">
    import { newEndpoint } from "../lib/types";
    import type { Collection, Folder, Endpoint } from "../lib/types";
    import {
        collections,
        expandedCollections,
        expandedFolders,
        patchCollection,
        openTabs,
        activeTabId,
    } from "../stores/index";
    import type { EndpointTab } from "../stores/index";
    import * as api from "../lib/api";
    import { Button, MethodBadge, Modal } from "../lib/ui";
    import SidebarShell from "./SidebarShell.svelte";
    import AssignServersModal from "./AssignServersModal.svelte";

    let newCollectionName = $state("");
    let showNewCollection = $state(false);

    let newFolderModal = $state<{ collectionId: string } | null>(null);
    let newFolderName = $state("");

    type CtxKind = "collection" | "folder" | "endpoint";
    let ctxMenu = $state<{
        type: CtxKind;
        collectionId: string;
        folderId?: string;
        endpointId?: string;
        x: number;
        y: number;
    } | null>(null);

    let renameModal = $state<{
        type: CtxKind;
        collectionId: string;
        folderId?: string;
        endpointId?: string;
        currentName: string;
    } | null>(null);
    let renameValue = $state("");

    let assignTarget = $state<
        | { type: "collection"; collection: Collection }
        | { type: "folder"; collection: Collection; folder: Folder }
        | { type: "endpoint"; collection: Collection; folder: Folder; endpoint: Endpoint }
        | null
    >(null);

    function endpointBp(ep: Endpoint): boolean {
        return !!ep.breakpoint;
    }

    function folderBp(folder: Folder): boolean {
        return folder.endpoints.some((e) => e.breakpoint);
    }

    function collectionBp(col: Collection): boolean {
        return col.folders.some((f) => f.endpoints.some((e) => e.breakpoint));
    }

    function toggleCollection(id: string) {
        expandedCollections.update((s) => {
            const next = new Set(s);
            if (next.has(id)) next.delete(id);
            else next.add(id);
            return next;
        });
    }

    function toggleFolder(id: string) {
        expandedFolders.update((s) => {
            const next = new Set(s);
            if (next.has(id)) next.delete(id);
            else next.add(id);
            return next;
        });
    }

    function selectEndpoint(collectionId: string, folderId: string, endpointId: string) {
        const existing = $openTabs.find((t) => t.endpointId === endpointId);
        if (existing) {
            activeTabId.set(existing.tabId);
            return;
        }
        const tab: EndpointTab = { tabId: endpointId, endpointId, collectionId, folderId };
        openTabs.update((tabs) => [...tabs, tab]);
        activeTabId.set(endpointId);
    }

    async function createCollection() {
        if (!newCollectionName.trim()) return;
        const c = await api.createCollection(newCollectionName.trim(), "");
        collections.update((list) => [...list, c]);
        expandedCollections.update((s) => {
            const n = new Set(s);
            n.add(c.id);
            return n;
        });
        newCollectionName = "";
        showNewCollection = false;
    }

    function openNewFolderModal(collectionId: string) {
        newFolderModal = { collectionId };
        newFolderName = "";
    }

    async function createFolder() {
        if (!newFolderModal || !newFolderName.trim()) return;
        const updated = await api.createFolder(newFolderModal.collectionId, newFolderName.trim());
        patchCollection(updated);
        expandedFolders.update((s) => {
            const folder = updated.folders[updated.folders.length - 1];
            const n = new Set(s);
            n.add(folder.id);
            return n;
        });
        newFolderModal = null;
        newFolderName = "";
    }

    async function createEndpoint(collectionId: string, folderId: string) {
        const updated = await api.createEndpoint(collectionId, folderId, newEndpoint());
        patchCollection(updated);
        const folder = updated.folders.find((f) => f.id === folderId);
        if (folder) {
            const ep = folder.endpoints[folder.endpoints.length - 1];
            selectEndpoint(collectionId, folderId, ep.id);
        }
    }

    async function deleteCollection(id: string) {
        if (!confirm("Delete this collection?")) return;
        await api.deleteCollection(id);
        collections.update((list) => list.filter((c) => c.id !== id));
    }

    async function deleteFolder(collectionId: string, folderId: string) {
        if (!confirm("Delete this folder?")) return;
        const updated = await api.deleteFolder(collectionId, folderId);
        patchCollection(updated);
    }

    async function deleteEndpoint(collectionId: string, folderId: string, endpointId: string) {
        if (!confirm("Delete this endpoint?")) return;
        const updated = await api.deleteEndpoint(collectionId, folderId, endpointId);
        patchCollection(updated);
        openTabs.update((tabs) => tabs.filter((t) => t.endpointId !== endpointId));
    }

    function openCtxMenu(
        e: MouseEvent,
        type: CtxKind,
        collectionId: string,
        folderId?: string,
        endpointId?: string,
    ) {
        e.preventDefault();
        e.stopPropagation();
        ctxMenu = { type, collectionId, folderId, endpointId, x: e.clientX, y: e.clientY };
    }

    function closeCtxMenu() {
        ctxMenu = null;
    }

    function startRename() {
        if (!ctxMenu) return;
        const { type, collectionId, folderId, endpointId } = ctxMenu;
        let currentName = "";
        if (type === "collection") {
            currentName = $collections.find((c) => c.id === collectionId)?.name ?? "";
        } else if (type === "folder") {
            const col = $collections.find((c) => c.id === collectionId);
            currentName = col?.folders.find((f) => f.id === folderId)?.name ?? "";
        } else {
            const col = $collections.find((c) => c.id === collectionId);
            const folder = col?.folders.find((f) => f.id === folderId);
            currentName = folder?.endpoints.find((e) => e.id === endpointId)?.name ?? "";
        }
        renameModal = { type, collectionId, folderId, endpointId, currentName };
        renameValue = currentName;
        ctxMenu = null;
    }

    async function doRename() {
        if (!renameModal || !renameValue.trim()) return;
        const { type, collectionId, folderId, endpointId } = renameModal;
        if (type === "collection") {
            const col = $collections.find((c) => c.id === collectionId);
            if (!col) return;
            const updated = await api.updateCollection(collectionId, {
                name: renameValue.trim(),
                description: col.description,
                variables: col.variables ?? {},
            });
            collections.update((list) => list.map((c) => (c.id === collectionId ? updated : c)));
        } else if (type === "folder" && folderId) {
            const updated = await api.renameFolder(collectionId, folderId, renameValue.trim());
            patchCollection(updated);
        } else if (type === "endpoint" && folderId && endpointId) {
            const col = $collections.find((c) => c.id === collectionId);
            const folder = col?.folders.find((f) => f.id === folderId);
            const ep = folder?.endpoints.find((e) => e.id === endpointId);
            if (!ep) return;
            const updated = await api.updateEndpoint(collectionId, folderId, endpointId, {
                ...ep,
                name: renameValue.trim(),
            });
            patchCollection(updated);
        }
        renameModal = null;
    }

    function ctxDelete() {
        const m = ctxMenu;
        ctxMenu = null;
        if (!m) return;
        if (m.type === "collection") deleteCollection(m.collectionId);
        else if (m.type === "folder" && m.folderId) deleteFolder(m.collectionId, m.folderId);
        else if (m.type === "endpoint" && m.folderId && m.endpointId)
            deleteEndpoint(m.collectionId, m.folderId, m.endpointId);
    }

    function ctxAssignServers() {
        const m = ctxMenu;
        ctxMenu = null;
        if (!m) return;
        const col = $collections.find((c) => c.id === m.collectionId);
        if (!col) return;
        if (m.type === "collection") {
            assignTarget = { type: "collection", collection: col };
        } else if (m.type === "folder" && m.folderId) {
            const folder = col.folders.find((f) => f.id === m.folderId);
            if (folder) assignTarget = { type: "folder", collection: col, folder };
        } else if (m.type === "endpoint" && m.folderId && m.endpointId) {
            const folder = col.folders.find((f) => f.id === m.folderId);
            const ep = folder?.endpoints.find((e) => e.id === m.endpointId);
            if (folder && ep)
                assignTarget = { type: "endpoint", collection: col, folder, endpoint: ep };
        }
    }

    async function ctxToggleBreakpoint() {
        const m = ctxMenu;
        ctxMenu = null;
        if (!m) return;
        const col = $collections.find((c) => c.id === m.collectionId);
        if (!col) return;

        if (m.type === "endpoint" && m.folderId && m.endpointId) {
            const folder = col.folders.find((f) => f.id === m.folderId);
            const ep = folder?.endpoints.find((e) => e.id === m.endpointId);
            if (!folder || !ep) return;
            const updated = await api.updateEndpoint(col.id, folder.id, ep.id, {
                ...ep,
                breakpoint: !ep.breakpoint,
            });
            patchCollection(updated);
            return;
        }

        if (m.type === "folder" && m.folderId) {
            const folder = col.folders.find((f) => f.id === m.folderId);
            if (!folder) return;
            const on = !folderBp(folder);
            for (const ep of folder.endpoints) {
                if (ep.breakpoint !== on) {
                    const updated = await api.updateEndpoint(col.id, folder.id, ep.id, {
                        ...ep,
                        breakpoint: on,
                    });
                    patchCollection(updated);
                }
            }
            return;
        }

        if (m.type === "collection") {
            const on = !collectionBp(col);
            for (const folder of col.folders) {
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
        }
    }

    let ctxBpLabel = $derived.by(() => {
        const m = ctxMenu;
        if (!m) return "";
        const col = $collections.find((c) => c.id === m.collectionId);
        if (!col) return "";
        if (m.type === "endpoint" && m.folderId && m.endpointId) {
            const ep = col.folders
                .find((f) => f.id === m.folderId)
                ?.endpoints.find((e) => e.id === m.endpointId);
            return ep?.breakpoint ? "Disable breakpoint" : "Enable breakpoint";
        }
        if (m.type === "folder" && m.folderId) {
            const folder = col.folders.find((f) => f.id === m.folderId);
            return folder && folderBp(folder) ? "Disable breakpoints" : "Enable breakpoints";
        }
        return collectionBp(col) ? "Disable breakpoints" : "Enable breakpoints";
    });
</script>

<svelte:window onclick={closeCtxMenu} />

<SidebarShell title="Collections">
    {#snippet actions()}
        <Button
            variant="icon"
            onclick={() => (showNewCollection = !showNewCollection)}
            title="New Collection">+</Button
        >
    {/snippet}

    {#if showNewCollection}
        <div
            class="flex shrink-0 flex-col gap-1.5 border-b border-wire bg-cave-surface px-2.5 py-2"
        >
            <input
                bind:value={newCollectionName}
                placeholder="Collection name"
                onkeydown={(e) => {
                    if (e.key === "Enter") createCollection();
                    if (e.key === "Escape") showNewCollection = false;
                }}
            />
            <div class="flex gap-1">
                <Button variant="primary" size="sm" onclick={createCollection}>Create</Button>
                <Button size="sm" onclick={() => (showNewCollection = false)}>Cancel</Button>
            </div>
        </div>
    {/if}

    <div class="flex-1 overflow-y-auto py-1">
        {#each $collections as collection (collection.id)}
            <div>
                <div
                    class="group flex min-h-[28px] items-center gap-0.5 px-2 py-1 text-ink-mid hover:bg-cave-raised hover:text-ink"
                    role="treeitem"
                    aria-selected="false"
                    aria-expanded={$expandedCollections.has(collection.id)}
                    tabindex="-1"
                    oncontextmenu={(e) => openCtxMenu(e, "collection", collection.id)}
                >
                    <button
                        class="flex h-4 w-4 shrink-0 items-center justify-center bg-transparent text-[0.625rem] text-ink-muted"
                        onclick={() => toggleCollection(collection.id)}
                        aria-label="Toggle collection"
                    >
                        {$expandedCollections.has(collection.id) ? "▾" : "▸"}
                    </button>
                    <span class="shrink-0 text-xs text-accent">⬡</span>
                    <span class="flex-1 truncate text-[0.8125rem] font-medium"
                        >{collection.name}</span
                    >
                    {#if collectionBp(collection)}
                        <span
                            class="h-2 w-2 shrink-0 rounded-[1px] bg-err"
                            title="Breakpoint enabled"
                        ></span>
                    {/if}
                    <div class="flex gap-0.5 opacity-0 group-hover:opacity-100">
                        <button
                            class="rounded-sm bg-transparent p-1 text-ink-muted hover:bg-cave-raised hover:text-ink"
                            title="Add folder"
                            onclick={(e) => {
                                e.stopPropagation();
                                openNewFolderModal(collection.id);
                            }}
                            aria-label="Add folder"
                        >
                            <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
                                <path
                                    d="M1 3.5A1.5 1.5 0 012.5 2h3.764c.958 0 1.76.56 2.177 1.368l.105.21c.175.35.53.572.93.572H13.5A1.5 1.5 0 0115 5.5v7A1.5 1.5 0 0113.5 14h-11A1.5 1.5 0 011 12.5v-9z"
                                />
                            </svg>
                        </button>
                        <button
                            class="rounded-sm bg-transparent p-1 text-ink-muted hover:bg-cave-raised hover:text-err"
                            title="Delete"
                            onclick={(e) => {
                                e.stopPropagation();
                                deleteCollection(collection.id);
                            }}
                            aria-label="Delete collection">✕</button
                        >
                    </div>
                </div>

                {#if $expandedCollections.has(collection.id)}
                    {#each collection.folders as folder (folder.id)}
                        <div>
                            <div
                                class="group flex min-h-[28px] items-center gap-0.5 py-1 pl-6 pr-2 text-ink-mid hover:bg-cave-raised hover:text-ink"
                                role="treeitem"
                                aria-selected="false"
                                aria-expanded={$expandedFolders.has(folder.id)}
                                tabindex="-1"
                                oncontextmenu={(e) =>
                                    openCtxMenu(e, "folder", collection.id, folder.id)}
                            >
                                <button
                                    class="flex h-4 w-4 shrink-0 items-center justify-center bg-transparent text-[0.625rem] text-ink-muted"
                                    onclick={() => toggleFolder(folder.id)}
                                    aria-label="Toggle folder"
                                >
                                    {$expandedFolders.has(folder.id) ? "▾" : "▸"}
                                </button>
                                <span class="shrink-0 text-xs"
                                    >{$expandedFolders.has(folder.id) ? "📂" : "📁"}</span
                                >
                                <span class="flex-1 truncate text-[0.8125rem] font-medium"
                                    >{folder.name}</span
                                >
                                {#if folderBp(folder)}
                                    <span
                                        class="h-2 w-2 shrink-0 rounded-[1px] bg-err"
                                        title="Breakpoint enabled"
                                    ></span>
                                {/if}
                                <div class="flex gap-0.5 opacity-0 group-hover:opacity-100">
                                    <button
                                        class="rounded-sm bg-transparent p-1 text-ink-muted hover:bg-cave-raised hover:text-ink"
                                        title="Add endpoint"
                                        onclick={(e) => {
                                            e.stopPropagation();
                                            createEndpoint(collection.id, folder.id);
                                        }}
                                        aria-label="Add endpoint">+</button
                                    >
                                    <button
                                        class="rounded-sm bg-transparent p-1 text-ink-muted hover:bg-cave-raised hover:text-err"
                                        title="Delete folder"
                                        onclick={(e) => {
                                            e.stopPropagation();
                                            deleteFolder(collection.id, folder.id);
                                        }}
                                        aria-label="Delete folder">✕</button
                                    >
                                </div>
                            </div>

                            {#if $expandedFolders.has(folder.id)}
                                {#each folder.endpoints as ep (ep.id)}
                                    <button
                                        class="group flex min-h-[28px] w-full items-center gap-1.5 py-1 pl-10 pr-2 text-left transition-colors
                      {$activeTabId === ep.id
                                            ? 'bg-accent/10 text-ink'
                                            : $openTabs.some((t) => t.endpointId === ep.id)
                                              ? 'text-ink hover:bg-cave-raised'
                                              : 'text-ink-mid hover:bg-cave-raised hover:text-ink'}"
                                        onclick={() =>
                                            selectEndpoint(collection.id, folder.id, ep.id)}
                                        oncontextmenu={(e) => {
                                            e.preventDefault();
                                            e.stopPropagation();
                                            openCtxMenu(
                                                e,
                                                "endpoint",
                                                collection.id,
                                                folder.id,
                                                ep.id,
                                            );
                                        }}
                                    >
                                        <MethodBadge method={ep.method} />
                                        <span class="flex-1 truncate text-[0.8125rem] font-medium"
                                            >{ep.name || ep.path || "Unnamed"}</span
                                        >
                                        {#if endpointBp(ep)}
                                            <span
                                                class="h-2 w-2 shrink-0 rounded-[1px] bg-err"
                                                title="Breakpoint enabled"
                                            ></span>
                                        {/if}
                                        <span
                                            class="rounded-sm bg-transparent p-0.5 text-[0.625rem] text-ink-muted opacity-0 hover:bg-cave-raised hover:text-err group-hover:opacity-100"
                                            role="button"
                                            tabindex="0"
                                            onclick={(e) => {
                                                e.stopPropagation();
                                                deleteEndpoint(collection.id, folder.id, ep.id);
                                            }}
                                            onkeydown={(e) => {
                                                if (e.key === "Enter" || e.key === " ") {
                                                    e.preventDefault();
                                                    e.stopPropagation();
                                                    deleteEndpoint(collection.id, folder.id, ep.id);
                                                }
                                            }}>✕</span
                                        >
                                    </button>
                                {/each}
                                {#if folder.endpoints.length === 0}
                                    <div class="py-1 pl-12 text-[0.6875rem] italic text-ink-muted">
                                        No endpoints
                                    </div>
                                {/if}
                            {/if}
                        </div>
                    {/each}

                    {#if collection.folders.length === 0}
                        <div class="flex flex-col items-start gap-1.5 py-2 pl-8 pr-2">
                            <span class="text-xs text-ink-muted">No folders yet</span>
                            <Button size="sm" onclick={() => openNewFolderModal(collection.id)}
                                >+ Add folder</Button
                            >
                        </div>
                    {/if}
                {/if}
            </div>
        {/each}

        {#if $collections.length === 0}
            <div class="flex flex-col items-center gap-2.5 px-5 py-10 text-center">
                <span class="text-3xl opacity-20">⬡</span>
                <p class="text-[0.8125rem] text-ink-muted">No collections yet.</p>
                <Button variant="primary" onclick={() => (showNewCollection = true)}
                    >+ New Collection</Button
                >
            </div>
        {/if}
    </div>
</SidebarShell>

{#if ctxMenu}
    <div
        class="fixed z-[2000] min-w-[160px] overflow-hidden rounded border border-wire bg-cave-elevated shadow-modal"
        style="left:{ctxMenu.x}px;top:{ctxMenu.y}px"
        onclick={(e) => e.stopPropagation()}
        onkeydown={(e) => e.stopPropagation()}
        role="menu"
        tabindex="-1"
    >
        <button
            class="flex w-full items-center bg-transparent px-3 py-2 text-left text-xs text-ink-mid hover:bg-cave-raised hover:text-ink"
            onclick={startRename}>Rename</button
        >
        <button
            class="flex w-full items-center bg-transparent px-3 py-2 text-left text-xs text-ink-mid hover:bg-cave-raised hover:text-ink"
            onclick={ctxAssignServers}>Assign servers…</button
        >
        <button
            class="flex w-full items-center bg-transparent px-3 py-2 text-left text-xs text-ink-mid hover:bg-cave-raised hover:text-ink"
            onclick={ctxToggleBreakpoint}>{ctxBpLabel}</button
        >
        <div class="border-t border-wire"></div>
        <button
            class="flex w-full items-center bg-transparent px-3 py-2 text-left text-xs text-err hover:bg-err/10"
            onclick={ctxDelete}>Delete</button
        >
    </div>
{/if}

{#if renameModal}
    <Modal open={true} title="Rename" onClose={() => (renameModal = null)}>
        <div class="flex flex-col gap-3">
            <input
                bind:value={renameValue}
                onkeydown={(e) => {
                    if (e.key === "Enter") doRename();
                    if (e.key === "Escape") renameModal = null;
                }}
            />
            <div class="flex justify-end gap-2">
                <Button variant="primary" onclick={doRename}>Rename</Button>
                <Button onclick={() => (renameModal = null)}>Cancel</Button>
            </div>
        </div>
    </Modal>
{/if}

{#if newFolderModal}
    <Modal
        open={true}
        title="New Folder"
        onClose={() => {
            newFolderModal = null;
            newFolderName = "";
        }}
    >
        <div class="flex flex-col gap-3">
            <input
                bind:value={newFolderName}
                placeholder="Folder name"
                onkeydown={(e) => {
                    if (e.key === "Enter") createFolder();
                    if (e.key === "Escape") {
                        newFolderModal = null;
                        newFolderName = "";
                    }
                }}
            />
            <div class="flex justify-end gap-2">
                <Button variant="primary" onclick={createFolder}>Create</Button>
                <Button
                    onclick={() => {
                        newFolderModal = null;
                        newFolderName = "";
                    }}>Cancel</Button
                >
            </div>
        </div>
    </Modal>
{/if}

{#if assignTarget}
    {@const t = assignTarget}
    <AssignServersModal target={t} onClose={() => (assignTarget = null)} />
{/if}
