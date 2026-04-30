<script lang="ts">
  import type { MockResponse } from '../lib/types'
  import { METHOD_COLORS, statusColor } from '../lib/types'
  import { pendingBreakpoints } from '../stores/index'
  import * as api from '../lib/api'
  import ResponseEditor from './ResponseEditor.svelte'
  import JsonEditor from './JsonEditor.svelte'

  let activeIdx = 0
  $: activeHit = $pendingBreakpoints[activeIdx] ?? null

  let editingResponse: MockResponse | null = null
  let prevHitId = ''

  $: if (activeHit?.id !== prevHitId) {
    prevHitId = activeHit?.id ?? ''
    editingResponse = activeHit ? (JSON.parse(JSON.stringify(activeHit.response)) as MockResponse) : null
  }

  function handleResponseUpdate(e: CustomEvent<MockResponse>) {
    editingResponse = e.detail
  }

  async function release() {
    if (!activeHit || !editingResponse) return
    const id = activeHit.id
    await api.releaseBreakpoint(id, editingResponse)
    pendingBreakpoints.update((list) => list.filter((b) => b.id !== id))
    activeIdx = Math.min(activeIdx, $pendingBreakpoints.length - 1)
  }

  async function discard() {
    if (!activeHit) return
    const id = activeHit.id
    await api.discardBreakpoint(id)
    pendingBreakpoints.update((list) => list.filter((b) => b.id !== id))
    activeIdx = Math.min(activeIdx, $pendingBreakpoints.length - 1)
  }

  async function releaseAll() {
    for (const hit of $pendingBreakpoints) {
      await api.releaseBreakpoint(hit.id, hit.response)
    }
    pendingBreakpoints.set([])
    activeIdx = 0
  }

  function formatHeaders(headers: Record<string, string>): string {
    return Object.entries(headers)
      .map(([k, v]) => `${k}: ${v}`)
      .join('\n')
  }

  function mc(m: string): string {
    return METHOD_COLORS[m.toUpperCase()] ?? '#64748b'
  }

  let reqTab: 'body' | 'headers' = 'body'
  let respTab: 'body' | 'headers' | 'cookies' = 'body'
</script>

<div
  class="fixed inset-0 z-[900] flex items-center justify-center bg-cave-base/65 backdrop-blur-sm"
>
  <div
    class="flex h-[640px] max-h-[90vh] w-[960px] max-w-[95vw] flex-col overflow-hidden rounded-lg border border-accent bg-cave-elevated shadow-bp"
  >
    <div class="flex shrink-0 items-center gap-3 border-b border-wire bg-cave-surface px-4 py-3">
      <div class="flex items-center gap-2 text-sm font-bold text-ink">
        <span class="text-base text-accent">⚡</span>
        <span>Breakpoint Hit</span>
        {#if $pendingBreakpoints.length > 1}
          <span
            class="rounded-full bg-accent/[.12] px-1.5 py-0.5 text-[11px] font-semibold text-accent"
          >
            {$pendingBreakpoints.length} pending
          </span>
        {/if}
      </div>

      {#if $pendingBreakpoints.length > 1}
        <div class="flex flex-1 gap-1 overflow-x-auto">
          {#each $pendingBreakpoints as hit, i}
            <button
              class="flex shrink-0 items-center gap-1.5 rounded-sm border px-2 py-0.5 text-[11px] transition-all
                {activeIdx === i
                ? 'border-accent bg-accent/[.12] text-ink'
                : 'border-wire bg-cave-raised text-ink-mid hover:border-wire-hi hover:text-ink'}"
              on:click={() => (activeIdx = i)}
            >
              <span class="font-mono">{hit.method}</span>
              <span class="max-w-[120px] overflow-hidden text-ellipsis whitespace-nowrap font-mono"
                >{hit.path}</span
              >
            </button>
          {/each}
        </div>
      {/if}

      <div class="ml-auto shrink-0">
        {#if $pendingBreakpoints.length > 1}
          <button class="btn btn-ghost text-xs" on:click={releaseAll}>Release All</button>
        {/if}
      </div>
    </div>

    {#if activeHit}
      <div class="grid flex-1 grid-cols-2 overflow-hidden">
        <div class="flex flex-col overflow-hidden border-r border-wire">
          <div
            class="flex shrink-0 flex-col gap-1.5 border-b border-wire bg-cave-surface px-3.5 py-2.5"
          >
            <span class="form-label">Incoming Request</span>
            <div class="flex items-center gap-2">
              <span
                class="shrink-0 rounded-[3px] px-1.5 py-0.5 font-mono text-[10px] font-bold uppercase"
                style="background:{mc(activeHit.method)};color:#0d0f17">{activeHit.method}</span
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

          <div class="flex flex-1 flex-col overflow-y-auto p-2.5">
            {#if reqTab === 'body'}
              {#if activeHit.reqBody}
                <JsonEditor value={activeHit.reqBody} readOnly={true} minHeight="200px" />
              {:else}
                <div class="py-10 text-center text-sm text-ink-muted">No request body</div>
              {/if}
            {:else}
              <pre
                class="whitespace-pre-wrap break-all font-mono text-xs leading-relaxed text-ink-mid">{formatHeaders(
                  activeHit.reqHeaders,
                )}</pre>
            {/if}
          </div>
        </div>

        <div class="flex flex-col overflow-hidden">
          <div
            class="flex shrink-0 items-center gap-3 border-b border-wire bg-cave-surface px-3.5 py-2.5"
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
            <div class="flex flex-1 flex-col overflow-hidden">
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
          {#if activeHit.endpoint.strategy === 'cycle'}
            <span
              class="rounded-[3px] bg-cave-raised px-1.5 py-0.5 font-mono text-[10px] text-ink-muted"
              >cycle</span
            >
          {/if}
        </div>
        <div class="flex gap-2">
          <button class="btn btn-danger" on:click={discard}>⊘ Discard</button>
          <button class="btn btn-primary px-4" on:click={release}>
            ▶ Release with {editingResponse?.statusCode}
          </button>
        </div>
      </div>
    {/if}
  </div>
</div>
