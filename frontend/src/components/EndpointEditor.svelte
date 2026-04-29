<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { Endpoint } from '../lib/types'
  import { METHOD_COLORS, newMockResponse } from '../lib/types'
  import ResponseEditor from './ResponseEditor.svelte'

  export let endpoint: Endpoint
  export let collectionId: string
  export let folderId: string

  const dispatch = createEventDispatcher()

  let editingEndpoint: Endpoint = JSON.parse(JSON.stringify(endpoint))
  let activeResponseIdx = editingEndpoint.activeIdx ?? 0
  let dirty = false

  const METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS', '*']

  $: if (endpoint) {
    editingEndpoint = JSON.parse(JSON.stringify(endpoint))
    activeResponseIdx = editingEndpoint.activeIdx ?? 0
    dirty = false
  }

  $: activeResponse = editingEndpoint.responses[activeResponseIdx] ?? null

  function markDirty() {
    dirty = true
  }

  function addResponse() {
    const resp = newMockResponse()
    resp.name = `Response ${editingEndpoint.responses.length + 1}`
    editingEndpoint.responses = [...editingEndpoint.responses, resp]
    activeResponseIdx = editingEndpoint.responses.length - 1
    dirty = true
  }

  function removeResponse(i: number) {
    if (editingEndpoint.responses.length <= 1) return
    editingEndpoint.responses = editingEndpoint.responses.filter((_, idx) => idx !== i)
    if (activeResponseIdx >= editingEndpoint.responses.length) {
      activeResponseIdx = editingEndpoint.responses.length - 1
    }
    dirty = true
  }

  function save() {
    editingEndpoint.activeIdx = activeResponseIdx
    dispatch('save', { collectionId, folderId, endpoint: editingEndpoint })
    dirty = false
  }

  function onResponseChange() {
    editingEndpoint.responses = [...editingEndpoint.responses]
    dirty = true
  }
</script>

