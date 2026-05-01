<script lang="ts">
    import { Switch } from "bits-ui";

    interface Props {
        checked: boolean;
        label?: string;
        disabled?: boolean;
    }

    let { checked = $bindable(), label, disabled = false }: Props = $props();
</script>

<label class="inline-flex cursor-pointer items-center gap-2.5 {disabled ? 'opacity-50' : ''}">
    {#if label}
        <span class="text-[0.6875rem] font-semibold uppercase tracking-[0.06em] text-ink-muted">
            {label}
        </span>
    {/if}
    <Switch.Root
        bind:checked
        {disabled}
        class="relative inline-flex h-[18px] w-8 shrink-0 cursor-pointer items-center rounded-full border border-wire bg-cave-raised transition-colors disabled:cursor-not-allowed data-[state=checked]:border-accent data-[state=checked]:bg-accent/20"
    >
        <Switch.Thumb
            class="block h-3 w-3 translate-x-0.5 rounded-full bg-ink-muted transition-transform data-[state=checked]:translate-x-[18px] data-[state=checked]:bg-accent"
        />
    </Switch.Root>
</label>
