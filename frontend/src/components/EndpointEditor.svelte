<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { Endpoint, MockResponse } from '../lib/types'
  import { METHOD_COLORS, newMockResponse, statusColor } from '../lib/types'
  import { tabEdits, activeTabId } from '../stores/index'
  import ResponseEditor from './ResponseEditor.svelte'

  function prettifyBody(body: string, bodyType: string): string {
    if (bodyType === 'json') {
      try {
        return JSON.stringify(JSON.parse(body), null, 4)
      } catch {
        return body
      }
    }
    if (bodyType === 'xml' || bodyType === 'html') {
      try {
        const mime = bodyType === 'html' ? 'text/html' : 'application/xml'
        const parser = new DOMParser()
        const doc = parser.parseFromString(body.trim(), mime)
        const root = bodyType === 'html' ? doc.body : doc.documentElement
        if (!root) return body
        return formatXmlNode(root, 0)
      } catch {
        return body
      }
    }
    return body
  }

  function formatXmlNode(node: Node, depth: number): string {
    const pad = '    '.repeat(depth)
    if (node.nodeType === Node.TEXT_NODE) {
      return node.textContent?.trim() ?? ''
    }
    if (node.nodeType === Node.ELEMENT_NODE) {
      const el = node as Element
      const attrs = Array.from(el.attributes).map((a) => ` ${a.name}="${a.value}"`).join('')
      const tag = el.tagName.toLowerCase()
      const nonEmpty = Array.from(el.childNodes).filter(
        (n) => n.nodeType !== Node.TEXT_NODE || (n.textContent?.trim() ?? ''),
      )
      if (nonEmpty.length === 0) return `${pad}<${tag}${attrs}/>`
      const hasEls = nonEmpty.some((n) => n.nodeType === Node.ELEMENT_NODE)
      if (!hasEls) return `${pad}<${tag}${attrs}>${el.textContent?.trim() ?? ''}</${tag}>`
      const children = nonEmpty.map((c) => formatXmlNode(c, depth + 1)).filter(Boolean).join('\n')
      return `${pad}<${tag}${attrs}>\n${children}\n${pad}</${tag}>`
    }
    return ''
  }

  export let endpoint: Endpoint
  export let collectionId: string
  export let folderId: string
  export let tabId: string

  const dispatch = createEventDispatcher()

  const METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', 'OPTIONS', '*']

  let editingEndpoint: Endpoint = JSON.parse(JSON.stringify(endpoint))
  let activeResponseIdx = editingEndpoint.activeIdx ?? 0
  let activeTab: 'body' | 'headers' | 'cookies' | 'configure' = 'body'
  let dirty = false
  let prevEndpointId = endpoint.id

  $: if (endpoint.id !== prevEndpointId) {
    prevEndpointId = endpoint.id
    editingEndpoint = JSON.parse(JSON.stringify(endpoint))
    activeResponseIdx = editingEndpoint.activeIdx ?? 0
    activeTab = 'body'
    dirty = false
  }

  $: tabEdits.update((m) => ({ ...m, [tabId]: editingEndpoint }))

  function markDirty() {
    dirty = true
    dispatch('dirtyChange', true)
  }

  function addResponse() {
    const resp = newMockResponse()
    resp.name = `Response ${editingEndpoint.responses.length + 1}`
    editingEndpoint.responses = [...editingEndpoint.responses, resp]
    activeResponseIdx = editingEndpoint.responses.length - 1
    markDirty()
  }

  function removeResponse(i: number) {
    if (editingEndpoint.responses.length <= 1) return
    editingEndpoint.responses = editingEndpoint.responses.filter((_, idx) => idx !== i)
    activeResponseIdx = Math.min(activeResponseIdx, editingEndpoint.responses.length - 1)
    markDirty()
  }

  function handleResponseUpdate(e: CustomEvent<MockResponse>) {
    editingEndpoint.responses = editingEndpoint.responses.map((r, i) =>
      i === activeResponseIdx ? e.detail : r,
    )
    markDirty()
  }

  function save() {
    editingEndpoint.activeIdx = activeResponseIdx
    dispatch('save', { collectionId, folderId, endpoint: editingEndpoint })
    dirty = false
    dispatch('dirtyChange', false)
  }

  function mc(m: string) {
    return METHOD_COLORS[m.toUpperCase()] || '#64748b'
  }

  $: activeResponse = editingEndpoint.responses[activeResponseIdx] ?? null

  function handleKeydown(e: KeyboardEvent) {
    if (e.ctrlKey && e.key === 's' && $activeTabId === tabId) {
      e.preventDefault()
      if (dirty) save()
    }
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="flex h-full flex-col overflow-hidden">
  <div class="flex shrink-0 items-center gap-2 border-b border-wire bg-cave-surface px-4 py-3">
    <select
      class="w-28 shrink-0 font-mono font-bold"
      value={editingEndpoint.method}
      style="color:{mc(editingEndpoint.method)}"
      on:change={(e) => {
        editingEndpoint.method = e.currentTarget.value
        markDirty()
      }}
    >
      {#each METHODS as m}
        <option value={m} style="color:{mc(m)}">{m === '*' ? 'ANY' : m}</option>
      {/each}
    </select>

    <input
      class="flex-1 font-mono"
      placeholder="/api/path/:param"
      bind:value={editingEndpoint.path}
      on:input={markDirty}
    />

    <div class="flex shrink-0 items-center gap-2">
      {#if dirty}
        <span class="h-2.5 w-2.5 rounded-full bg-warn" title="Unsaved changes" />
      {/if}
      <button class="btn btn-primary" on:click={save} disabled={!dirty}>
        {dirty ? 'Save' : 'Saved ✓'}
      </button>
    </div>
  </div>

  <div class="flex shrink-0 items-center gap-1 border-b border-wire bg-cave-deep px-3 py-1.5">
    <span class="mr-1 shrink-0 text-[11px] font-semibold uppercase tracking-[0.08em] text-ink-muted">
      Responses
    </span>
    <div class="flex flex-1 gap-0.5 overflow-x-auto">
      {#each editingEndpoint.responses as resp, i}
        <div class="flex shrink-0 items-center">
          <button
            class="flex items-center gap-1.5 rounded-sm border px-2 py-1 text-xs transition-all
              {activeResponseIdx === i
              ? 'border-wire-hi bg-cave-elevated text-ink'
              : 'border-transparent text-ink-mid hover:border-wire hover:text-ink'}"
            on:click={() => (activeResponseIdx = i)}
          >
            <span class="font-mono font-bold" style="color:{statusColor(resp.statusCode)}"
              >{resp.statusCode}</span
            >
            <span class="max-w-[80px] overflow-hidden text-ellipsis whitespace-nowrap"
              >{resp.name}</span
            >
            {#if editingEndpoint.strategy === 'fixed' && activeResponseIdx === i}
              <span class="h-1.5 w-1.5 shrink-0 rounded-full bg-accent" />
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
    <button class="btn btn-ghost shrink-0 px-2 py-1 text-xs" on:click={addResponse}>+ Add</button>
  </div>

  <div class="tab-strip items-center">
    <button
      class="tab-btn {activeTab === 'body' ? 'active' : ''}"
      on:click={() => (activeTab = 'body')}
    >
      Body{activeResponse && activeResponse.body.trim() ? ' ●' : ''}
    </button>
    <button
      class="tab-btn {activeTab === 'headers' ? 'active' : ''}"
      on:click={() => (activeTab = 'headers')}
    >
      Headers{activeResponse
        ? activeResponse.headers.filter((h) => h.enabled && h.key).length > 0
          ? ` (${activeResponse.headers.filter((h) => h.enabled && h.key).length})`
          : ''
        : ''}
    </button>
    <button
      class="tab-btn {activeTab === 'cookies' ? 'active' : ''}"
      on:click={() => (activeTab = 'cookies')}
    >
      Cookies{activeResponse && activeResponse.cookies.length > 0
        ? ` (${activeResponse.cookies.length})`
        : ''}
    </button>
    <button
      class="tab-btn {activeTab === 'configure' ? 'active' : ''}"
      on:click={() => (activeTab = 'configure')}
    >
      Configure
    </button>
    {#if activeTab === 'body' && activeResponse && ['json', 'xml', 'html'].includes(activeResponse.bodyType)}
      <button
        class="btn btn-ghost ml-auto px-2.5 py-0.5 text-xs"
        on:click={() => {
          if (!activeResponse) return
          const formatted = prettifyBody(activeResponse.body, activeResponse.bodyType)
          if (formatted === activeResponse.body) return
          editingEndpoint.responses = editingEndpoint.responses.map((r, i) =>
            i === activeResponseIdx ? { ...r, body: formatted } : r,
          )
          save()
        }}>Prettify</button
      >
    {/if}
  </div>

  <div class="flex flex-1 flex-col overflow-hidden">
    {#if activeTab === 'configure'}
      <div class="flex flex-col gap-5 overflow-y-auto p-5">
        <div class="flex flex-wrap gap-6">
          <div class="flex flex-col gap-1.5">
            <span class="form-label">Delay (ms)</span>
            <input
              type="number"
              min="0"
              max="60000"
              class="w-28"
              bind:value={editingEndpoint.delayMs}
              on:input={markDirty}
            />
          </div>

          <div class="flex flex-col gap-1.5">
            <span class="form-label">Strategy</span>
            <select bind:value={editingEndpoint.strategy} on:change={markDirty}>
              <option value="fixed">Fixed</option>
              <option value="cycle">Cycle</option>
              <option value="random">Random</option>
            </select>
          </div>

          <div class="flex flex-col gap-3 self-end pb-1">
            <label class="flex cursor-pointer items-center gap-3">
              <span class="form-label">WebSocket</span>
              <button
                class="flex items-center bg-transparent {editingEndpoint.wsEnabled
                  ? 'toggle-on'
                  : ''}"
                on:click={() => {
                  editingEndpoint.wsEnabled = !editingEndpoint.wsEnabled
                  markDirty()
                }}
              >
                <span class="toggle-track"><span class="toggle-thumb" /></span>
              </button>
            </label>
          </div>
        </div>

        <div class="flex flex-col gap-1.5">
          <span class="form-label">Proxy URL</span>
          {#if editingEndpoint.proxyUrl !== undefined}
            <div class="flex items-center gap-2">
              <input
                class="flex-1 font-mono text-sm"
                placeholder="https://api.example.com"
                bind:value={editingEndpoint.proxyUrl}
                on:input={markDirty}
              />
              <button
                class="btn-icon hover:text-err"
                on:click={() => {
                  editingEndpoint.proxyUrl = undefined
                  markDirty()
                }}>✕</button
              >
            </div>
          {:else}
            <button
              class="btn btn-ghost self-start"
              on:click={() => {
                editingEndpoint.proxyUrl = ''
                markDirty()
              }}>+ Add Proxy</button
            >
          {/if}
        </div>
      </div>
    {:else if activeResponse}
      {#key activeResponseIdx}
        <ResponseEditor
          response={activeResponse}
          {activeTab}
          on:update={handleResponseUpdate}
        />
      {/key}
    {/if}
  </div>
</div>