<div class="flex h-full flex-col overflow-hidden">
  <!-- Header -->
  <div class="flex shrink-0 flex-col gap-2.5 border-b border-wire bg-cave-surface px-4 py-3">
    <!-- Method + path + save -->
    <div class="flex items-center gap-2">
      <div class="flex shrink-0 gap-0.5">
        {#each METHODS as m}
          <button
            class="rounded-sm border px-2 py-1 font-mono text-[10px] font-bold uppercase tracking-[0.04em] transition-all"
            style="--mc:{METHOD_COLORS[m] || '#64748b'};
              {editingEndpoint.method === m
              ? `background:color-mix(in srgb,var(--mc) 15%,transparent);color:var(--mc);border-color:var(--mc)`
              : 'background:transparent;color:#4a5368;border-color:#252a40'}"
            on:click={() => {
              editingEndpoint.method = m
              markDirty()
            }}>{m === '*' ? 'ANY' : m}</button
          >
        {/each}
      </div>

      <input
        class="flex-1 font-mono text-sm"
        placeholder="/api/path/:param"
        bind:value={editingEndpoint.path}
        on:input={markDirty}
      />

      <button class="btn btn-primary shrink-0" on:click={save} disabled={!dirty}>
        {dirty ? 'Save' : 'Saved ✓'}
      </button>
    </div>

    <!-- Meta row -->
    <div class="flex flex-wrap items-center gap-4">
      <div class="flex flex-col gap-1">
        <span class="form-label">Name</span>
        <input
          class="w-40"
          placeholder="Endpoint name"
          bind:value={editingEndpoint.name}
          on:input={markDirty}
        />
      </div>

      <div class="flex flex-col gap-1">
        <span class="form-label">Delay</span>
        <div class="flex items-center gap-1">
          <input
            type="number"
            min="0"
            max="60000"
            class="w-20"
            bind:value={editingEndpoint.delayMs}
            on:input={markDirty}
          />
          <span class="text-[11px] text-ink-muted">ms</span>
        </div>
      </div>

      <div class="flex flex-col gap-1">
        <span class="form-label">Strategy</span>
        <select bind:value={editingEndpoint.strategy} on:change={markDirty}>
          <option value="fixed">Fixed</option>
          <option value="cycle">Cycle</option>
          <option value="random">Random</option>
        </select>
      </div>

      <!-- Breakpoint toggle -->
      <div class="flex items-center gap-2">
        <span class="form-label">Breakpoint</span>
        <button
          class="flex items-center gap-1.5 bg-transparent {editingEndpoint.breakpoint
            ? 'toggle-on'
            : ''}"
          on:click={() => {
            editingEndpoint.breakpoint = !editingEndpoint.breakpoint
            markDirty()
          }}
          title="Pause incoming requests at this endpoint for live editing"
        >
          <span class="toggle-track"><span class="toggle-thumb" /></span>
          {#if editingEndpoint.breakpoint}
            <span class="text-[10px] font-bold tracking-[0.05em] text-accent">ON</span>
          {/if}
        </button>
      </div>

      <!-- WebSocket toggle -->
      <div class="flex items-center gap-2">
        <span class="form-label">WebSocket</span>
        <button
          class="flex items-center bg-transparent {editingEndpoint.wsEnabled ? 'toggle-on' : ''}"
          on:click={() => {
            editingEndpoint.wsEnabled = !editingEndpoint.wsEnabled
            markDirty()
          }}
        >
          <span class="toggle-track"><span class="toggle-thumb" /></span>
        </button>
      </div>

      <!-- Proxy -->
      {#if editingEndpoint.proxyUrl !== undefined}
        <div class="flex min-w-[200px] flex-1 flex-col gap-1">
          <span class="form-label">Proxy URL</span>
          <input
            class="w-full font-mono text-xs"
            placeholder="https://api.example.com"
            bind:value={editingEndpoint.proxyUrl}
            on:input={markDirty}
          />
        </div>
      {:else}
        <button
          class="btn btn-ghost self-end px-2.5 py-1 text-[11px]"
          on:click={() => {
            editingEndpoint.proxyUrl = ''
            markDirty()
          }}>+ Add Proxy</button
        >
      {/if}
    </div>
  </div>

  <!-- Response tabs bar -->
  <div class="flex shrink-0 items-center gap-2 border-b border-wire bg-cave-deep px-3 py-2">
    <span class="mr-1 shrink-0 text-[10px] font-semibold uppercase tracking-[0.08em] text-ink-muted"
      >Responses</span
    >
    <div class="flex flex-1 gap-1 overflow-x-auto pb-0.5">
      {#each editingEndpoint.responses as resp, i}
        <div class="flex shrink-0 items-center">
          <button
            class="flex items-center gap-1.5 rounded-sm border px-2.5 py-1 text-xs transition-all
              {activeResponseIdx === i
              ? 'border-wire-hi bg-cave-elevated text-ink'
              : 'border-transparent bg-cave-raised text-ink-mid hover:border-wire hover:text-ink'}"
            on:click={() => (activeResponseIdx = i)}
          >
            <span
              class="font-mono text-[11px] font-bold"
              style="color:{resp.statusCode >= 500
                ? '#f43f5e'
                : resp.statusCode >= 400
                ? '#f59e0b'
                : resp.statusCode >= 300
                ? '#3b82f6'
                : '#22d3a0'}">{resp.statusCode}</span
            >
            <span class="max-w-[100px] overflow-hidden text-ellipsis whitespace-nowrap"
              >{resp.name}</span
            >
            {#if editingEndpoint.strategy === 'fixed' && activeResponseIdx === i}
              <span class="h-1.5 w-1.5 shrink-0 rounded-full bg-accent" title="Active response" />
            {/if}
          </button>
          {#if editingEndpoint.responses.length > 1}
            <button
              class="ml-[-2px] p-0.5 text-[9px] text-ink-muted hover:text-err"
              on:click|stopPropagation={() => removeResponse(i)}>✕</button
            >
          {/if}
        </div>
      {/each}
    </div>
    <button class="btn btn-ghost shrink-0 px-2.5 py-1 text-[11px]" on:click={addResponse}>
      + Add
    </button>
  </div>

  <!-- Response editor -->
  <div class="flex flex-1 flex-col overflow-hidden">
    {#if activeResponse}
      <ResponseEditor
        bind:response={editingEndpoint.responses[activeResponseIdx]}
        on:change={onResponseChange}
      />
    {/if}
  </div>
</div>
