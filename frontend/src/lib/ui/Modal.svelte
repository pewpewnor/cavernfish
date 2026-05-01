<script lang="ts">
    import { Dialog } from "bits-ui";
    import type { Snippet } from "svelte";

    interface Props {
        open: boolean;
        title?: string;
        width?: string;
        onClose?: () => void;
        children?: Snippet;
        footer?: Snippet;
    }

    let { open = $bindable(), title, width = "440px", onClose, children, footer }: Props = $props();

    function handleOpenChange(o: boolean) {
        open = o;
        if (!o) onClose?.();
    }
</script>

<Dialog.Root {open} onOpenChange={handleOpenChange}>
    <Dialog.Portal>
        <Dialog.Overlay class="fixed inset-0 z-[1000] bg-cave-base/75 backdrop-blur-sm" />
        <Dialog.Content
            class="fixed left-1/2 top-1/2 z-[1001] flex max-h-[90vh] -translate-x-1/2 -translate-y-1/2 flex-col overflow-hidden rounded-lg border border-wire-hi bg-cave-elevated shadow-modal focus:outline-none"
            style="width: {width}"
        >
            {#if title}
                <div class="flex items-center justify-between border-b border-wire px-4 py-3">
                    <Dialog.Title class="text-sm font-semibold text-ink">{title}</Dialog.Title>
                    <Dialog.Close
                        class="rounded-sm bg-transparent p-1 text-ink-muted transition-colors hover:bg-cave-raised hover:text-ink"
                    >
                        ✕
                    </Dialog.Close>
                </div>
            {/if}
            <div class="flex-1 overflow-y-auto p-4">
                {@render children?.()}
            </div>
            {#if footer}
                <div
                    class="flex shrink-0 items-center justify-end gap-2 border-t border-wire px-4 py-3"
                >
                    {@render footer()}
                </div>
            {/if}
        </Dialog.Content>
    </Dialog.Portal>
</Dialog.Root>
