<script lang="ts">
  import { onMount } from 'svelte'
  import TopBar from './components/TopBar.svelte'
  import Sidebar from './components/Sidebar.svelte'
  import ServerPanel from './components/ServerPanel.svelte'
  import CollectionsView from './components/CollectionsView.svelte'
  import SettingsPanel from './components/SettingsPanel.svelte'
  import BreakpointPanel from './components/BreakpointPanel.svelte'
  import { activeView, collections, servers, pendingBreakpoints } from './stores/index'
  import * as api from './lib/api'

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

<div class="flex h-full flex-col overflow-hidden">
  <TopBar />

  <div class="flex min-h-0 flex-1 overflow-hidden">
    {#if $activeView === 'collections'}
      <Sidebar />
      <CollectionsView />
    {:else if $activeView === 'servers'}
      <ServerPanel />
    {:else if $activeView === 'settings'}
      <SettingsPanel />
    {/if}
  </div>

  {#if $pendingBreakpoints.length > 0}
    <BreakpointPanel />
  {/if}
</div>
