<script lang="ts">
    import { collections, patchCollection } from "../stores/index";
    import * as api from "../lib/api";
    import { Button, Modal, TabStrip } from "../lib/ui";

    let showImport = $state(false);
    let importTab = $state<"native" | "openapi">("native");
    let importText = $state("");
    let importError = $state("");
    let importLoading = $state(false);

    let editingCollection = $state<string | null>(null);
    let editingVars = $state<{ key: string; value: string }[]>([]);
    let editingName = $state("");
    let editingDesc = $state("");

    function openVarEditor(cid: string) {
        const col = $collections.find((c) => c.id === cid);
        if (!col) return;
        editingCollection = cid;
        editingName = col.name;
        editingDesc = col.description;
        editingVars = Object.entries(col.variables ?? {}).map(([key, value]) => ({ key, value }));
    }

    async function saveVars() {
        if (!editingCollection) return;
        const variables: Record<string, string> = {};
        editingVars.forEach((v) => {
            if (v.key.trim()) variables[v.key.trim()] = v.value;
        });
        const updated = await api.updateCollection(editingCollection, {
            name: editingName,
            description: editingDesc,
            variables,
        });
        patchCollection(updated);
        editingCollection = null;
    }

    async function exportCollection(id: string) {
        const data = await api.exportCollection(id);
        const blob = new Blob([data], { type: "application/json" });
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        const col = $collections.find((c) => c.id === id);
        a.href = url;
        a.download = `${col?.name ?? "collection"}.cavernfish.json`;
        a.click();
        URL.revokeObjectURL(url);
    }

    async function doImport() {
        if (!importText.trim()) return;
        importError = "";
        importLoading = true;
        try {
            let result;
            if (importTab === "openapi") {
                result = await api.importFromOpenAPI(importText);
            } else {
                result = await api.importCollection(importText);
            }
            if (result) {
                collections.update((list) => [...list, result]);
                showImport = false;
                importText = "";
            }
        } catch (e: unknown) {
            importError = e instanceof Error ? e.message : String(e);
        } finally {
            importLoading = false;
        }
    }

    async function handleFileImport(e: Event) {
        const file = (e.target as HTMLInputElement).files?.[0];
        if (!file) return;
        importText = await file.text();
    }
</script>

