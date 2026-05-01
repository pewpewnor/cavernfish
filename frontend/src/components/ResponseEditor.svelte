<script lang="ts">
    import type { BodyType, KVPair, MockCookie } from "../lib/types";
    import { statusColor, newCookie } from "../lib/types";
    import JsonEditor from "./JsonEditor.svelte";
    import KeyValueEditor from "./KeyValueEditor.svelte";
    import { Button } from "../lib/ui";

    interface ResponseShape {
        statusCode: number;
        body: string;
        bodyType: BodyType;
        headers: KVPair[];
        cookies: MockCookie[];
    }

    interface Props {
        response: ResponseShape;
        activeTab: "body" | "headers" | "cookies";
        showStatusBar?: boolean;
        onChange: (next: ResponseShape) => void;
    }

    let { response, activeTab, showStatusBar = true, onChange }: Props = $props();

    function patch(p: Partial<ResponseShape>) {
        onChange({ ...response, ...p });
    }

    function addCookie() {
        onChange({ ...response, cookies: [...response.cookies, newCookie()] });
    }

    function updateCookie(i: number, p: Partial<MockCookie>) {
        onChange({
            ...response,
            cookies: response.cookies.map((c, idx) => (idx === i ? { ...c, ...p } : c)),
        });
    }

    function removeCookie(i: number) {
        onChange({ ...response, cookies: response.cookies.filter((_, idx) => idx !== i) });
    }

    let validity = $derived.by<"valid" | "invalid" | "empty" | "na">(() => {
        if (response.bodyType === "none" || response.bodyType === "text") return "na";
        if (!response.body.trim()) return "empty";
        if (response.bodyType === "json") {
            try {
                JSON.parse(response.body);
                return "valid";
            } catch {
                return "invalid";
            }
        }
        try {
            const mime = response.bodyType === "html" ? "text/html" : "application/xml";
            const doc = new DOMParser().parseFromString(response.body, mime);
            return doc.getElementsByTagName("parsererror").length > 0 ? "invalid" : "valid";
        } catch {
            return "invalid";
        }
    });

    let autoHeaders = $derived.by<KVPair[]>(() => {
        if (response.bodyType === "none") return [];
        const userKeys = new Set(
            response.headers.filter((h) => h.enabled && h.key).map((h) => h.key.toLowerCase()),
        );
        if (userKeys.has("content-type")) return [];
        const ct =
            response.bodyType === "json"
                ? "application/json"
                : response.bodyType === "html"
                  ? "text/html"
                  : response.bodyType === "xml"
                    ? "application/xml"
                    : "text/plain";
        return [{ key: "Content-Type", value: ct, enabled: true }];
    });
</script>

<div class="flex h-full flex-col">
    {#if showStatusBar}
        <div
            class="flex flex-wrap items-center gap-4 border-b border-wire bg-cave-surface px-4 py-2.5"
        >
            <div class="flex flex-col gap-1">
                <span class="form-label">Status</span>
                <input
                    type="number"
                    min="100"
                    max="599"
                    class="w-16 border-0 bg-transparent p-0 font-mono text-xl font-bold focus:ring-0"
                    value={response.statusCode}
                    onchange={(e) => patch({ statusCode: parseInt(e.currentTarget.value) || 200 })}
                    style="color:{statusColor(response.statusCode)}"
                />
            </div>

            <div class="flex flex-col gap-1">
                <span class="form-label">Body Type</span>
                <div class="flex items-center gap-2">
                    <select
                        value={response.bodyType}
                        onchange={(e) => patch({ bodyType: e.currentTarget.value as BodyType })}
                    >
                        <option value="json">JSON</option>
                        <option value="text">Text</option>
                        <option value="html">HTML</option>
                        <option value="xml">XML</option>
                        <option value="none">None</option>
                    </select>
                    {#if validity !== "na" && validity !== "empty"}
                        <span
                            class="flex items-center gap-1.5 rounded-sm px-1.5 py-0.5 text-[0.6875rem] font-medium {validity ===
                            'valid'
                                ? 'bg-ok/10 text-ok'
                                : 'bg-err/10 text-err'}"
                            title={validity === "valid"
                                ? "Body parses cleanly"
                                : "Body cannot be parsed"}
                        >
                            <span
                                class="h-1.5 w-1.5 rounded-full {validity === 'valid'
                                    ? 'bg-ok'
                                    : 'bg-err'}"
                            ></span>
                            {validity === "valid" ? "Valid" : "Invalid"}
                        </span>
                    {/if}
                </div>
            </div>
        </div>
    {/if}

    <div class="flex flex-1 flex-col overflow-hidden">
        {#if activeTab === "body"}
            {#if response.bodyType === "none"}
                <div class="flex flex-1 items-center justify-center py-10 text-sm text-ink-muted">
                    No body
                </div>
            {:else}
                <JsonEditor value={response.body} onChange={(v) => patch({ body: v })} />
            {/if}
        {:else if activeTab === "headers"}
            <div class="overflow-y-auto p-3">
                <KeyValueEditor
                    pairs={response.headers}
                    extras={autoHeaders}
                    extrasLabel="Auto-added by server"
                    keyPlaceholder="Header name"
                    valuePlaceholder="Value"
                    onChange={(headers) => patch({ headers })}
                />
            </div>
        {:else if activeTab === "cookies"}
            <div class="flex flex-col gap-1.5 overflow-y-auto p-3">
                {#each response.cookies as cookie, i (i)}
                    <div
                        class="flex flex-wrap items-center gap-1.5 rounded-sm px-1 py-1 hover:bg-cave-raised"
                    >
                        <input
                            class="min-w-0 flex-1"
                            placeholder="Name"
                            value={cookie.name}
                            oninput={(e) => updateCookie(i, { name: e.currentTarget.value })}
                        />
                        <input
                            class="min-w-0 flex-1"
                            placeholder="Value"
                            value={cookie.value}
                            oninput={(e) => updateCookie(i, { value: e.currentTarget.value })}
                        />
                        <input
                            class="w-20"
                            placeholder="Path"
                            value={cookie.path}
                            oninput={(e) => updateCookie(i, { path: e.currentTarget.value })}
                        />
                        <label
                            class="flex cursor-pointer items-center gap-1.5 whitespace-nowrap text-xs text-ink-mid"
                        >
                            <input
                                type="checkbox"
                                checked={cookie.httpOnly}
                                onchange={(e) =>
                                    updateCookie(i, { httpOnly: e.currentTarget.checked })}
                            />
                            HttpOnly
                        </label>
                        <label
                            class="flex cursor-pointer items-center gap-1.5 whitespace-nowrap text-xs text-ink-mid"
                        >
                            <input
                                type="checkbox"
                                checked={cookie.secure}
                                onchange={(e) =>
                                    updateCookie(i, { secure: e.currentTarget.checked })}
                            />
                            Secure
                        </label>
                        <Button variant="icon" onclick={() => removeCookie(i)}>✕</Button>
                    </div>
                {/each}
                <Button size="sm" onclick={addCookie}>+ Add Cookie</Button>
            </div>
        {/if}
    </div>
</div>
