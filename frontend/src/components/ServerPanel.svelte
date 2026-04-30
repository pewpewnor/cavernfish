<script lang="ts">
  import { onMount } from 'svelte'
  import type { ServerInfo, ServerConfig } from '../lib/types'
  import { METHOD_COLORS, statusColor } from '../lib/types'
  import { servers, collections, requestLogs, addLog } from '../stores/index'
  import * as api from '../lib/api'
  import ResizeHandle from './ResizeHandle.svelte'

  let listWidth = 224
  let selectedId: string | null = null
  let editForm: ServerConfig | null = null
  let dirty = false
  let detailTab: 'serve' | 'logs' = 'serve'
  let logSearch = ''

  $: selectedServer = $servers.find((s) => s.id === selectedId) ?? null

  $: filteredLogs = ($requestLogs[selectedId ?? ''] ?? []).filter(
    (e) =>
      !logSearch ||
      e.path.toLowerCase().includes(logSearch.toLowerCase()) ||
      e.method.toLowerCase().includes(logSearch.toLowerCase()) ||
      String(e.statusCode).includes(logSearch),
  )

  function selectServer(id: string | null) {
    selectedId = id
    const srv = $servers.find((s) => s.id === id)
    editForm = srv
      ? {
          id: srv.id,
          name: srv.name,
          port: srv.port,
          collectionIds: [...(srv.collectionIds ?? [])],
          https: srv.https,
          disabledEndpointIds: [...(srv.disabledEndpointIds ?? [])],
        }
      : null
    dirty = false
    detailTab = 'serve'
    logSearch = ''
    if (id) {
      api.getRequestLog(id).then((log) => {
        requestLogs.update((r) => ({ ...r, [id]: log }))
      })
    }
  }

  function markDirty() {
    dirty = true
  }

  async function createServer() {
    const cfg: ServerConfig = {
      id: crypto.randomUUID(),
      name: 'New Server',
      port: 8000,
      collectionIds: [],
      https: false,
      disabledEndpointIds: [],
    }
    const created = await api.createServer(cfg)
    servers.update((list) => [...list, created])
    selectServer(created.id)
  }

  async function saveServer() {
    if (!editForm) return
    if (!editForm.port || editForm.port < 1 || editForm.port > 65535) {
      alert('Invalid port')
      return
    }
    const updated = await api.updateServer(editForm)
    servers.update((list) => list.map((s) => (s.id === updated.id ? updated : s)))
    dirty = false
  }

  async function toggleServer(s: ServerInfo) {
    if (s.status === 'running') {
      await api.stopServer(s.id)
    } else {
      await api.startServer(s.id)
    }
  }

  async function deleteServer(id: string) {
    if (!confirm('Delete this server?')) return
    await api.deleteServer(id)
    servers.update((list) => list.filter((s) => s.id !== id))
    if (selectedId === id) selectServer(null)
  }

  function isCollectionEnabled(colId: string): boolean {
    return (editForm?.collectionIds ?? []).includes(colId)
  }

  function toggleCollection(colId: string) {
    if (!editForm) return
    const ids = editForm.collectionIds
    if (ids.includes(colId)) {
      const col = $collections.find((c) => c.id === colId)
      const colEpIds = col ? col.folders.flatMap((f) => f.endpoints.map((e) => e.id)) : []
      editForm = {
        ...editForm,
        collectionIds: ids.filter((id) => id !== colId),
        disabledEndpointIds: editForm.disabledEndpointIds.filter((id) => !colEpIds.includes(id)),
      }
    } else {
      editForm = { ...editForm, collectionIds: [...ids, colId] }
    }
    markDirty()
  }

  function getFolderState(folderEpIds: string[]): 'all' | 'some' | 'none' {
    if (!editForm || folderEpIds.length === 0) return 'all'
    const disabledCount = folderEpIds.filter((id) =>
      (editForm?.disabledEndpointIds ?? []).includes(id),
    ).length
    if (disabledCount === 0) return 'all'
    if (disabledCount === folderEpIds.length) return 'none'
    return 'some'
  }

  function toggleFolder(folderEpIds: string[], state: 'all' | 'some' | 'none') {
    if (!editForm) return
    let disabled = [...editForm.disabledEndpointIds]
    if (state === 'none') {
      disabled = disabled.filter((id) => !folderEpIds.includes(id))
    } else {
      disabled = [...new Set([...disabled, ...folderEpIds])]
    }
    editForm = { ...editForm, disabledEndpointIds: disabled }
    markDirty()
  }

  function isEndpointEnabled(epId: string): boolean {
    return !(editForm?.disabledEndpointIds ?? []).includes(epId)
  }

  function toggleEndpoint(epId: string) {
    if (!editForm) return
    const disabled = editForm.disabledEndpointIds
    if (disabled.includes(epId)) {
      editForm = { ...editForm, disabledEndpointIds: disabled.filter((id) => id !== epId) }
    } else {
      editForm = { ...editForm, disabledEndpointIds: [...disabled, epId] }
    }
    markDirty()
  }

  async function clearLog(id: string) {
    await api.clearRequestLog(id)
    requestLogs.update((r) => ({ ...r, [id]: [] }))
  }

  function formatTime(ts: string): string {
    return new Date(ts).toLocaleTimeString()
  }

  function mc(m: string): string {
    return METHOD_COLORS[m.toUpperCase()] ?? '#64748b'
  }

  onMount(() => {
    api.onRequestLogged((entry) => addLog(entry))
    api.onServerStatus((updated) => servers.set(updated))
  })
