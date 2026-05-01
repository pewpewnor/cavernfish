<script lang="ts" generics="T extends string">
    import type { Snippet } from "svelte";

    interface TabItem {
        value: T;
        label: string;
        badge?: string | number | null;
    }

    interface Props {
        value: T;
        tabs: TabItem[];
        onChange?: (v: T) => void;
        trailing?: Snippet;
    }

    let { value = $bindable(), tabs, onChange, trailing }: Props = $props();

    function pick(v: T) {
        value = v;
        onChange?.(v);
    }
</script>

<div class="flex shrink-0 items-center border-b border-wire bg-cave-deep px-3">
    {#each tabs as t (t.value)}
        <button
            class="-mb-px border-b-2 bg-transparent px-3 py-2.5 text-sm transition-colors {value ===
            t.value
                ? 'border-accent text-accent'
                : 'border-transparent text-ink-mid hover:text-ink'}"
            onclick={() => pick(t.value)}
        >
            {t.label}{t.badge !== undefined && t.badge !== null && t.badge !== 0 && t.badge !== ""
                ? ` (${t.badge})`
                : ""}
        </button>
    {/each}
    {#if trailing}
        <div class="ml-auto flex items-center">{@render trailing()}</div>
    {/if}
</div>
