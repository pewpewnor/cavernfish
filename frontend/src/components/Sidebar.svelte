<script lang="ts">
  import { METHOD_COLORS } from '../lib/types'
  import {
    collections,
    expandedCollections,
    expandedFolders,
    patchCollection,
    openTabs,
    activeTabId,
  } from '../stores/index'
  import type { EndpointTab } from '../stores/index'
  import * as api from '../lib/api'
  import Modal from './Modal.svelte'
  import ResizeHandle from './ResizeHandle.svelte'

  let sidebarWidth = 240

  // ── New collection inline form ───────────────────────────────────
  let newCollectionName = ''
  let showNewCollection = false

  // ── New folder modal (replaces native prompt) ────────────────────
  let newFolderModal: { collectionId: string } | null = null
  let newFolderName = ''

  // ── Context menu ─────────────────────────────────────────────────
  let ctxMenu: {
    type: 'collection' | 'folder' | 'endpoint'
    collectionId: string
    folderId?: string
    endpointId?: string
    x: number
    y: number
  } | null = null

  // ── Rename modal ─────────────────────────────────────────────────
  let renameModal: {
    type: 'collection' | 'folder' | 'endpoint'
    collectionId: string
    folderId?: string
    endpointId?: string
    currentName: string
  } | null = null
  let renameValue = ''

  // ── Expand / select ──────────────────────────────────────────────
  function toggleCollection(id: string) {
    expandedCollections.update((s) => {
      const next = new Set(s)
      next.has(id) ? next.delete(id) : next.add(id)
      return next
    })
  }

  function toggleFolder(id: string) {
    expandedFolders.update((s) => {
      const next = new Set(s)
      next.has(id) ? next.delete(id) : next.add(id)
      return next
    })
  }

  function selectEndpoint(collectionId: string, folderId: string, endpointId: string) {
    const existing = $openTabs.find((t) => t.endpointId === endpointId)
    if (existing) {
      activeTabId.set(existing.tabId)
      return
    }
    const tab: EndpointTab = { tabId: endpointId, endpointId, collectionId, folderId }
    openTabs.update((tabs) => [...tabs, tab])
    activeTabId.set(endpointId)
  }

  // ── Create ───────────────────────────────────────────────────────
  async function createCollection() {
    if (!newCollectionName.trim()) return
    const c = await api.createCollection(newCollectionName.trim(), '')
    collections.update((list) => [...list, c])
    expandedCollections.update((s) => {
      s.add(c.id)
      return new Set(s)
    })
    newCollectionName = ''
    showNewCollection = false
  }

  function openNewFolderModal(collectionId: string) {
    newFolderModal = { collectionId }
    newFolderName = ''
  }

  async function createFolder() {
    if (!newFolderModal || !newFolderName.trim()) return
    const updated = await api.createFolder(newFolderModal.collectionId, newFolderName.trim())
    patchCollection(updated)
    expandedFolders.update((s) => {
      const folder = updated.folders[updated.folders.length - 1]
      s.add(folder.id)
      return new Set(s)
    })
    newFolderModal = null
    newFolderName = ''
  }

  async function createEndpoint(collectionId: string, folderId: string) {
    const updated = await api.createEndpoint(collectionId, folderId, {
      id: crypto.randomUUID(),
      name: 'New Endpoint',
      method: 'GET',
      path: '/api/new',
      responses: [],
      strategy: 'fixed',
      activeIdx: 0,
      delayMs: 0,
      breakpoint: false,
      wsEnabled: false,
    })
    patchCollection(updated)
    const folder = updated.folders.find((f) => f.id === folderId)
    if (folder) {
      const ep = folder.endpoints[folder.endpoints.length - 1]
      selectEndpoint(collectionId, folderId, ep.id)
    }
  }

  // ── Delete ───────────────────────────────────────────────────────
  async function deleteCollection(id: string) {
    if (!confirm('Delete this collection?')) return
    await api.deleteCollection(id)
    collections.update((list) => list.filter((c) => c.id !== id))
  }

  async function deleteFolder(collectionId: string, folderId: string) {
    if (!confirm('Delete this folder?')) return
    const updated = await api.deleteFolder(collectionId, folderId)
    patchCollection(updated)
  }

  async function deleteEndpoint(collectionId: string, folderId: string, endpointId: string) {
    if (!confirm('Delete this endpoint?')) return
    const updated = await api.deleteEndpoint(collectionId, folderId, endpointId)
    patchCollection(updated)
    openTabs.update((tabs) => tabs.filter((t) => t.endpointId !== endpointId))
  }

  // ── Context menu ─────────────────────────────────────────────────
  function openCtxMenu(
    e: MouseEvent,
    type: 'collection' | 'folder' | 'endpoint',
    collectionId: string,
    folderId?: string,
    endpointId?: string,
  ) {
    e.preventDefault()
    e.stopPropagation()
    ctxMenu = { type, collectionId, folderId, endpointId, x: e.clientX, y: e.clientY }
  }

  function closeCtxMenu() {
    ctxMenu = null
  }

  function startRename() {
    if (!ctxMenu) return
    const { type, collectionId, folderId, endpointId } = ctxMenu
    let currentName = ''
    if (type === 'collection') {
      currentName = $collections.find((c) => c.id === collectionId)?.name ?? ''
    } else if (type === 'folder') {
      const col = $collections.find((c) => c.id === collectionId)
      currentName = col?.folders.find((f) => f.id === folderId)?.name ?? ''
    } else {
      const col = $collections.find((c) => c.id === collectionId)
      const folder = col?.folders.find((f) => f.id === folderId)
      currentName = folder?.endpoints.find((e) => e.id === endpointId)?.name ?? ''
    }
    renameModal = { type, collectionId, folderId, endpointId, currentName }
    renameValue = currentName
    ctxMenu = null
  }

  async function doRename() {
    if (!renameModal || !renameValue.trim()) return
    const { type, collectionId, folderId, endpointId } = renameModal
    if (type === 'collection') {
      const col = $collections.find((c) => c.id === collectionId)
      if (!col) return
      const updated = await api.updateCollection(collectionId, {
        name: renameValue.trim(),
        description: col.description,
        variables: col.variables ?? {},
      })
      collections.update((list) => list.map((c) => (c.id === collectionId ? updated : c)))
    } else if (type === 'folder' && folderId) {
      const updated = await api.renameFolder(collectionId, folderId, renameValue.trim())
      patchCollection(updated)
    } else if (type === 'endpoint' && folderId && endpointId) {
      const col = $collections.find((c) => c.id === collectionId)
      const folder = col?.folders.find((f) => f.id === folderId)
      const ep = folder?.endpoints.find((e) => e.id === endpointId)
      if (!ep) return
      const updated = await api.updateEndpoint(collectionId, folderId, endpointId, {
        ...ep,
        name: renameValue.trim(),
      })
      patchCollection(updated)
    }
    renameModal = null
  }

  function mc(m: string) {
    return METHOD_COLORS[m.toUpperCase()] || '#64748b'
  }
