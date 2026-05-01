<script lang="ts">
    import { Button, Modal } from "../lib/ui";
    import { servers } from "../stores/index";
    import * as api from "../lib/api";
    import type { ServerInfo, Collection, Folder, Endpoint } from "../lib/types";

    type Target =
        | { type: "collection"; collection: Collection }
        | { type: "folder"; collection: Collection; folder: Folder }
        | { type: "endpoint"; collection: Collection; folder: Folder; endpoint: Endpoint };

    interface Props {
        target: Target;
        onClose: () => void;
    }

    let { target, onClose }: Props = $props();

    function endpointIds(t: Target): string[] {
        if (t.type === "collection")
            return t.collection.folders.flatMap((f) => f.endpoints.map((e) => e.id));
        if (t.type === "folder") return t.folder.endpoints.map((e) => e.id);
        return [t.endpoint.id];
    }

    function isAssigned(s: ServerInfo): boolean {
        const colId = target.collection.id;
        if (!s.collectionIds.includes(colId)) return false;
        const eps = endpointIds(target);
        if (eps.length === 0) return true;
        return eps.every((id) => !s.disabledEndpointIds.includes(id));
    }

    function isPartial(s: ServerInfo): boolean {
        const colId = target.collection.id;
        if (!s.collectionIds.includes(colId)) return false;
        const eps = endpointIds(target);
        if (eps.length === 0) return false;
        const enabled = eps.filter((id) => !s.disabledEndpointIds.includes(id)).length;
        return enabled > 0 && enabled < eps.length;
    }

    async function toggleServer(s: ServerInfo) {
        const colId = target.collection.id;
        const eps = endpointIds(target);
        const turningOn = !isAssigned(s);
        let collectionIds = [...s.collectionIds];
        let disabled = [...s.disabledEndpointIds];
        if (turningOn) {
            if (!collectionIds.includes(colId)) collectionIds.push(colId);
            disabled = disabled.filter((id) => !eps.includes(id));
        } else if (target.type === "collection") {
            collectionIds = collectionIds.filter((id) => id !== colId);
            disabled = disabled.filter((id) => !eps.includes(id));
        } else {
            disabled = [...new Set([...disabled, ...eps])];
        }
        const updated = await api.updateServer({
            id: s.id,
            name: s.name,
            port: s.port,
            collectionIds,
            https: s.https,
            disabledEndpointIds: disabled,
        });
        servers.update((list) => list.map((srv) => (srv.id === updated.id ? updated : srv)));
    }

    let title = $derived(
        target.type === "collection"
            ? `Assign servers — ${target.collection.name}`
            : target.type === "folder"
              ? `Assign servers — ${target.folder.name}`
              : `Assign servers — ${target.endpoint.name || target.endpoint.path}`,
    );
</script>

<Modal open={true} {title} width="480px" {onClose}>
    <div class="flex flex-col gap-2">
        {#if $servers.length === 0}
            <p class="py-4 text-center text-sm text-ink-muted">
                No servers yet. Create one in the Servers panel.
            </p>
        {:else}
            {#each $servers as s (s.id)}
                {@const assigned = isAssigned(s)}
                {@const partial = isPartial(s)}
                <label
                    class="flex cursor-pointer items-center gap-3 rounded-sm border border-wire bg-cave-surface px-3 py-2 hover:bg-cave-raised"
                >
                    <input
                        type="checkbox"
                        checked={assigned}
                        indeterminate={partial}
                        onchange={() => toggleServer(s)}
                    />
                    <div class="flex min-w-0 flex-1 flex-col">
                        <span class="truncate text-[0.8125rem] font-medium text-ink"
                            >{s.name || "Unnamed"}</span
                        >
                        <span class="font-mono text-[0.6875rem] text-ink-muted"
                            >:{s.port} · {s.status}</span
                        >
                    </div>
                    {#if partial}
                        <span class="text-[0.6875rem] italic text-ink-muted">partial</span>
                    {/if}
                </label>
            {/each}
        {/if}
        <div class="flex justify-end gap-2 pt-2">
            <Button variant="primary" onclick={onClose}>Done</Button>
        </div>
    </div>
</Modal>
