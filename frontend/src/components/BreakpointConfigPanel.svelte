<script lang="ts">
  import type { MockResponse, Endpoint, Folder, Collection } from '../lib/types'
  import { METHOD_COLORS, statusColor } from '../lib/types'
  import { collections, pendingBreakpoints, patchCollection } from '../stores/index'
  import * as api from '../lib/api'
  import ResponseEditor from './ResponseEditor.svelte'
  import JsonEditor from './JsonEditor.svelte'

  let panelTab: 'endpoints' | 'queue' = 'endpoints'

  let expandedCols = new Set<string>()
  let expandedFolders = new Set<string>()

  function toggleColExpand(id: string) {
    const next = new Set(expandedCols)
    next.has(id) ? next.delete(id) : next.add(id)
    expandedCols = next
  }

  function toggleFolderExpand(id: string) {
    const next = new Set(expandedFolders)
    next.has(id) ? next.delete(id) : next.add(id)
    expandedFolders = next
  }

  function getColBpState(col: Collection): 'all' | 'some' | 'none' {
    const eps = col.folders.flatMap((f) => f.endpoints)
    if (eps.length === 0) return 'none'
    const n = eps.filter((e) => e.breakpoint).length
    if (n === 0) return 'none'
    if (n === eps.length) return 'all'
    return 'some'
  }

  function getFolderBpState(folder: Folder): 'all' | 'some' | 'none' {
    if (folder.endpoints.length === 0) return 'none'
    const n = folder.endpoints.filter((e) => e.breakpoint).length
    if (n === 0) return 'none'
    if (n === folder.endpoints.length) return 'all'
    return 'some'
  }

  async function toggleEpBp(colId: string, folderId: string, ep: Endpoint) {
    const updated = await api.updateEndpoint(colId, folderId, ep.id, {
      ...ep,
      breakpoint: !ep.breakpoint,
    })
    patchCollection(updated)
  }

  async function toggleFolderBp(
    colId: string,
    folder: Folder,
    state: 'all' | 'some' | 'none',
  ) {
    const newBp = state === 'none'
    for (const ep of folder.endpoints) {
      if (ep.breakpoint !== newBp) {
        const updated = await api.updateEndpoint(colId, folder.id, ep.id, {
          ...ep,
          breakpoint: newBp,
        })
        patchCollection(updated)
      }
    }
  }

  async function toggleColBp(col: Collection, state: 'all' | 'some' | 'none') {
    const newBp = state === 'none'
    for (const folder of col.folders) {
      for (const ep of folder.endpoints) {
        if (ep.breakpoint !== newBp) {
          const updated = await api.updateEndpoint(col.id, folder.id, ep.id, {
            ...ep,
            breakpoint: newBp,
          })
          patchCollection(updated)
        }
      }
    }
  }

  let activeQueueIdx = 0
  $: activeHit = $pendingBreakpoints[activeQueueIdx] ?? null

  let editingResponse: MockResponse | null = null
  let prevHitId = ''
  let respTab: 'body' | 'headers' | 'cookies' = 'body'
  let reqTab: 'body' | 'headers' = 'body'

  $: if (activeHit?.id !== prevHitId) {
    prevHitId = activeHit?.id ?? ''
    editingResponse = activeHit
      ? (JSON.parse(JSON.stringify(activeHit.response)) as MockResponse)
      : null
    respTab = 'body'
    reqTab = 'body'
  }

  function handleResponseUpdate(e: CustomEvent<MockResponse>) {
    editingResponse = e.detail
  }

  async function release() {
    if (!activeHit || !editingResponse) return
    const id = activeHit.id
    await api.releaseBreakpoint(id, editingResponse)
    pendingBreakpoints.update((list) => list.filter((b) => b.id !== id))
    activeQueueIdx = Math.min(activeQueueIdx, $pendingBreakpoints.length - 1)
  }

  async function discard() {
    if (!activeHit) return
    const id = activeHit.id
    await api.discardBreakpoint(id)
    pendingBreakpoints.update((list) => list.filter((b) => b.id !== id))
    activeQueueIdx = Math.min(activeQueueIdx, $pendingBreakpoints.length - 1)
  }

  function formatHeaders(h: Record<string, string>): string {
    return Object.entries(h)
      .map(([k, v]) => `${k}: ${v}`)
      .join('\n')
  }

  function mc(m: string): string {
    return METHOD_COLORS[m.toUpperCase()] ?? '#64748b'
  }

  $: totalBpEndpoints = $collections
    .flatMap((c) => c.folders.flatMap((f) => f.endpoints))
    .filter((e) => e.breakpoint).length
