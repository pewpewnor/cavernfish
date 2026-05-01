<script lang="ts">
    import type { KVPair } from "../lib/types";
    import { newKVPair } from "../lib/types";
    import { Button } from "../lib/ui";

    interface Props {
        pairs: KVPair[];
        keyPlaceholder?: string;
        valuePlaceholder?: string;
        extras?: KVPair[];
        extrasLabel?: string;
        onChange: (pairs: KVPair[]) => void;
    }

    let {
        pairs,
        keyPlaceholder = "Key",
        valuePlaceholder = "Value",
        extras = [],
        extrasLabel = "Inferred",
        onChange,
    }: Props = $props();

    function add() {
        onChange([...pairs, newKVPair()]);
    }

    function remove(i: number) {
        onChange(pairs.filter((_, idx) => idx !== i));
    }

    function toggle(i: number) {
        onChange(pairs.map((p, idx) => (idx === i ? { ...p, enabled: !p.enabled } : p)));
    }

    function setKey(i: number, v: string) {
        onChange(pairs.map((p, idx) => (idx === i ? { ...p, key: v } : p)));
    }

    function setValue(i: number, v: string) {
        onChange(pairs.map((p, idx) => (idx === i ? { ...p, value: v } : p)));
    }
</script>

<div class="flex flex-col gap-0.5">
    {#if extras.length > 0 || pairs.length > 0}
        <div
            class="grid grid-cols-[20px_1fr_1fr_28px] gap-1 px-1 pb-1 text-[0.625rem] font-semibold uppercase tracking-[0.06em] text-ink-muted"
        >
            <span></span>
            <span>{keyPlaceholder}</span>
            <span>{valuePlaceholder}</span>
            <span></span>
        </div>
    {/if}

    {#if extras.length > 0}
        <div class="mb-1 flex items-center gap-1.5 px-1 pb-0.5 pt-1">
            <span
                class="rounded-[3px] bg-cave-elevated px-1.5 py-0.5 text-[0.5625rem] font-semibold uppercase tracking-[0.06em] text-ink-muted"
                >{extrasLabel}</span
            >
            <span class="text-[0.625rem] italic text-ink-muted"
                >read-only — set the same key below to override</span
            >
        </div>
        {#each extras as pair, i (i)}
            <div
                class="grid grid-cols-[20px_1fr_1fr_28px] items-center gap-1 rounded-sm px-1 py-0.5 opacity-60"
            >
                <span
                    class="flex h-5 w-5 items-center justify-center text-[0.625rem] text-ink-muted"
                    title="Inferred — not editable"
                >
                    ⓘ
                </span>
                <input class="cursor-not-allowed italic" readonly value={pair.key} />
                <input class="cursor-not-allowed italic" readonly value={pair.value} />
                <span></span>
            </div>
        {/each}
    {/if}

    {#each pairs as pair, i (i)}
        <div
            class="grid grid-cols-[20px_1fr_1fr_28px] items-center gap-1 rounded-sm px-1 py-0.5 transition-colors hover:bg-cave-raised {pair.enabled
                ? ''
                : 'opacity-45'}"
        >
            <button
                class="flex h-5 w-5 items-center justify-center bg-transparent"
                onclick={() => toggle(i)}
                title={pair.enabled ? "Disable" : "Enable"}
                aria-label={pair.enabled ? "Disable pair" : "Enable pair"}
            >
                <span
                    class="block h-2 w-2 rounded-full transition-colors {pair.enabled
                        ? 'bg-accent'
                        : 'bg-ink-muted'}"
                ></span>
            </button>
            <input
                placeholder={keyPlaceholder}
                value={pair.key}
                oninput={(e) => setKey(i, e.currentTarget.value)}
            />
            <input
                placeholder={valuePlaceholder}
                value={pair.value}
                oninput={(e) => setValue(i, e.currentTarget.value)}
            />
            <Button variant="icon" onclick={() => remove(i)}>✕</Button>
        </div>
    {/each}

    <div class="mt-1">
        <Button size="sm" onclick={add}>+ Add</Button>
    </div>
</div>