</script>

<div class="flex min-w-0 flex-1 overflow-hidden">
  <aside
    class="relative flex shrink-0 flex-col border-r border-wire bg-cave-deep"
    style="width:{listWidth}px"
  >
    <ResizeHandle
      currentWidth={listWidth}
      min={160}
      max={400}
      on:resize={(e) => (listWidth = e.detail)}
    />

    <div class="flex shrink-0 items-center justify-between border-b border-wire px-3 py-3">
      <span class="text-[11px] font-semibold uppercase tracking-[0.08em] text-ink-muted"
        >Servers</span
      >
      <button
        class="btn-icon text-lg leading-none hover:text-accent"
        on:click={createServer}
        title="New Server">+</button
      >
    </div>

    <div class="flex-1 overflow-y-auto py-1">
      {#each $servers as srv}
        <button
          class="flex w-full items-center gap-3 px-3 py-3 text-left transition-colors
            {selectedId === srv.id
            ? 'bg-accent/10 text-ink'
            : 'text-ink-mid hover:bg-cave-raised hover:text-ink'}"
          on:click={() => selectServer(srv.id)}
        >
          <span
            class="h-2 w-2 shrink-0 rounded-full
              {srv.status === 'running'
              ? 'bg-ok'
              : srv.status === 'error'
                ? 'bg-err'
                : 'bg-ink-muted'}"
          />
          <div class="min-w-0 flex-1">
            <div class="truncate font-medium">{srv.name || 'Unnamed'}</div>
            <div class="flex items-center gap-1.5 font-mono text-xs text-ink-muted">
              <span>:{srv.port}</span>
              {#if srv.https}<span class="text-accent">HTTPS</span>{/if}
            </div>
          </div>
        </button>
      {/each}

      {#if $servers.length === 0}
        <div class="flex flex-col items-center gap-2.5 px-4 py-10 text-center">
          <p class="text-xs text-ink-muted">No servers yet</p>
          <button class="btn btn-primary" on:click={createServer}>+ New Server</button>
        </div>
      {/if}
    </div>
  </aside>

  <div class="flex min-w-0 flex-1 flex-col overflow-hidden">
    {#if editForm && selectedServer}
      <div class="flex shrink-0 flex-col gap-3 border-b border-wire bg-cave-surface px-5 py-4">
        <div class="flex items-center gap-3">
          <input
            class="flex-1 text-base font-semibold"
            bind:value={editForm.name}
            placeholder="Server name"
            on:input={markDirty}
          />
          <div class="flex shrink-0 items-center gap-2">
            {#if dirty}
              <span class="h-2.5 w-2.5 rounded-full bg-warn" title="Unsaved changes" />
            {/if}
            <button class="btn btn-ghost" on:click={saveServer} disabled={!dirty}>
              {dirty ? 'Save' : 'Saved ✓'}
            </button>
            <button
              class="btn {selectedServer.status === 'running' ? 'btn-ghost' : 'btn-primary'} min-w-[80px]"
              on:click={() => selectedServer && toggleServer(selectedServer)}
              disabled={selectedServer.status === 'error'}
            >
              {selectedServer.status === 'running' ? '⏹ Stop' : '▶ Start'}
            </button>
            <button class="btn-icon hover:text-err" on:click={() => deleteServer(selectedServer.id)}
              >🗑</button
            >
          </div>
        </div>

        <div class="flex flex-wrap items-end gap-5">
          <div class="flex flex-col gap-1">
            <span class="form-label">Port</span>
            <input
              type="number"
              min="1"
              max="65535"
              class="w-24 font-mono"
              bind:value={editForm.port}
              on:input={markDirty}
            />
          </div>
          <label class="flex cursor-pointer items-center gap-2 pb-1">
            <input type="checkbox" bind:checked={editForm.https} on:change={markDirty} />
            <span class="text-sm text-ink-mid">HTTPS</span>
            {#if editForm.https}
              <span class="font-mono text-xs text-ink-muted"
                >(HTTP also on :{editForm.port + 1})</span
              >
            {/if}
          </label>
          <div class="flex items-center gap-2 pb-1">
            <span class="status-pill status-{selectedServer.status}">{selectedServer.status}</span>
            {#if selectedServer.requestCount > 0}
              <span class="font-mono text-xs text-ink-muted">{selectedServer.requestCount} req</span>
            {/if}
            {#if selectedServer.errorMsg}
              <span
                class="max-w-[240px] truncate text-xs text-err"
                title={selectedServer.errorMsg}>⚠ {selectedServer.errorMsg}</span
              >
            {/if}
          </div>
        </div>
      </div>

      <div class="tab-strip">
        <button
          class="tab-btn {detailTab === 'serve' ? 'active' : ''}"
          on:click={() => (detailTab = 'serve')}>Serve</button
        >
        <button
          class="tab-btn {detailTab === 'logs' ? 'active' : ''}"
          on:click={() => {
            detailTab = 'logs'
            if (selectedId) {
              const sid = selectedId
              api.getRequestLog(sid).then((log) => {
                requestLogs.update((r) => ({ ...r, [sid]: log }))
              })
            }
          }}>Logs{selectedServer.requestCount > 0 ? ` (${selectedServer.requestCount})` : ''}</button
        >
      </div>

      <div class="flex min-h-0 flex-1 flex-col overflow-hidden">
        {#if detailTab === 'serve'}
          <div class="flex-1 overflow-y-auto p-3">
            {#if $collections.length === 0}
              <p class="py-6 text-center text-sm text-ink-muted">
                No collections. Create one in the Collections tab.
              </p>
            {:else}
              <div class="flex flex-col gap-2">
                {#each $collections as col}
                  {@const colEnabled = isCollectionEnabled(col.id)}
                  <div class="overflow-hidden rounded-lg border border-wire">
                    <label
                      class="flex cursor-pointer items-center gap-2.5 bg-cave-surface px-3 py-2.5 hover:bg-cave-raised"
                    >
                      <input
                        type="checkbox"
                        checked={colEnabled}
                        on:change={() => toggleCollection(col.id)}
                      />
                      <span class="text-sm text-accent">⬡</span>
                      <span class="flex-1 font-medium text-ink">{col.name}</span>
                      {#if col.description}
                        <span class="max-w-[200px] truncate text-xs text-ink-muted"
                          >{col.description}</span
                        >
                      {/if}
                    </label>

                    {#if colEnabled && col.folders.length > 0}
                      <div class="bg-cave-deep pb-1 pt-0.5">
                        {#each col.folders as folder}
                          {@const folderEpIds = folder.endpoints.map((e) => e.id)}
                          {@const folderState = getFolderState(folderEpIds)}
                          <div class="px-3">
                            <label class="flex cursor-pointer items-center gap-2 py-1.5">
                              <input
                                type="checkbox"
                                checked={folderState !== 'none'}
                                on:change={() => toggleFolder(folderEpIds, folderState)}
                              />
                              <span class="text-xs"
                                >{folderState !== 'none' ? '📂' : '📁'}</span
                              >
                              <span class="text-sm text-ink-mid">{folder.name}</span>
                              {#if folderState === 'some'}
                                <span class="ml-1 text-[11px] italic text-ink-muted">partial</span>
                              {/if}
                            </label>
                            {#each folder.endpoints as ep}
                              <label class="flex cursor-pointer items-center gap-2 py-1 pl-5">
                                <input
                                  type="checkbox"
                                  checked={isEndpointEnabled(ep.id)}
                                  on:change={() => toggleEndpoint(ep.id)}
                                />
                                <span
                                  class="tag method-badge shrink-0"
                                  style="background:{mc(ep.method)}"
                                  >{ep.method === '*' ? 'ANY' : ep.method}</span
                                >
                                <span class="min-w-0 flex-1 truncate font-mono text-xs text-ink-mid"
                                  >{ep.path}</span
                                >
                                <span class="max-w-[120px] truncate text-xs text-ink-muted"
                                  >{ep.name}</span
                                >
                              </label>
                            {/each}
                          </div>
                        {/each}
                      </div>
                    {/if}
                  </div>
                {/each}
              </div>
            {/if}
          </div>
        {:else}
          <div class="flex shrink-0 items-center gap-2 border-b border-wire px-3 py-2">
            <input
              class="flex-1"
              placeholder="Search by path, method, or status…"
              bind:value={logSearch}
            />
            <button
              class="btn btn-ghost px-2.5 py-1 text-xs"
              on:click={() => clearLog(selectedServer.id)}>Clear</button
            >
          </div>
          <div class="flex-1 overflow-y-auto">
            {#if filteredLogs.length === 0}
              <div class="py-10 text-center text-sm text-ink-muted">
                {logSearch ? 'No matching entries' : 'No requests yet'}
              </div>
            {:else}
              {#each [...filteredLogs].reverse() as entry}
                <div
                  class="flex items-center gap-2 border-b border-wire/40 px-3 py-2 text-xs
                    {!entry.matched ? 'opacity-50' : ''}"
                >
                  <span class="w-16 shrink-0 font-mono text-ink-muted"
                    >{formatTime(entry.timestamp)}</span
                  >
                  <span
                    class="tag method-badge shrink-0"
                    style="background:{mc(entry.method)};font-size:9px;padding:1px 5px"
                    >{entry.method}</span
                  >
                  <span class="min-w-0 flex-1 truncate font-mono text-ink-mid"
                    >{entry.path}{entry.query ? '?' + entry.query : ''}</span
                  >
                  <span
                    class="shrink-0 font-mono font-bold"
                    style="color:{statusColor(entry.statusCode)}">{entry.statusCode}</span
                  >
                  <span class="w-12 shrink-0 text-right font-mono text-ink-muted"
                    >{entry.latencyMs}ms</span
                  >
                  {#if entry.isProxy}
                    <span
                      class="shrink-0 rounded-[3px] bg-cave-elevated px-1 py-0.5 font-mono text-[9px] text-ink-muted"
                      >proxy</span
                    >
                  {/if}
                  {#if entry.isWs}
                    <span
                      class="shrink-0 rounded-[3px] bg-cave-elevated px-1 py-0.5 font-mono text-[9px] text-ink-muted"
                      >ws</span
                    >
                  {/if}
                  {#if entry.breakpointed}
                    <span
                      class="shrink-0 rounded-[3px] bg-accent/[.12] px-1 py-0.5 font-mono text-[9px] text-accent"
                      >⚡bp</span
                    >
                  {/if}
                </div>
              {/each}
            {/if}
          </div>
        {/if}
      </div>
    {:else}
      <div class="flex flex-1 flex-col items-center justify-center gap-3 p-10 text-center">
        <div class="text-5xl opacity-[0.12]">◎</div>
        <p class="text-sm text-ink-mid">Select a server or create a new one</p>
        <button class="btn btn-primary" on:click={createServer}>+ New Server</button>
      </div>
    {/if}
  </div>
</div>
