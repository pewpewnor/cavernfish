<script lang="ts">
    import { METHOD_COLORS } from "../types";

    interface Props {
        method: string;
        size?: "sm" | "md";
    }

    let { method, size = "md" }: Props = $props();

    let color = $derived(METHOD_COLORS[method.toUpperCase()] ?? "#64748b");
    let cls = $derived(
        size === "sm" ? "text-[0.5625rem] px-1 py-0" : "text-[0.625rem] px-1.5 py-px",
    );
</script>

<span
    class="inline-block shrink-0 rounded-[3px] font-mono font-bold uppercase tracking-wider {cls}"
    style="background:{color};color:#1a1f30"
>
    {method === "*" ? "ANY" : method}
</span>