</script>

<svelte:window on:click={closeCtxMenu} />

<aside
  class="relative flex shrink-0 flex-col overflow-hidden border-r border-wire bg-cave-deep"
  style="width:{sidebarWidth}px"
>
  <ResizeHandle currentWidth={sidebarWidth} min={160} max={520} on:resize={(e) => (sidebarWidth = e.detail)} />
  <!-- Header -->
  <div class="flex shrink-0 items-center justify-between border-b border-wire px-3 py-2.5">
    <span class="text-[11px] font-semibold uppercase tracking-[0.08em] text-ink-muted"
      >Collections</span
    >
    <button
      class="btn-icon text-base leading-none hover:text-accent"
      on:click={() => (showNewCollection = !showNewCollection)}
      title="New Collection">+</button
    >
  </div>

  {#if showNewCollection}
    <div class="flex shrink-0 flex-col gap-1.5 border-b border-wire bg-cave-surface px-2.5 py-2">
      <input
        bind:value={newCollectionName}
        placeholder="Collection name"
        on:keydown={(e) => {
          if (e.key === 'Enter') createCollection()
          if (e.key === 'Escape') showNewCollection = false
        }}
        autofocus
      />
      <div class="flex gap-1">
        <button class="btn btn-primary px-2.5 py-0.5 text-[11px]" on:click={createCollection}
          >Create</button
        >
        <button
          class="btn btn-ghost px-2 py-0.5 text-[11px]"
          on:click={() => (showNewCollection = false)}>Cancel</button
        >
      </div>
    </div>
  {/if}

  <!-- Tree -->
  <div class="flex-1 overflow-y-auto py-1">
    {#each $collections as collection}
      <div>
        <!-- Collection row -->
        <div
          class="group flex min-h-[28px] items-center gap-0.5 px-2 py-1 text-ink-mid hover:bg-cave-raised hover:text-ink"
          on:contextmenu={(e) => openCtxMenu(e, 'collection', collection.id)}
        >
          <button
            class="flex h-4 w-4 shrink-0 items-center justify-center bg-transparent text-[10px] text-ink-muted"
            on:click={() => toggleCollection(collection.id)}
          >
            {$expandedCollections.has(collection.id) ? '▾' : '▸'}
          </button>
          <span class="shrink-0 text-xs text-accent">⬡</span>
          <span class="flex-1 truncate text-[13px] font-medium">{collection.name}</span>
          <div class="flex gap-0.5 opacity-0 group-hover:opacity-100">
            <button
              class="btn-icon text-[11px]"
              title="Add folder"
              on:click|stopPropagation={() => openNewFolderModal(collection.id)}
            >
              <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor">
                <path
                  d="M1 3.5A1.5 1.5 0 012.5 2h3.764c.958 0 1.76.56 2.177 1.368l.105.21c.175.35.53.572.93.572H13.5A1.5 1.5 0 0115 5.5v7A1.5 1.5 0 0113.5 14h-11A1.5 1.5 0 011 12.5v-9z"
                />
              </svg>
            </button>
            <button
              class="btn-icon text-[11px] hover:text-err"
              title="Delete"
              on:click|stopPropagation={() => deleteCollection(collection.id)}>✕</button
            >
          </div>
        </div>

        {#if $expandedCollections.has(collection.id)}
          {#each collection.folders as folder}
            <div>
              <!-- Folder row -->
              <div
                class="group flex min-h-[28px] items-center gap-0.5 py-1 pl-6 pr-2 text-ink-mid hover:bg-cave-raised hover:text-ink"
                on:contextmenu={(e) => openCtxMenu(e, 'folder', collection.id, folder.id)}
              >
                <button
                  class="flex h-4 w-4 shrink-0 items-center justify-center bg-transparent text-[10px] text-ink-muted"
                  on:click={() => toggleFolder(folder.id)}
                >
                  {$expandedFolders.has(folder.id) ? '▾' : '▸'}
                </button>
                <span class="shrink-0 text-xs"
                  >{$expandedFolders.has(folder.id) ? '📂' : '📁'}</span
                >
                <span class="flex-1 truncate text-[13px] font-medium">{folder.name}</span>
                <div class="flex gap-0.5 opacity-0 group-hover:opacity-100">
                  <button
                    class="btn-icon text-[11px]"
                    title="Add endpoint"
                    on:click|stopPropagation={() => createEndpoint(collection.id, folder.id)}
                    >+</button
                  >
                  <button
                    class="btn-icon text-[11px] hover:text-err"
                    title="Delete folder"
                    on:click|stopPropagation={() => deleteFolder(collection.id, folder.id)}
                    >✕</button
                  >
                </div>
              </div>

              {#if $expandedFolders.has(folder.id)}
                {#each folder.endpoints as ep}
                  <button
                    class="group flex w-full min-h-[28px] items-center gap-1.5 py-1 pl-10 pr-2 text-left transition-colors
                      {$activeTabId === ep.id
                      ? 'bg-accent/10 text-ink'
                      : $openTabs.some((t) => t.endpointId === ep.id)
                        ? 'text-ink hover:bg-cave-raised'
                        : 'text-ink-mid hover:bg-cave-raised hover:text-ink'}"
                    on:click={() => selectEndpoint(collection.id, folder.id, ep.id)}
                    on:contextmenu|preventDefault|stopPropagation={(e) =>
                      openCtxMenu(e, 'endpoint', collection.id, folder.id, ep.id)}
                  >
                    <span
                      class="method-badge tag shrink-0 text-center"
                      style="background:{mc(ep.method)};min-width:34px"
                      >{ep.method === '*' ? 'ANY' : ep.method}</span
                    >
                    <span class="flex-1 truncate text-sm">{ep.name || ep.path || 'Unnamed'}</span>
                    {#if ep.breakpoint}
                      <span class="shrink-0 text-[10px]" title="Breakpoint">⚡</span>
                    {/if}
                    <button
                      class="btn-icon shrink-0 text-[10px] opacity-0 hover:text-err group-hover:opacity-100"
                      title="Delete"
                      on:click|stopPropagation={() =>
                        deleteEndpoint(collection.id, folder.id, ep.id)}>✕</button
                    >
                  </button>
                {/each}
                {#if folder.endpoints.length === 0}
                  <div class="py-1 pl-12 text-[11px] italic text-ink-muted">No endpoints</div>
                {/if}
              {/if}
            </div>
          {/each}

          {#if collection.folders.length === 0}
            <div class="flex flex-col items-start gap-1.5 py-2 pl-8 pr-2">
              <span class="text-xs text-ink-muted">No folders yet</span>
              <button
                class="btn btn-ghost px-2.5 py-0.5 text-[11px]"
                on:click={() => openNewFolderModal(collection.id)}>+ Add folder</button
              >
            </div>
          {/if}
        {/if}
      </div>
    {/each}

    {#if $collections.length === 0}
      <div class="flex flex-col items-center gap-2.5 px-5 py-10 text-center">
        <span class="text-3xl opacity-20">⬡</span>
        <p class="text-[13px] text-ink-muted">No collections yet.</p>
        <button
          class="btn btn-primary px-3.5 py-1.5 text-xs"
          on:click={() => (showNewCollection = true)}>+ New Collection</button
        >
      </div>
    {/if}
  </div>
</aside>

<!-- ── Context menu ──────────────────────────────────────────────── -->
{#if ctxMenu}
  <div
    class="fixed z-[2000] min-w-[140px] overflow-hidden rounded border border-wire bg-cave-elevated shadow-modal"
    style="left:{ctxMenu.x}px;top:{ctxMenu.y}px"
    on:click|stopPropagation
  >
    <button
      class="flex w-full items-center px-3 py-2 text-left text-xs text-ink-mid hover:bg-cave-raised hover:text-ink"
      on:click={startRename}>Rename</button
    >
    <div class="border-t border-wire" />
    <button
      class="flex w-full items-center px-3 py-2 text-left text-xs text-err hover:bg-err/10"
      on:click={() => {
        const m = ctxMenu
        ctxMenu = null
        if (!m) return
        if (m.type === 'collection') deleteCollection(m.collectionId)
        else if (m.type === 'folder' && m.folderId) deleteFolder(m.collectionId, m.folderId)
        else if (m.type === 'endpoint' && m.folderId && m.endpointId)
          deleteEndpoint(m.collectionId, m.folderId, m.endpointId)
      }}>Delete</button
    >
  </div>
{/if}

<!-- ── Rename modal ───────────────────────────────────────────────── -->
{#if renameModal}
  <Modal title="Rename" on:close={() => (renameModal = null)}>
    <div class="flex flex-col gap-3">
      <input
        bind:value={renameValue}
        autofocus
        on:keydown={(e) => {
          if (e.key === 'Enter') doRename()
          if (e.key === 'Escape') renameModal = null
        }}
      />
      <div class="flex justify-end gap-2">
        <button class="btn btn-primary" on:click={doRename}>Rename</button>
        <button class="btn btn-ghost" on:click={() => (renameModal = null)}>Cancel</button>
      </div>
    </div>
  </Modal>
{/if}

<!-- ── New folder modal ───────────────────────────────────────────── -->
{#if newFolderModal}
  <Modal
    title="New Folder"
    on:close={() => {
      newFolderModal = null
      newFolderName = ''
    }}
  >
    <div class="flex flex-col gap-3">
      <input
        bind:value={newFolderName}
        placeholder="Folder name"
        autofocus
        on:keydown={(e) => {
          if (e.key === 'Enter') createFolder()
          if (e.key === 'Escape') {
            newFolderModal = null
            newFolderName = ''
          }
        }}
      />
      <div class="flex justify-end gap-2">
        <button class="btn btn-primary" on:click={createFolder}>Create</button>
        <button
          class="btn btn-ghost"
          on:click={() => {
            newFolderModal = null
            newFolderName = ''
          }}>Cancel</button
        >
      </div>
    </div>
  </Modal>
{/if}
