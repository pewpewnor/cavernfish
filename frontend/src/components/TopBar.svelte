<script lang="ts">
  import type { View } from '../stores/index'
  import { activeView, pendingBreakpoints } from '../stores/index'

  const views: { id: View; label: string; icon: string }[] = [
    { id: 'servers', label: 'Servers', icon: '◎' },
    { id: 'collections', label: 'Collections', icon: '⬡' },
    { id: 'settings', label: 'Settings', icon: '⚙' },
  ]
</script>

<header
  class="flex h-11 shrink-0 items-center gap-0 border-b border-wire bg-cave-deep px-4"
  style="-webkit-app-region: drag"
>
  <div class="mr-5 flex items-center gap-1.5" style="-webkit-app-region: no-drag">
    <span class="text-lg">🐟</span>
    <span class="text-sm font-bold tracking-tight text-ink">Cavernfish</span>
  </div>

  <nav class="flex gap-0.5" style="-webkit-app-region: no-drag">
    {#each views as v}
      <button
        class="flex items-center gap-1.5 rounded px-3 py-1.5 text-[13px] font-medium transition-all
          {$activeView === v.id
          ? 'bg-cave-surface text-ink'
          : 'bg-transparent text-ink-mid hover:bg-cave-raised hover:text-ink'}"
        on:click={() => activeView.set(v.id)}
      >
        <span class="text-sm {$activeView === v.id ? 'text-accent' : 'text-ink-muted'}"
          >{v.icon}</span
        >
        <span>{v.label}</span>
      </button>
    {/each}
  </nav>

  <div class="ml-auto" style="-webkit-app-region: no-drag">
    {#if $pendingBreakpoints.length > 0}
      <div
        class="flex items-center gap-1.5 rounded-full border border-accent bg-accent/10 px-2.5 py-1"
      >
        <span class="block h-2 w-2 animate-pulse-slow rounded-full bg-accent" />
        <span class="text-xs font-semibold text-accent">
          {$pendingBreakpoints.length} breakpoint{$pendingBreakpoints.length !== 1 ? 's' : ''} waiting
        </span>
      </div>
    {/if}
  </div>
</header>