<div class="flex flex-1 flex-col gap-6 overflow-y-auto p-5">
    <div class="shrink-0">
        <h2 class="text-lg font-semibold text-ink">Settings</h2>
    </div>

    <div class="mx-auto flex w-full max-w-3xl flex-col gap-6">
        <section class="overflow-hidden rounded-lg border border-wire bg-cave-surface">
            <div class="section-header">
                <span class="section-title">Collections</span>
                <Button variant="primary" onclick={() => (showImport = true)}>⬆ Import</Button>
            </div>

            {#if $collections.length === 0}
                <div class="px-4 py-5 text-[0.8125rem] text-ink-muted">No collections yet.</div>
            {:else}
                <div class="flex flex-col">
                    <div
                        class="grid grid-cols-[1fr_80px_80px_80px_180px] gap-3 border-b border-wire px-4 py-2 text-[0.625rem] font-semibold uppercase tracking-[0.06em] text-ink-muted"
                    >
                        <span>Name</span>
                        <span>Folders</span>
                        <span>Endpoints</span>
                        <span>Variables</span>
                        <span></span>
                    </div>
                    {#each $collections as col (col.id)}
                        {@const totalEndpoints = col.folders.reduce(
                            (s, f) => s + f.endpoints.length,
                            0,
                        )}
                        <div
                            class="grid grid-cols-[1fr_80px_80px_80px_180px] items-center gap-3 border-b border-wire px-4 py-2.5 transition-colors last:border-b-0 hover:bg-cave-raised"
                        >
                            <div class="flex min-w-0 items-center gap-2">
                                <span class="shrink-0 text-base text-accent">⬡</span>
                                <div class="min-w-0">
                                    <div class="truncate text-[0.8125rem] font-medium text-ink">
                                        {col.name}
                                    </div>
                                    {#if col.description}
                                        <div class="truncate text-[0.6875rem] text-ink-muted">
                                            {col.description}
                                        </div>
                                    {/if}
                                </div>
                            </div>
                            <span class="font-mono text-[0.8125rem] text-ink-mid"
                                >{col.folders.length}</span
                            >
                            <span class="font-mono text-[0.8125rem] text-ink-mid"
                                >{totalEndpoints}</span
                            >
                            <span class="font-mono text-[0.8125rem] text-ink-mid"
                                >{Object.keys(col.variables ?? {}).length}</span
                            >
                            <div class="flex justify-end gap-1.5">
                                <Button size="sm" onclick={() => openVarEditor(col.id)}
                                    >✏ Edit</Button
                                >
                                <Button size="sm" onclick={() => exportCollection(col.id)}
                                    >⬇ Export</Button
                                >
                            </div>
                        </div>
                    {/each}
                </div>
            {/if}
        </section>

        <section class="overflow-hidden rounded-lg border border-wire bg-cave-surface">
            <div class="section-header">
                <span class="section-title">About</span>
            </div>
            <div class="flex flex-col gap-2 px-4 py-3">
                <div class="flex items-center gap-4">
                    <span class="w-28 shrink-0 text-xs text-ink-muted">App</span>
                    <span class="text-[0.8125rem] text-ink">Cavernfish</span>
                </div>
                <div class="flex items-center gap-4">
                    <span class="w-28 shrink-0 text-xs text-ink-muted">Version</span>
                    <span class="font-mono text-[0.8125rem] text-ink">0.1.0</span>
                </div>
                <div class="flex items-center gap-4">
                    <span class="w-28 shrink-0 text-xs text-ink-muted">Data directory</span>
                    <span class="font-mono text-[0.8125rem] text-ink">~/.cavernfish/</span>
                </div>
            </div>
        </section>
    </div>
</div>

{#if showImport}
    <Modal open={true} title="Import Collection" width="560px" onClose={() => (showImport = false)}>
        <div class="flex flex-col gap-3">
            <TabStrip
                bind:value={importTab}
                tabs={[
                    { value: "native", label: "Cavernfish JSON" },
                    { value: "openapi", label: "OpenAPI 3.x" },
                ]}
            />

            <div class="flex items-center gap-2.5 pt-2">
                <label
                    class="inline-flex cursor-pointer items-center gap-1.5 rounded-sm border border-wire bg-cave-raised px-3 py-1.5 text-xs text-ink-mid transition-all hover:border-wire-hi hover:text-ink"
                >
                    📂 Choose file
                    <input
                        type="file"
                        class="hidden"
                        accept=".json,.yaml,.yml"
                        onchange={handleFileImport}
                    />
                </label>
                <span class="text-[0.6875rem] text-ink-muted">or paste below</span>
            </div>

            <textarea
                class="w-full resize-y font-mono text-[0.6875rem] leading-relaxed"
                bind:value={importText}
                placeholder={importTab === "openapi"
                    ? '{ "openapi": "3.0.0", ... }'
                    : '{ "name": "My Collection", "folders": [...] }'}
                rows="10"
            ></textarea>

            {#if importError}
                <div
                    class="rounded-sm border border-err/20 bg-err/[.08] px-2.5 py-2 text-xs text-err"
                >
                    {importError}
                </div>
            {/if}

            <div class="flex justify-end gap-2">
                <Button
                    variant="primary"
                    onclick={doImport}
                    disabled={importLoading || !importText.trim()}
                >
                    {importLoading ? "Importing…" : "Import"}
                </Button>
                <Button onclick={() => (showImport = false)}>Cancel</Button>
            </div>
        </div>
    </Modal>
{/if}

{#if editingCollection}
    <Modal
        open={true}
        title="Edit Collection"
        width="500px"
        onClose={() => (editingCollection = null)}
    >
        <div class="flex flex-col gap-3.5">
            <label class="flex flex-col gap-1.5">
                <span class="form-label">Name</span>
                <input bind:value={editingName} />
            </label>

            <label class="flex flex-col gap-1.5">
                <span class="form-label">Description</span>
                <input bind:value={editingDesc} placeholder="Optional description" />
            </label>

            <div class="flex flex-col gap-1.5">
                <div class="mb-0.5 flex items-baseline gap-2.5">
                    <span class="form-label">Variables</span>
                    <span class="text-[0.6875rem] text-ink-muted"
                        >Use &#123;&#123;variable_name&#125;&#125; in paths and bodies</span
                    >
                </div>
                {#each editingVars as v, i (i)}
                    <div class="flex items-center gap-1.5">
                        <input
                            class="flex-1 font-mono"
                            placeholder="variable_name"
                            bind:value={v.key}
                        />
                        <span class="shrink-0 font-mono text-sm text-ink-muted">=</span>
                        <input class="flex-1 font-mono" placeholder="value" bind:value={v.value} />
                        <Button
                            variant="icon"
                            onclick={() =>
                                (editingVars = editingVars.filter((_, idx) => idx !== i))}>✕</Button
                        >
                    </div>
                {/each}
                <Button
                    size="sm"
                    onclick={() => (editingVars = [...editingVars, { key: "", value: "" }])}
                    >+ Add Variable</Button
                >
            </div>

            <div class="flex justify-end gap-2">
                <Button variant="primary" onclick={saveVars}>Save</Button>
                <Button onclick={() => (editingCollection = null)}>Cancel</Button>
            </div>
        </div>
    </Modal>
{/if}
