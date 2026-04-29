<script lang="ts">
  import { METHOD_COLORS } from '../lib/types'
  import {
    collections,
    expandedCollections,
    expandedFolders,
    patchCollection,
    selectedEndpointId,
    selectedFolderId,
    selectedCollectionId,
  } from '../stores/index'
  import * as api from '../lib/api'

  let newCollectionName = ''
  let showNewCollection = false

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
    selectedCollectionId.set(collectionId)
    selectedFolderId.set(folderId)
    selectedEndpointId.set(endpointId)
  }

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

  async function createFolder(collectionId: string) {
    const name = prompt('Folder name:')
    if (!name?.trim()) return
    const updated = await api.createFolder(collectionId, name.trim())
    patchCollection(updated)
    expandedFolders.update((s) => {
      const folder = updated.folders[updated.folders.length - 1]
      s.add(folder.id)
      return new Set(s)
    })
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

  async function deleteCollection(id: string) {
    if (!confirm('Delete this collection?')) return
    await api.deleteCollection(id)
    collections.update((list) => list.filter((c) => c.id !== id))
    selectedCollectionId.update((v) => (v === id ? null : v))
  }

  async function deleteFolder(collectionId: string, folderId: string) {
    if (!confirm('Delete this folder?')) return
    const updated = await api.deleteFolder(collectionId, folderId)
    patchCollection(updated)
    selectedFolderId.update((v) => (v === folderId ? null : v))
  }

  async function deleteEndpoint(collectionId: string, folderId: string, endpointId: string) {
    if (!confirm('Delete this endpoint?')) return
    const updated = await api.deleteEndpoint(collectionId, folderId, endpointId)
    patchCollection(updated)
    selectedEndpointId.update((v) => (v === endpointId ? null : v))
  }

  function mc(m: string) {
    return METHOD_COLORS[m.toUpperCase()] || '#64748b'
  }
</script>

<aside class="flex w-60 shrink-0 flex-col overflow-hidden border-r border-wire bg-cave-deep">
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
              on:click={() => createFolder(collection.id)}
            >
              <svg width="12" height="12" viewBox="0 0 16 16" fill="currentColor"
                ><path
                  d="M1 3.5A1.5 1.5 0 012.5 2h3.764c.958 0 1.76.56 2.177 1.368l.105.21c.175.35.53.572.93.572H13.5A1.5 1.5 0 0115 5.5v7A1.5 1.5 0 0113.5 14h-11A1.5 1.5 0 011 12.5v-9z"
                /></svg
              >
            </button>
            <button
              class="btn-icon text-[11px] hover:text-err"
              title="Delete"
              on:click={() => deleteCollection(collection.id)}>✕</button
            >
          </div>
        </div>

        {#if $expandedCollections.has(collection.id)}
          {#each collection.folders as folder}
            <div>
              <!-- Folder row -->
              <div
                class="group flex min-h-[28px] items-center gap-0.5 py-1 pl-6 pr-2 text-ink-mid hover:bg-cave-raised hover:text-ink"
              >
                <button
                  class="flex h-4 w-4 shrink-0 items-center justify-center bg-transparent text-[10px] text-ink-muted"
                  on:click={() => toggleFolder(folder.id)}
                >
                  {$expandedFolders.has(folder.id) ? '▾' : '▸'}
                </button>
                <span class="shrink-0 text-xs">{$expandedFolders.has(folder.id) ? '📂' : '📁'}</span
                >
                <span class="flex-1 truncate text-[13px] font-medium">{folder.name}</span>
                <div class="flex gap-0.5 opacity-0 group-hover:opacity-100">
                  <button
                    class="btn-icon text-[11px]"
                    title="Add endpoint"
                    on:click={() => createEndpoint(collection.id, folder.id)}>+</button
                  >
                  <button
                    class="btn-icon text-[11px] hover:text-err"
                    title="Delete folder"
                    on:click={() => deleteFolder(collection.id, folder.id)}>✕</button
                  >
                </div>
              </div>

              {#if $expandedFolders.has(folder.id)}
                {#each folder.endpoints as ep}
                  <button
                    class="group flex min-h-[28px] w-full items-center gap-1.5 py-1 pl-10 pr-2 text-left transition-colors
                      {$selectedEndpointId === ep.id
                      ? 'bg-accent/10 text-ink'
                      : 'text-ink-mid hover:bg-cave-raised hover:text-ink'}"
                    on:click={() => selectEndpoint(collection.id, folder.id, ep.id)}
                  >
                    <span
                      class="method-badge tag shrink-0 text-center"
                      style="background:{mc(ep.method)};min-width:36px"
                      >{ep.method === '*' ? 'ANY' : ep.method}</span
                    >
                    <span class="flex-1 truncate font-mono text-[11px] text-ink-muted"
                      >{ep.path || '/...'}</span
                    >
                    {#if ep.breakpoint}<span class="shrink-0 text-[10px]" title="Breakpoint"
                        >⚡</span
                      >{/if}
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
                on:click={() => createFolder(collection.id)}>+ Add folder</button
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
          on:click={() => (showNewCollection = true)}
        >
          + New Collection
        </button>
      </div>
    {/if}
  </div>
</aside>
