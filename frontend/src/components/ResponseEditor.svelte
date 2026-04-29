<script lang="ts">
  import type { MockResponse } from '../lib/types'
  import { newCookie, statusColor } from '../lib/types'
  import JsonEditor from './JsonEditor.svelte'
  import KeyValueEditor from './KeyValueEditor.svelte'

  export let response: MockResponse

  let activeTab: 'body' | 'headers' | 'cookies' = 'body'

  const bodyTypes = ['json', 'text', 'html', 'xml', 'none'] as const
  const commonStatuses = [200, 201, 204, 301, 302, 400, 401, 403, 404, 409, 422, 500, 502, 503]

  function formatPretty() {
    if (response.bodyType !== 'json') return
    try {
      response.body = JSON.stringify(JSON.parse(response.body), null, 2)
    } catch {
      /* leave as-is */
    }
  }
</script>

<div class="flex h-full flex-col">
  <!-- Meta row -->
  <div class="flex flex-wrap items-start gap-4 border-b border-wire bg-cave-surface px-4 py-3">
    <!-- Status -->
    <div class="flex flex-col gap-1">
      <span class="form-label">Status</span>
      <input
        type="number"
        min="100"
        max="599"
        class="w-16 border-0 bg-transparent p-0 font-mono text-xl font-bold focus:ring-0"
        bind:value={response.statusCode}
        style="color: {statusColor(response.statusCode)}"
      />
      <div class="flex flex-wrap gap-0.5">
        {#each commonStatuses as s}
          <button
            class="rounded-[3px] border px-1.5 py-0.5 font-mono text-[10px] transition-all
              {response.statusCode === s
              ? 'border-accent bg-accent/10 text-accent'
              : 'border-wire bg-cave-raised text-ink-muted hover:border-wire-hi hover:text-ink'}"
            style={response.statusCode === s ? '' : `--hc:${statusColor(s)}`}
            on:click={() => (response.statusCode = s)}>{s}</button
          >
        {/each}
      </div>
    </div>

    <!-- Body type -->
    <div class="flex flex-col gap-1.5">
      <span class="form-label">Body type</span>
      <div class="flex gap-1">
        {#each bodyTypes as t}
          <button
            class="rounded-full border px-2.5 py-0.5 text-[11px] uppercase tracking-[0.05em] transition-all
              {response.bodyType === t
              ? 'border-accent bg-accent/10 text-accent'
              : 'border-wire bg-cave-raised text-ink-mid hover:border-wire-hi hover:text-ink'}"
            on:click={() => (response.bodyType = t)}>{t}</button
          >
        {/each}
      </div>
    </div>

    <!-- Name -->
    <div class="flex min-w-[140px] flex-1 flex-col gap-1.5">
      <span class="form-label">Response name</span>
      <input class="w-full" bind:value={response.name} placeholder="e.g. 200 OK" />
    </div>
  </div>

  <!-- Tab strip -->
  <div class="tab-strip items-center">
    <button
      class="tab-btn {activeTab === 'body' ? 'active' : ''}"
      on:click={() => (activeTab = 'body')}
    >
      Body {response.body.trim() ? '●' : ''}
    </button>
    <button
      class="tab-btn {activeTab === 'headers' ? 'active' : ''}"
      on:click={() => (activeTab = 'headers')}
    >
      Headers {response.headers.filter((h) => h.enabled && h.key).length > 0
        ? `(${response.headers.filter((h) => h.enabled && h.key).length})`
        : ''}
    </button>
    <button
      class="tab-btn {activeTab === 'cookies' ? 'active' : ''}"
      on:click={() => (activeTab = 'cookies')}
    >
      Cookies {response.cookies.length > 0 ? `(${response.cookies.length})` : ''}
    </button>
    {#if activeTab === 'body' && response.bodyType === 'json'}
      <button class="btn btn-ghost ml-auto px-2.5 py-0.5 text-[11px]" on:click={formatPretty}>
        Format
      </button>
    {/if}
  </div>

  <!-- Tab content -->
  <div class="flex flex-1 flex-col overflow-y-auto p-3">
    {#if activeTab === 'body'}
      {#if response.bodyType === 'none'}
        <div class="py-10 text-center text-[13px] text-ink-muted">No body</div>
      {:else}
        <div class="flex min-h-[150px] flex-1 flex-col">
          <JsonEditor bind:value={response.body} />
        </div>
      {/if}
    {:else if activeTab === 'headers'}
      <KeyValueEditor
        bind:pairs={response.headers}
        keyPlaceholder="Header name"
        valuePlaceholder="Value"
      />
    {:else if activeTab === 'cookies'}
      <div class="flex flex-col gap-1.5">
        {#each response.cookies as cookie, i}
          <div class="flex items-center gap-1.5 rounded-sm px-1 py-1 hover:bg-cave-raised">
            <input class="flex-1 font-mono" placeholder="Name" bind:value={cookie.name} />
            <input class="flex-1 font-mono" placeholder="Value" bind:value={cookie.value} />
            <input class="w-20 font-mono" placeholder="Path" bind:value={cookie.path} />
            <label
              class="flex cursor-pointer items-center gap-1 whitespace-nowrap text-xs text-ink-mid"
            >
              <input type="checkbox" bind:checked={cookie.httpOnly} />
              HttpOnly
            </label>
            <label
              class="flex cursor-pointer items-center gap-1 whitespace-nowrap text-xs text-ink-mid"
            >
              <input type="checkbox" bind:checked={cookie.secure} />
              Secure
            </label>
            <button
              class="btn-icon text-[10px] hover:text-err"
              on:click={() => (response.cookies = response.cookies.filter((_, idx) => idx !== i))}
              >✕</button
            >
          </div>
        {/each}
        <button
          class="btn btn-ghost self-start px-2.5 py-1 text-[11px]"
          on:click={() =>
            (response.cookies = [
              ...response.cookies,
              {
                name: '',
                value: '',
                path: '/',
                domain: '',
                maxAge: 0,
                httpOnly: false,
                secure: false,
                sameSite: 'Lax',
              },
            ])}
        >
          + Add Cookie
        </button>
      </div>
    {/if}
  </div>
</div>
