<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { MockResponse } from '../lib/types'
  import { statusColor } from '../lib/types'
  import JsonEditor from './JsonEditor.svelte'
  import KeyValueEditor from './KeyValueEditor.svelte'

  export let response: MockResponse
  export let activeTab: 'body' | 'headers' | 'cookies' = 'body'

  const dispatch = createEventDispatcher()

  let r: MockResponse = JSON.parse(JSON.stringify(response))

  $: if (response.id !== r.id) {
    r = JSON.parse(JSON.stringify(response))
  }

  function emit() {
    dispatch('update', { ...r })
  }

  function set(patch: Partial<MockResponse>) {
    r = { ...r, ...patch }
    emit()
  }

  function setBodyType(v: string) {
    set({ bodyType: v as MockResponse['bodyType'] })
  }

  function onBodyChange(e: CustomEvent<string>) {
    r = { ...r, body: e.detail }
    emit()
  }

  function onHeadersChange(e: CustomEvent<typeof r.headers>) {
    r = { ...r, headers: e.detail }
    emit()
  }

  function addCookie() {
    r = {
      ...r,
      cookies: [
        ...r.cookies,
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
      ],
    }
    emit()
  }

  function removeCookie(i: number) {
    r = { ...r, cookies: r.cookies.filter((_, idx) => idx !== i) }
    emit()
  }
</script>

<div class="flex h-full flex-col">
  <div class="flex flex-wrap items-center gap-4 border-b border-wire bg-cave-surface px-4 py-2.5">
    <div class="flex flex-col gap-1">
      <span class="form-label">Status</span>
      <input
        type="number"
        min="100"
        max="599"
        class="w-16 border-0 bg-transparent p-0 font-mono text-xl font-bold focus:ring-0"
        value={r.statusCode}
        on:change={(e) => set({ statusCode: parseInt(e.currentTarget.value) || 200 })}
        style="color:{statusColor(r.statusCode)}"
      />
    </div>

    <div class="flex flex-col gap-1">
      <span class="form-label">Body Type</span>
      <select value={r.bodyType} on:change={(e) => setBodyType(e.currentTarget.value)}>
        <option value="json">JSON</option>
        <option value="text">Text</option>
        <option value="html">HTML</option>
        <option value="xml">XML</option>
        <option value="none">None</option>
      </select>
    </div>
  </div>

  <div class="flex flex-1 flex-col overflow-hidden">
    {#if activeTab === 'body'}
      {#if r.bodyType === 'none'}
        <div class="flex flex-1 items-center justify-center py-10 text-sm text-ink-muted">No body</div>
      {:else}
        <JsonEditor value={r.body} on:change={onBodyChange} />
      {/if}
    {:else if activeTab === 'headers'}
      <div class="overflow-y-auto p-3">
        <KeyValueEditor
          pairs={r.headers}
          keyPlaceholder="Header name"
          valuePlaceholder="Value"
          on:change={onHeadersChange}
        />
      </div>
    {:else if activeTab === 'cookies'}
      <div class="flex flex-col gap-1.5 overflow-y-auto p-3">
        {#each r.cookies as cookie, i}
          <div class="flex flex-wrap items-center gap-1.5 rounded-sm px-1 py-1 hover:bg-cave-raised">
            <input
              class="min-w-0 flex-1"
              placeholder="Name"
              value={cookie.name}
              on:input={(e) => {
                r.cookies[i] = { ...cookie, name: e.currentTarget.value }
                r = { ...r }
                emit()
              }}
            />
            <input
              class="min-w-0 flex-1"
              placeholder="Value"
              value={cookie.value}
              on:input={(e) => {
                r.cookies[i] = { ...cookie, value: e.currentTarget.value }
                r = { ...r }
                emit()
              }}
            />
            <input
              class="w-20"
              placeholder="Path"
              value={cookie.path}
              on:input={(e) => {
                r.cookies[i] = { ...cookie, path: e.currentTarget.value }
                r = { ...r }
                emit()
              }}
            />
            <label
              class="flex cursor-pointer items-center gap-1.5 whitespace-nowrap text-xs text-ink-mid"
            >
              <input
                type="checkbox"
                checked={cookie.httpOnly}
                on:change={(e) => {
                  r.cookies[i] = { ...cookie, httpOnly: e.currentTarget.checked }
                  r = { ...r }
                  emit()
                }}
              />
              HttpOnly
            </label>
            <label
              class="flex cursor-pointer items-center gap-1.5 whitespace-nowrap text-xs text-ink-mid"
            >
              <input
                type="checkbox"
                checked={cookie.secure}
                on:change={(e) => {
                  r.cookies[i] = { ...cookie, secure: e.currentTarget.checked }
                  r = { ...r }
                  emit()
                }}
              />
              Secure
            </label>
            <button
              class="btn-icon text-[10px] hover:text-err"
              on:click={() => removeCookie(i)}>✕</button
            >
          </div>
        {/each}
        <button class="btn btn-ghost self-start px-2.5 py-1 text-xs" on:click={addCookie}>
          + Add Cookie
        </button>
      </div>
    {/if}
  </div>
</div>