</script>

<div class="flex min-w-0 flex-1 flex-col overflow-hidden">
  <div class="flex shrink-0 items-center gap-3 border-b border-wire bg-cave-surface px-5 py-3">
    <span class="text-base text-accent">⚡</span>
    <h2 class="font-semibold text-ink">Breakpoints</h2>
    {#if totalBpEndpoints > 0}
      <span class="rounded-full bg-accent/10 px-2 py-0.5 text-xs text-accent"
        >{totalBpEndpoints} active</span
      >
    {/if}
    {#if $pendingBreakpoints.length > 0}
      <span class="rounded-full bg-warn/10 px-2 py-0.5 text-xs text-warn"
        >{$pendingBreakpoints.length} pending</span
      >
    {/if}
  </div>

  <div class="tab-strip">
    <button
      class="tab-btn {panelTab === 'endpoints' ? 'active' : ''}"
      on:click={() => (panelTab = 'endpoints')}>Endpoints</button
    >
    <button
      class="tab-btn {panelTab === 'queue' ? 'active' : ''}"
      on:click={() => (panelTab = 'queue')}
    >
      Queue{$pendingBreakpoints.length > 0 ? ` (${$pendingBreakpoints.length})` : ''}
    </button>
  </div>

  {#if panelTab === 'endpoints'}
    <div class="flex-1 overflow-y-auto p-4">
      <p class="mb-3 text-xs text-ink-muted">
        Enable breakpoint mode on specific endpoints. Matched requests will pause and appear in the
        Queue.
      </p>
      {#if $collections.length === 0}
        <p class="py-6 text-center text-sm text-ink-muted">
          No collections. Create one in the Collections tab.
        </p>
      {:else}
        <div class="flex flex-col gap-2">
          {#each $collections as col}
            {@const colState = getColBpState(col)}
            <div class="overflow-hidden rounded-lg border border-wire">
              <div class="flex items-center gap-2 bg-cave-surface px-3 py-2.5">
                <button
                  class="flex h-4 w-4 shrink-0 items-center justify-center bg-transparent text-[10px] text-ink-muted"
                  on:click={() => toggleColExpand(col.id)}
                >
                  {expandedCols.has(col.id) ? '▾' : '▸'}
                </button>
                <label class="flex flex-1 cursor-pointer items-center gap-2">
                  <input
                    type="checkbox"
                    checked={colState !== 'none'}
                    on:change={() => toggleColBp(col, colState)}
                  />
                  <span class="text-sm text-accent">⬡</span>
                  <span class="flex-1 font-medium text-ink">{col.name}</span>
                  {#if colState === 'some'}
                    <span class="text-[11px] italic text-ink-muted">partial</span>
                  {:else if colState === 'all'}
                    <span class="text-[11px] text-accent">all active</span>
                  {/if}
                </label>
              </div>

              {#if expandedCols.has(col.id)}
                <div class="bg-cave-deep pb-1 pt-0.5">
                  {#each col.folders as folder}
                    {@const folderState = getFolderBpState(folder)}
                    <div class="px-3">
                      <div class="flex items-center gap-2 py-1.5">
                        <button
                          class="flex h-4 w-4 shrink-0 items-center justify-center bg-transparent text-[10px] text-ink-muted"
                          on:click={() => toggleFolderExpand(folder.id)}
                        >
                          {expandedFolders.has(folder.id) ? '▾' : '▸'}
                        </button>
                        <label class="flex flex-1 cursor-pointer items-center gap-2">
                          <input
                            type="checkbox"
                            checked={folderState !== 'none'}
                            on:change={() => toggleFolderBp(col.id, folder, folderState)}
                          />
                          <span class="text-xs"
                            >{expandedFolders.has(folder.id) ? '📂' : '📁'}</span
                          >
                          <span class="text-sm text-ink-mid">{folder.name}</span>
                          {#if folderState === 'some'}
                            <span class="ml-1 text-[11px] italic text-ink-muted">partial</span>
                          {/if}
                        </label>
                      </div>

                      {#if expandedFolders.has(folder.id)}
                        {#each folder.endpoints as ep}
                          <label
                            class="flex cursor-pointer items-center gap-2 py-1 pl-7 {ep.breakpoint
                              ? 'text-ink'
                              : 'text-ink-mid'}"
                          >
                            <input
                              type="checkbox"
                              checked={ep.breakpoint}
                              on:change={() => toggleEpBp(col.id, folder.id, ep)}
                            />
                            <span
                              class="tag method-badge shrink-0"
                              style="background:{mc(ep.method)}"
                              >{ep.method === '*' ? 'ANY' : ep.method}</span
                            >
                            <span class="min-w-0 flex-1 truncate font-mono text-xs"
                              >{ep.path}</span
                            >
                            <span class="max-w-[140px] truncate text-xs text-ink-muted"
                              >{ep.name}</span
                            >
                            {#if ep.breakpoint}
                              <span class="text-xs text-accent">⚡</span>
                            {/if}
                          </label>
                        {/each}
                      {/if}
                    </div>
                  {/each}
                </div>
              {/if}
            </div>
          {/each}
        </div>
      {/if}
    </div>
  {:else if panelTab === 'queue'}
    {#if $pendingBreakpoints.length === 0}
      <div class="flex flex-1 flex-col items-center justify-center gap-3 text-center">
        <span class="text-4xl opacity-20">⚡</span>
        <p class="text-sm text-ink-mid">No pending breakpoints</p>
        <p class="max-w-xs text-xs text-ink-muted">
          Enable breakpoints on endpoints and trigger requests to capture them here.
        </p>
      </div>
    {:else}
      <div class="flex min-h-0 flex-1 overflow-hidden">
        <aside class="flex w-52 shrink-0 flex-col border-r border-wire bg-cave-deep">
          <div
            class="flex shrink-0 items-center justify-between border-b border-wire px-3 py-2.5"
          >
            <span class="text-[11px] font-semibold uppercase tracking-[0.08em] text-ink-muted"
              >Queue ({$pendingBreakpoints.length})</span
            >
          </div>
          <div class="flex-1 overflow-y-auto py-1">
            {#each $pendingBreakpoints as hit, i}
              <button
                class="flex w-full items-center gap-2 px-3 py-2.5 text-left transition-colors
                  {activeQueueIdx === i
                  ? 'bg-accent/10 text-ink'
                  : 'text-ink-mid hover:bg-cave-raised hover:text-ink'}"
                on:click={() => (activeQueueIdx = i)}
              >
                <span class="h-1.5 w-1.5 shrink-0 animate-pulse-slow rounded-full bg-accent" />
                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-1.5">
                    <span
                      class="tag method-badge shrink-0 text-[9px]"
                      style="background:{mc(hit.method)};padding:1px 4px">{hit.method}</span
                    >
                    <span class="truncate font-mono text-xs">{hit.path}</span>
                  </div>
                  <span class="text-[11px] text-ink-muted"
                    >{new Date(hit.timestamp).toLocaleTimeString()}</span
                  >
                </div>
              </button>
            {/each}
          </div>
        </aside>

        {#if activeHit}
          <div class="flex min-w-0 flex-1 flex-col overflow-hidden">
            <div class="grid flex-1 grid-cols-2 overflow-hidden">
              <div class="flex flex-col overflow-hidden border-r border-wire">
                <div
                  class="flex shrink-0 flex-col gap-1 border-b border-wire bg-cave-surface px-4 py-2.5"
                >
                  <span class="form-label">Incoming Request</span>
                  <div class="flex items-center gap-2">
                    <span
                      class="shrink-0 rounded-[3px] px-1.5 py-0.5 font-mono text-[10px] font-bold uppercase"
                      style="background:{mc(activeHit.method)};color:#0d0f17"
                      >{activeHit.method}</span
                    >
                    <span class="truncate font-mono text-xs text-ink">
                      {activeHit.path}{activeHit.query ? '?' + activeHit.query : ''}
                    </span>
                  </div>
                </div>

                <div class="tab-strip">
                  <button
                    class="tab-btn {reqTab === 'body' ? 'active' : ''}"
                    on:click={() => (reqTab = 'body')}>Body</button
                  >
                  <button
                    class="tab-btn {reqTab === 'headers' ? 'active' : ''}"
                    on:click={() => (reqTab = 'headers')}>Headers</button
                  >
                </div>

                <div class="flex flex-1 flex-col overflow-y-auto p-3">
                  {#if reqTab === 'body'}
                    {#if activeHit.reqBody}
                      <JsonEditor value={activeHit.reqBody} readOnly={true} />
                    {:else}
                      <div class="py-8 text-center text-sm text-ink-muted">No request body</div>
                    {/if}
                  {:else}
                    <pre
                      class="whitespace-pre-wrap break-all font-mono text-xs leading-relaxed text-ink-mid"
                      >{formatHeaders(activeHit.reqHeaders)}</pre
                    >
                  {/if}
                </div>
              </div>

              <div class="flex flex-col overflow-hidden">
                <div
                  class="flex shrink-0 items-center gap-3 border-b border-wire bg-cave-surface px-4 py-2.5"
                >
                  <span class="form-label">Response (editable)</span>
                  {#if editingResponse}
                    <span
                      class="font-mono text-lg font-bold"
                      style="color:{statusColor(editingResponse.statusCode)}"
                      >{editingResponse.statusCode}</span
                    >
                  {/if}
                </div>

                <div class="tab-strip">
                  <button
                    class="tab-btn {respTab === 'body' ? 'active' : ''}"
                    on:click={() => (respTab = 'body')}>Body</button
                  >
                  <button
                    class="tab-btn {respTab === 'headers' ? 'active' : ''}"
                    on:click={() => (respTab = 'headers')}>Headers</button
                  >
                  <button
                    class="tab-btn {respTab === 'cookies' ? 'active' : ''}"
                    on:click={() => (respTab = 'cookies')}>Cookies</button
                  >
                </div>

                {#if editingResponse}
                  <div class="flex min-h-0 flex-1 flex-col overflow-hidden">
                    <ResponseEditor
                      response={editingResponse}
                      activeTab={respTab}
                      on:update={handleResponseUpdate}
                    />
                  </div>
                {/if}
              </div>
            </div>

            <div
              class="flex shrink-0 items-center justify-between border-t border-wire bg-cave-surface px-4 py-3"
            >
              <div class="flex items-center gap-2">
                <span class="font-mono text-xs text-ink-muted">
                  {activeHit.endpoint.name || activeHit.endpoint.path}
                </span>
              </div>
              <div class="flex gap-2">
                <button class="btn btn-danger" on:click={discard}>⊘ Discard</button>
                <button class="btn btn-primary px-4" on:click={release}>
                  ▶ Release with {editingResponse?.statusCode}
                </button>
              </div>
            </div>
          </div>
        {/if}
      </div>
    {/if}
  {/if}
</div>
