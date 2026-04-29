<script lang="ts">
  import type { KVPair } from '../lib/types'
  import { newKVPair } from '../lib/types'

  export let pairs: KVPair[] = []
  export let keyPlaceholder = 'Key'
  export let valuePlaceholder = 'Value'

  function add() {
    pairs = [...pairs, newKVPair()]
  }
  function remove(i: number) {
    pairs = pairs.filter((_, idx) => idx !== i)
  }
  function toggle(i: number) {
    pairs = pairs.map((p, idx) => (idx === i ? { ...p, enabled: !p.enabled } : p))
  }
</script>

<div class="flex flex-col gap-0.5">
  {#if pairs.length > 0}
    <div
      class="grid grid-cols-[20px_1fr_1fr_28px] gap-1 px-1 pb-1 text-[10px] font-semibold uppercase tracking-[0.06em] text-ink-muted"
    >
      <span />
      <span>{keyPlaceholder}</span>
      <span>{valuePlaceholder}</span>
      <span />
    </div>
    {#each pairs as pair, i}
      <div
        class="grid grid-cols-[20px_1fr_1fr_28px] items-center gap-1 rounded-sm px-1 py-0.5 transition-colors hover:bg-cave-raised {pair.enabled
          ? ''
          : 'opacity-45'}"
      >
        <button
          class="flex h-5 w-5 items-center justify-center bg-transparent"
          on:click={() => toggle(i)}
          title={pair.enabled ? 'Disable' : 'Enable'}
        >
          <span
            class="block h-2 w-2 rounded-full transition-colors {pair.enabled
              ? 'bg-accent'
              : 'bg-ink-muted'}"
          />
        </button>
        <input
          class="w-full border-transparent bg-transparent px-1.5 py-1 font-mono text-xs hover:border-wire hover:bg-cave-raised focus:border-accent focus:bg-cave-raised"
          placeholder={keyPlaceholder}
          value={pair.key}
          on:input={(e) => {
            const v = e.currentTarget.value
            pairs = pairs.map((p, idx) => (idx === i ? { ...p, key: v } : p))
          }}
        />
        <input
          class="w-full border-transparent bg-transparent px-1.5 py-1 font-mono text-xs hover:border-wire hover:bg-cave-raised focus:border-accent focus:bg-cave-raised"
          placeholder={valuePlaceholder}
          value={pair.value}
          on:input={(e) => {
            const v = e.currentTarget.value
            pairs = pairs.map((p, idx) => (idx === i ? { ...p, value: v } : p))
          }}
        />
        <button
          class="btn-icon text-[10px] hover:text-err"
          on:click={() => remove(i)}
          title="Remove">✕</button
        >
      </div>
    {/each}
  {/if}
  <button class="btn btn-ghost mt-1 self-start px-2.5 py-1 text-[11px]" on:click={add}>
    + Add
  </button>
</div>
