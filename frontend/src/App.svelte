<script lang="ts">
  import { onMount } from 'svelte'
  import NavRail from './components/NavRail.svelte'
  import Sidebar from './components/Sidebar.svelte'
  import ServerPanel from './components/ServerPanel.svelte'
  import CollectionsView from './components/CollectionsView.svelte'
  import SettingsPanel from './components/SettingsPanel.svelte'
  import BreakpointConfigPanel from './components/BreakpointConfigPanel.svelte'
  import { activeView, collections, servers, pendingBreakpoints } from './stores/index'
  import * as api from './lib/api'

  let zoom = 1

  function handleKeydown(e: KeyboardEvent) {
    if (!e.ctrlKey) return
    if (e.key === '=' || e.key === '+') {
      e.preventDefault()
      zoom = Math.min(+(zoom + 0.1).toFixed(1), 2)
    } else if (e.key === '-') {
      e.preventDefault()
      zoom = Math.max(+(zoom - 0.1).toFixed(1), 0.5)
    } else if (e.key === '0') {
      e.preventDefault()
      zoom = 1
    }
  }

  $: document.documentElement.style.zoom = String(zoom)

  onMount(async () => {
    const [cols, srvs] = await Promise.all([api.getCollections(), api.getServers()])
    collections.set(cols ?? [])
    servers.set(srvs ?? [])
    api.onServerStatus((updated) => servers.set(updated))
    api.onBreakpointHit((hit) => {
      pendingBreakpoints.update((list) => {
        if (list.some((b) => b.id === hit.id)) return list
        return [...list, hit]
      })
    })
    api.onBreakpointReleased((id) => {
      pendingBreakpoints.update((list) => list.filter((b) => b.id !== id))
    })
  })
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="flex h-full overflow-hidden">
  <NavRail />
  <div class="flex min-w-0 flex-1 overflow-hidden">
    {#if $activeView === 'collections'}
      <Sidebar />
      <CollectionsView />
    {:else if $activeView === 'servers'}
      <ServerPanel />
    {:else if $activeView === 'breakpoints'}
      <BreakpointConfigPanel />
    {:else if $activeView === 'settings'}
      <SettingsPanel />
    {/if}
  </div>
</div>
