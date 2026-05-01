<script lang="ts">
    import type { Snippet } from "svelte";

    interface Props {
        active?: boolean;
        onclick?: (e: MouseEvent) => void;
        oncontextmenu?: (e: MouseEvent) => void;
        indent?: 0 | 1 | 2 | 3;
        children: Snippet;
        actions?: Snippet;
        leading?: Snippet;
        title?: string;
    }

    let {
        active = false,
        onclick,
        oncontextmenu,
        indent = 0,
        children,
        actions,
        leading,
        title,
    }: Props = $props();

    const padX: Record<number, string> = {
        0: "pl-3 pr-2",
        1: "pl-6 pr-2",
        2: "pl-10 pr-2",
        3: "pl-14 pr-2",
    };
</script>

<button
    class="group flex w-full items-center gap-2 bg-transparent py-2 text-left transition-colors {padX[
        indent
    ]} {active ? 'bg-accent/10 text-ink' : 'text-ink-mid hover:bg-cave-raised hover:text-ink'}"
    {onclick}
    {oncontextmenu}
    {title}
>
    {#if leading}
        {@render leading()}
    {/if}
    <div class="min-w-0 flex-1">
        {@render children()}
    </div>
    {#if actions}
        <div class="flex shrink-0 gap-0.5 opacity-0 group-hover:opacity-100">
            {@render actions()}
        </div>
    {/if}
</button>
