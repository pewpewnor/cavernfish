<script lang="ts">
    interface Props {
        min?: number;
        max?: number;
        currentWidth: number;
        onResize: (w: number) => void;
    }

    let { min = 160, max = 520, currentWidth, onResize }: Props = $props();

    let active = $state(false);

    function onMousedown(e: MouseEvent) {
        const startX = e.clientX;
        const startW = currentWidth;
        active = true;
        e.preventDefault();

        function onMove(ev: MouseEvent) {
            onResize(Math.max(min, Math.min(max, startW + ev.clientX - startX)));
        }

        function onUp() {
            active = false;
            window.removeEventListener("mousemove", onMove);
            window.removeEventListener("mouseup", onUp);
        }

        window.addEventListener("mousemove", onMove);
        window.addEventListener("mouseup", onUp);
    }
</script>

<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<div
    role="separator"
    aria-orientation="vertical"
    tabindex="-1"
    class="absolute right-0 top-0 z-10 h-full w-1.5 cursor-col-resize select-none transition-colors {active
        ? 'bg-accent/50'
        : 'hover:bg-accent/30'}"
    onmousedown={onMousedown}
></div>
