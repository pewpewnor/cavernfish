<script lang="ts">
    import type { Endpoint } from "../lib/types";
    import { METHOD_COLORS } from "../lib/types";
    import { tabEdits, activeTabId } from "../stores/index";
    import ResponseEditor from "./ResponseEditor.svelte";
    import { Button, TabStrip, Toggle } from "../lib/ui";

    interface Props {
        endpoint: Endpoint;
        collectionId: string;
        folderId: string;
        tabId: string;
        onSave: (data: { collectionId: string; folderId: string; endpoint: Endpoint }) => void;
        onDirtyChange: (dirty: boolean) => void;
    }

    let { endpoint, collectionId, folderId, tabId, onSave, onDirtyChange }: Props = $props();

    const METHODS = ["GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS", "*"];

    type EditorTab = "body" | "headers" | "cookies" | "configure";

    // svelte-ignore state_referenced_locally
    let editingEndpoint: Endpoint = $state(JSON.parse(JSON.stringify(endpoint)));
    let activeTab = $state<EditorTab>("body");
    let dirty = $state(false);
    // svelte-ignore state_referenced_locally
    let prevEndpointId = $state(endpoint.id);

    $effect(() => {
        if (endpoint.id !== prevEndpointId) {
            prevEndpointId = endpoint.id;
            editingEndpoint = JSON.parse(JSON.stringify(endpoint));
            activeTab = "body";
            dirty = false;
        }
    });

    $effect(() => {
        tabEdits.update((m) => ({ ...m, [tabId]: editingEndpoint }));
    });

    function markDirty() {
        if (!dirty) {
            dirty = true;
            onDirtyChange(true);
        }
    }

    function save() {
        onSave({ collectionId, folderId, endpoint: $state.snapshot(editingEndpoint) });
        dirty = false;
        onDirtyChange(false);
    }

    function mc(m: string) {
        return METHOD_COLORS[m.toUpperCase()] || "#64748b";
    }

    function prettifyBody(body: string, bodyType: string): string {
        if (bodyType === "json") {
            try {
                return JSON.stringify(JSON.parse(body), null, 4);
            } catch {
                return body;
            }
        }
        if (bodyType === "xml" || bodyType === "html") {
            try {
                const mime = bodyType === "html" ? "text/html" : "application/xml";
                const parser = new DOMParser();
                const doc = parser.parseFromString(body.trim(), mime);
                const root = bodyType === "html" ? doc.body : doc.documentElement;
                if (!root) return body;
                return formatXmlNode(root, 0);
            } catch {
                return body;
            }
        }
        return body;
    }

    function formatXmlNode(node: Node, depth: number): string {
        const pad = "    ".repeat(depth);
        if (node.nodeType === Node.TEXT_NODE) {
            return node.textContent?.trim() ?? "";
        }
        if (node.nodeType === Node.ELEMENT_NODE) {
            const el = node as Element;
            const attrs = Array.from(el.attributes)
                .map((a) => ` ${a.name}="${a.value}"`)
                .join("");
            const tag = el.tagName.toLowerCase();
            const nonEmpty = Array.from(el.childNodes).filter(
                (n) => n.nodeType !== Node.TEXT_NODE || (n.textContent?.trim() ?? ""),
            );
            if (nonEmpty.length === 0) return `${pad}<${tag}${attrs}/>`;
            const hasEls = nonEmpty.some((n) => n.nodeType === Node.ELEMENT_NODE);
            if (!hasEls) return `${pad}<${tag}${attrs}>${el.textContent?.trim() ?? ""}</${tag}>`;
            const children = nonEmpty
                .map((c) => formatXmlNode(c, depth + 1))
                .filter(Boolean)
                .join("\n");
            return `${pad}<${tag}${attrs}>\n${children}\n${pad}</${tag}>`;
        }
        return "";
    }

    function handleKeydown(e: KeyboardEvent) {
        if (e.ctrlKey && e.key === "s" && $activeTabId === tabId) {
            e.preventDefault();
            if (dirty) save();
        }
    }

    let tabs = $derived<{ value: EditorTab; label: string }[]>([
        { value: "body", label: "Body" },
        {
            value: "headers",
            label:
                "Headers" +
                (editingEndpoint.headers.filter((h) => h.enabled && h.key).length > 0
                    ? ` (${editingEndpoint.headers.filter((h) => h.enabled && h.key).length})`
                    : ""),
        },
        {
            value: "cookies",
            label:
                "Cookies" +
                (editingEndpoint.cookies.length > 0 ? ` (${editingEndpoint.cookies.length})` : ""),
        },
        { value: "configure", label: "Configure" },
    ]);

    let canPrettify = $derived(
        activeTab === "body" && ["json", "xml", "html"].includes(editingEndpoint.bodyType),
    );

    function doPrettify() {
        const formatted = prettifyBody(editingEndpoint.body, editingEndpoint.bodyType);
        if (formatted === editingEndpoint.body) return;
        editingEndpoint = { ...editingEndpoint, body: formatted };
        markDirty();
    }
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="flex h-full flex-col overflow-hidden">
    <div class="flex shrink-0 items-center gap-2 border-b border-wire bg-cave-surface px-4 py-3">
        <select
            class="w-28 shrink-0 font-mono font-bold"
            value={editingEndpoint.method}
            style="color:{mc(editingEndpoint.method)}"
            onchange={(e) => {
                editingEndpoint.method = e.currentTarget.value;
                markDirty();
            }}
        >
            {#each METHODS as m (m)}
                <option value={m} style="color:{mc(m)}">{m === "*" ? "ANY" : m}</option>
            {/each}
        </select>

        <input
            class="flex-1 font-mono"
            placeholder="/api/path/:param"
            bind:value={editingEndpoint.path}
            oninput={markDirty}
        />

        <div class="flex shrink-0 items-center gap-2">
            {#if dirty}
                <span class="h-2.5 w-2.5 rounded-full bg-warn" title="Unsaved changes"></span>
            {/if}
            <Button variant="primary" onclick={save} disabled={!dirty}>
                {dirty ? "Save" : "Saved ✓"}
            </Button>
        </div>
    </div>

    {#snippet prettifyBtn()}
        <Button size="sm" onclick={doPrettify}>Prettify</Button>
    {/snippet}

    <TabStrip bind:value={activeTab} {tabs} trailing={canPrettify ? prettifyBtn : undefined} />

    <div class="flex flex-1 flex-col overflow-hidden">
        {#if activeTab === "configure"}
            <div class="flex flex-col gap-5 overflow-y-auto p-5">
                <div class="flex flex-wrap gap-6">
                    <div class="flex flex-col gap-1.5">
                        <span class="form-label">Delay (ms)</span>
                        <input
                            type="number"
                            min="0"
                            max="60000"
                            class="w-28"
                            bind:value={editingEndpoint.delayMs}
                            oninput={markDirty}
                        />
                    </div>

                    <div class="flex flex-col gap-1.5 self-end pb-1">
                        <Toggle bind:checked={editingEndpoint.wsEnabled} label="WebSocket" />
                    </div>
                </div>

                <div class="flex flex-col gap-1.5">
                    <span class="form-label">Proxy URL</span>
                    {#if editingEndpoint.proxyUrl !== undefined}
                        <div class="flex items-center gap-2">
                            <input
                                class="flex-1 font-mono text-sm"
                                placeholder="https://api.example.com"
                                bind:value={editingEndpoint.proxyUrl}
                                oninput={markDirty}
                            />
                            <Button
                                variant="icon"
                                onclick={() => {
                                    editingEndpoint.proxyUrl = undefined;
                                    markDirty();
                                }}>✕</Button
                            >
                        </div>
                    {:else}
                        <div>
                            <Button
                                onclick={() => {
                                    editingEndpoint.proxyUrl = "";
                                    markDirty();
                                }}>+ Add Proxy</Button
                            >
                        </div>
                    {/if}
                </div>
            </div>
        {:else}
            <ResponseEditor
                response={{
                    statusCode: editingEndpoint.statusCode,
                    body: editingEndpoint.body,
                    bodyType: editingEndpoint.bodyType,
                    headers: editingEndpoint.headers,
                    cookies: editingEndpoint.cookies,
                }}
                {activeTab}
                showStatusBar={true}
                onChange={(r) => {
                    editingEndpoint = { ...editingEndpoint, ...r };
                    markDirty();
                }}
            />
        {/if}
    </div>
</div>
