<script lang="ts">
  import type { View } from '../stores/index'
  import { activeView, pendingBreakpoints } from '../stores/index'

  const items: { id: View; label: string }[] = [
    { id: 'servers', label: 'Servers' },
    { id: 'collections', label: 'Collections' },
    { id: 'breakpoints', label: 'Breakpoints' },
    { id: 'settings', label: 'Settings' },
  ]
</script>

<nav class="flex w-14 shrink-0 flex-col items-center gap-1 border-r border-wire bg-cave-deep py-3">
  <div class="mb-3 text-2xl">🐟</div>

  {#each items as item}
    <button
      class="relative flex h-11 w-11 flex-col items-center justify-center rounded-lg transition-all
        {$activeView === item.id
        ? 'bg-cave-surface text-accent'
        : 'text-ink-muted hover:bg-cave-raised hover:text-ink'}"
      on:click={() => activeView.set(item.id)}
      title={item.label}
    >
      {#if item.id === 'servers'}
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
          <rect x="2" y="2" width="20" height="8" rx="2" />
          <rect x="2" y="14" width="20" height="8" rx="2" />
          <circle cx="6" cy="6" r="1" fill="currentColor" stroke="none" />
          <circle cx="6" cy="18" r="1" fill="currentColor" stroke="none" />
        </svg>
      {:else if item.id === 'collections'}
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" />
        </svg>
      {:else if item.id === 'breakpoints'}
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
          <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
        </svg>
        {#if $pendingBreakpoints.length > 0}
          <span
            class="absolute right-1.5 top-1.5 flex h-4 w-4 items-center justify-center rounded-full bg-accent text-[9px] font-bold text-white"
          >
            {$pendingBreakpoints.length}
          </span>
        {/if}
      {:else if item.id === 'settings'}
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="3" />
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06A1.65 1.65 0 0 0 19.4 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z" />
        </svg>
      {/if}
      {#if $activeView === item.id}
        <span class="absolute left-0 top-1/2 h-5 w-0.5 -translate-y-1/2 rounded-r bg-accent" />
      {/if}
    </button>
  {/each}
</nav>
