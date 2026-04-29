<script lang="ts">
  import { createEventDispatcher } from 'svelte'

  export let title: string = ''
  export let width: string = '440px'

  const dispatch = createEventDispatcher()

  function onKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') dispatch('close')
  }
</script>

<svelte:window on:keydown={onKeydown} />

<div
  class="fixed inset-0 z-[1000] flex items-center justify-center bg-cave-base/75 backdrop-blur-sm"
  on:click|self={() => dispatch('close')}
>
  <div
    class="flex max-h-[90vh] flex-col overflow-hidden rounded-lg border border-wire-hi bg-cave-elevated shadow-modal"
    style="width: {width}"
  >
    {#if title}
      <div class="flex items-center justify-between border-b border-wire px-4 py-3">
        <span class="text-sm font-semibold text-ink">{title}</span>
        <button class="btn-icon text-xs" on:click={() => dispatch('close')}>✕</button>
      </div>
    {/if}
    <div class="flex-1 overflow-y-auto p-4">
      <slot />
    </div>
  </div>
</div>
