<script lang="ts">
  import type { Endpoint } from '../lib/types'
  import { METHOD_COLORS } from '../lib/types'
  import {
    collections,
    openTabs,
    activeTabId,
    dirtyTabs,
    tabEdits,
    patchCollection,
  } from '../stores/index'
  import * as api from '../lib/api'
  import EndpointEditor from './EndpointEditor.svelte'
  import Modal from './Modal.svelte'

  let closeConfirmTabId: string | null = null

  function getEndpointForTab(tabId: string): Endpoint | null {
    const tab = $openTabs.find((t) => t.tabId === tabId)
    if (!tab) return null
    const col = $collections.find((c) => c.id === tab.collectionId)
    if (!col) return null
    const folder = col.folders.find((f) => f.id === tab.folderId)
    if (!folder) return null
    return folder.endpoints.find((e) => e.id === tab.endpointId) ?? null
  }

  function requestCloseTab(tabId: string) {
    if ($dirtyTabs.has(tabId)) {
      closeConfirmTabId = tabId
    } else {
      doCloseTab(tabId)
    }
  }

  function doCloseTab(tabId: string) {
    const tabs = $openTabs
    const idx = tabs.findIndex((t) => t.tabId === tabId)
    openTabs.update((list) => list.filter((t) => t.tabId !== tabId))
    dirtyTabs.update((s) => {
      s.delete(tabId)
      return new Set(s)
    })
    tabEdits.update((m) => {
      const next = { ...m }
      delete next[tabId]
      return next
    })
    if ($activeTabId === tabId) {
      const remaining = $openTabs
      if (remaining.length > 0) {
        const nextIdx = Math.min(idx, remaining.length - 1)
        activeTabId.set(remaining[nextIdx >= 0 ? nextIdx : 0].tabId)
      } else {
        activeTabId.set(null)
      }
    }
    closeConfirmTabId = null
  }

  async function saveAndClose(tabId: string) {
    const tab = $openTabs.find((t) => t.tabId === tabId)
    const ep = $tabEdits[tabId]
    if (tab && ep) {
      const updated = await api.updateEndpoint(tab.collectionId, tab.folderId, ep.id, ep)
      patchCollection(updated)
    }
    doCloseTab(tabId)
  }

  async function handleSave(
    e: CustomEvent<{ collectionId: string; folderId: string; endpoint: Endpoint }>,
  ) {
    const { collectionId, folderId, endpoint } = e.detail
    const updated = await api.updateEndpoint(collectionId, folderId, endpoint.id, endpoint)
    patchCollection(updated)
  }

  function handleDirtyChange(tabId: string, dirty: boolean) {
    dirtyTabs.update((s) => {
      if (dirty) s.add(tabId)
      else s.delete(tabId)
      return new Set(s)
    })
  }

  function mc(m: string) {
    return METHOD_COLORS[m.toUpperCase()] || '#64748b'
  }
</script>

<div class="flex min-w-0 flex-1 flex-col overflow-hidden">
  {#if $openTabs.length > 0}
    <div
      class="flex shrink-0 items-stretch overflow-x-auto border-b border-wire bg-cave-deep"
      style="min-height:36px"
    >
      {#each $openTabs as tab}
        {@const ep = getEndpointForTab(tab.tabId)}
        <div
          class="group flex shrink-0 items-stretch border-r border-wire
            {$activeTabId === tab.tabId ? 'bg-cave-surface' : 'hover:bg-cave-raised'}"
        >
          <button
            class="flex items-center gap-1.5 px-3 py-2 text-xs transition-colors
              {$activeTabId === tab.tabId ? 'text-ink' : 'text-ink-mid hover:text-ink'}"
            on:click={() => activeTabId.set(tab.tabId)}
          >
            {#if ep}
              <span
                class="tag method-badge shrink-0"
                style="background:{mc(ep.method)};font-size:9px;padding:1px 5px"
                >{ep.method === '*' ? 'ANY' : ep.method}</span
              >
              <span class="max-w-[140px] truncate">{ep.name || ep.path || 'Unnamed'}</span>
            {:else}
              <span class="text-ink-muted italic">Deleted</span>
            {/if}
            {#if $dirtyTabs.has(tab.tabId)}
              <span class="h-1.5 w-1.5 shrink-0 rounded-full bg-warn" />
            {/if}
          </button>
          <button
            class="px-2 text-ink-muted opacity-0 transition-all hover:text-ink group-hover:opacity-100
              {$activeTabId === tab.tabId ? 'opacity-100' : ''}"
            on:click={() => requestCloseTab(tab.tabId)}>✕</button
          >
        </div>
      {/each}
    </div>

    {#each $openTabs as tab}
      {@const ep = getEndpointForTab(tab.tabId)}
      <div
        class="flex min-h-0 flex-1 flex-col overflow-hidden"
        class:hidden={$activeTabId !== tab.tabId}
      >
        {#if ep}
          <EndpointEditor
            tabId={tab.tabId}
            endpoint={ep}
            collectionId={tab.collectionId}
            folderId={tab.folderId}
            on:save={handleSave}
            on:dirtyChange={(e) => handleDirtyChange(tab.tabId, e.detail)}
          />
        {:else}
          <div class="flex flex-1 items-center justify-center">
            <p class="text-sm text-ink-muted">This endpoint has been deleted.</p>
          </div>
        {/if}
      </div>
    {/each}
  {:else}
    <div class="flex flex-1 flex-col items-center justify-center gap-3 p-10 text-center">
      <div class="relative mb-2 flex items-center justify-center">
        <span class="relative z-10 text-[52px] opacity-20">⬡</span>
        <div class="absolute h-24 w-24 animate-breathe rounded-full bg-accent/10" />
      </div>
      <h3 class="text-base font-semibold text-ink">Select an endpoint</h3>
      <p class="max-w-sm text-sm leading-relaxed text-ink-muted">
        Click an endpoint in the sidebar to open it in a tab.
      </p>
    </div>
  {/if}
</div>

{#if closeConfirmTabId}
  {@const confirmTabEp = getEndpointForTab(closeConfirmTabId)}
  <Modal title="Unsaved Changes" on:close={() => (closeConfirmTabId = null)}>
    <div class="flex flex-col gap-4">
      <p class="text-sm text-ink-mid">
        <span class="font-medium text-ink">{confirmTabEp?.name || confirmTabEp?.path || 'This endpoint'}</span> has unsaved
        changes. What would you like to do?
      </p>
      <div class="flex justify-end gap-2">
        <button
          class="btn btn-primary"
          on:click={() => {
            if (closeConfirmTabId) saveAndClose(closeConfirmTabId)
          }}>Save & Close</button
        >
        <button
          class="btn btn-ghost"
          on:click={() => {
            if (closeConfirmTabId) doCloseTab(closeConfirmTabId)
          }}>Discard & Close</button
        >
        <button class="btn btn-ghost" on:click={() => (closeConfirmTabId = null)}>Cancel</button>
      </div>
    </div>
  </Modal>
{/if}
