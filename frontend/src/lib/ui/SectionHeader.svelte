<script lang="ts">
    import type { Snippet } from "svelte";

    interface Props {
        title: string;
        children?: Snippet;
    }

    let { title, children }: Props = $props();
</script>

<div class="flex shrink-0 items-center justify-between border-b border-wire px-3 py-2.5">
    <span class="text-[0.6875rem] font-semibold uppercase tracking-[0.08em] text-ink-muted">
        {title}
    </span>
    {#if children}
        {@render children()}
    {/if}
</div>
