<script lang="ts">
    import type { Snippet } from "svelte";
    import { ResizeHandle } from "../lib/ui";
    import { pxToRem, remToPx } from "../lib/zoom";
    import { workspaceSidebarWidthRem } from "../stores/index";

    interface Props {
        title: string;
        actions?: Snippet;
        children: Snippet;
        min?: number;
        max?: number;
    }

    let { title, actions, children, min = 10, max = 32 }: Props = $props();
</script>

<aside
    class="relative flex shrink-0 flex-col overflow-hidden border-r border-wire bg-cave-deep"
    style="width:{$workspaceSidebarWidthRem}rem"
>
    <ResizeHandle
        currentWidth={remToPx($workspaceSidebarWidthRem)}
        min={remToPx(min)}
        max={remToPx(max)}
        onResize={(w) => workspaceSidebarWidthRem.set(pxToRem(w))}
    />
    <div class="flex shrink-0 items-center justify-between border-b border-wire px-3 py-2.5">
        <span class="text-[0.6875rem] font-semibold uppercase tracking-[0.08em] text-ink-muted">
            {title}
        </span>
        {#if actions}
            {@render actions()}
        {/if}
    </div>
    <div class="flex flex-1 flex-col overflow-hidden">
        {@render children()}
    </div>
</aside>
