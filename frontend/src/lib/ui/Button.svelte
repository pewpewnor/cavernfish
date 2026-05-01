<script lang="ts">
    import type { Snippet } from "svelte";
    import type { HTMLButtonAttributes } from "svelte/elements";

    type Variant = "primary" | "ghost" | "danger" | "icon";
    type Size = "sm" | "md";

    interface Props extends Omit<HTMLButtonAttributes, "children" | "class"> {
        variant?: Variant;
        size?: Size;
        class?: string;
        children?: Snippet;
    }

    // svelte-ignore custom_element_props_identifier
    let {
        variant = "ghost",
        size = "md",
        class: className = "",
        children,
        ...rest
    }: Props = $props();

    const variantClass: Record<Variant, string> = {
        primary: "bg-accent text-white hover:bg-accent-hi",
        ghost: "border border-wire bg-transparent text-ink-mid hover:border-wire-hi hover:bg-cave-raised hover:text-ink",
        danger: "border border-transparent bg-transparent text-err hover:bg-err/10",
        icon: "rounded-sm bg-transparent p-1 text-ink-muted hover:bg-cave-raised hover:text-ink",
    };

    const sizeClass: Record<Size, string> = {
        sm: "px-2.5 py-1 text-[0.6875rem]",
        md: "px-3.5 py-1.5 text-sm",
    };
</script>

<button
    class="inline-flex select-none items-center gap-1.5 rounded-sm font-medium transition-all disabled:cursor-not-allowed disabled:opacity-50 {variantClass[
        variant
    ]} {variant !== 'icon' ? sizeClass[size] : ''} {className}"
    {...rest}
>
    {@render children?.()}
</button>
