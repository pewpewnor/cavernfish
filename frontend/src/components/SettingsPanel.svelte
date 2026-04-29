<script lang="ts">
  import { collections, patchCollection } from '../stores/index'
  import * as api from '../lib/api'
  import Modal from './Modal.svelte'

  let showImport = false
  let importTab: 'native' | 'openapi' = 'native'
  let importText = ''
  let importError = ''
  let importLoading = false

  let editingCollection: string | null = null
  let editingVars: { key: string; value: string }[] = []
  let editingName = ''
  let editingDesc = ''

  function openVarEditor(cid: string) {
    const col = $collections.find((c) => c.id === cid)
    if (!col) return
    editingCollection = cid
    editingName = col.name
    editingDesc = col.description
    editingVars = Object.entries(col.variables ?? {}).map(([key, value]) => ({ key, value }))
  }

  async function saveVars() {
    if (!editingCollection) return
    const variables: Record<string, string> = {}
    editingVars.forEach((v) => {
      if (v.key.trim()) variables[v.key.trim()] = v.value
    })
    const updated = await api.updateCollection(editingCollection, {
      name: editingName,
      description: editingDesc,
      variables,
    })
    patchCollection(updated)
    editingCollection = null
  }

  async function exportCollection(id: string) {
    const data = await api.exportCollection(id)
    const blob = new Blob([data], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    const col = $collections.find((c) => c.id === id)
    a.href = url
    a.download = `${col?.name ?? 'collection'}.cavernfish.json`
    a.click()
    URL.revokeObjectURL(url)
  }

  async function doImport() {
    if (!importText.trim()) return
    importError = ''
    importLoading = true
    try {
      let result
      if (importTab === 'openapi') {
        result = await api.importFromOpenAPI(importText)
      } else {
        result = await api.importCollection(importText)
      }
      if (result) {
        collections.update((list) => [...list, result])
        showImport = false
        importText = ''
      }
    } catch (e: unknown) {
      importError = e instanceof Error ? e.message : String(e)
    } finally {
      importLoading = false
    }
  }

  async function handleFileImport(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    importText = await file.text()
  }
</script>

<div class="flex flex-1 flex-col gap-6 overflow-y-auto p-5">
  <div class="shrink-0">
    <h2 class="text-lg font-semibold text-ink">Settings</h2>
  </div>

  <div class="flex max-w-3xl flex-col gap-6">
    <!-- Collections section -->
    <section class="overflow-hidden rounded-lg border border-wire bg-cave-surface">
      <div class="section-header">
        <span class="section-title">Collections</span>
        <button class="btn btn-primary" on:click={() => (showImport = true)}>⬆ Import</button>
      </div>

      {#if $collections.length === 0}
        <div class="px-4 py-5 text-[13px] text-ink-muted">No collections yet.</div>
      {:else}
        <div class="flex flex-col">
          <div
            class="grid grid-cols-[1fr_80px_80px_80px_180px] gap-3 border-b border-wire px-4 py-2 text-[10px] font-semibold uppercase tracking-[0.06em] text-ink-muted"
          >
            <span>Name</span>
            <span>Folders</span>
            <span>Endpoints</span>
            <span>Variables</span>
            <span />
          </div>
          {#each $collections as col}
            {@const totalEndpoints = col.folders.reduce((s, f) => s + f.endpoints.length, 0)}
            <div
              class="grid grid-cols-[1fr_80px_80px_80px_180px] items-center gap-3 border-b border-wire px-4 py-2.5 transition-colors last:border-b-0 hover:bg-cave-raised"
            >
              <div class="flex min-w-0 items-center gap-2">
                <span class="shrink-0 text-base text-accent">⬡</span>
                <div class="min-w-0">
                  <div class="truncate text-[13px] font-medium text-ink">{col.name}</div>
                  {#if col.description}
                    <div class="truncate text-[11px] text-ink-muted">{col.description}</div>
                  {/if}
                </div>
              </div>
              <span class="font-mono text-[13px] text-ink-mid">{col.folders.length}</span>
              <span class="font-mono text-[13px] text-ink-mid">{totalEndpoints}</span>
              <span class="font-mono text-[13px] text-ink-mid"
                >{Object.keys(col.variables ?? {}).length}</span
              >
              <div class="flex justify-end gap-1.5">
                <button
                  class="btn btn-ghost px-2.5 py-1 text-[11px]"
                  on:click={() => openVarEditor(col.id)}>✏ Edit</button
                >
                <button
                  class="btn btn-ghost px-2.5 py-1 text-[11px]"
                  on:click={() => exportCollection(col.id)}>⬇ Export</button
                >
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </section>

    <!-- About section -->
    <section class="overflow-hidden rounded-lg border border-wire bg-cave-surface">
      <div class="section-header">
        <span class="section-title">About</span>
      </div>
      <div class="flex flex-col gap-2 px-4 py-3">
        <div class="flex items-center gap-4">
          <span class="w-28 shrink-0 text-xs text-ink-muted">App</span>
          <span class="text-[13px] text-ink">Cavernfish</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="w-28 shrink-0 text-xs text-ink-muted">Version</span>
          <span class="font-mono text-[13px] text-ink">0.1.0</span>
        </div>
        <div class="flex items-center gap-4">
          <span class="w-28 shrink-0 text-xs text-ink-muted">Data directory</span>
          <span class="font-mono text-[13px] text-ink">~/.cavernfish/</span>
        </div>
      </div>
    </section>
  </div>
</div>

<!-- Import modal -->
{#if showImport}
  <Modal title="Import Collection" width="560px" on:close={() => (showImport = false)}>
    <div class="flex flex-col gap-3">
      <div class="-mx-4 -mt-1 flex border-b border-wire px-3">
        <button
          class="tab-btn {importTab === 'native' ? 'active' : ''}"
          on:click={() => (importTab = 'native')}>Cavernfish JSON</button
        >
        <button
          class="tab-btn {importTab === 'openapi' ? 'active' : ''}"
          on:click={() => (importTab = 'openapi')}>OpenAPI 3.x</button
        >
      </div>

      <div class="flex items-center gap-2.5 pt-2">
        <label
          class="inline-flex cursor-pointer items-center gap-1.5 rounded-sm border border-wire bg-cave-raised px-3 py-1.5 text-xs text-ink-mid transition-all hover:border-wire-hi hover:text-ink"
        >
          📂 Choose file
          <input
            type="file"
            class="hidden"
            accept=".json,.yaml,.yml"
            on:change={handleFileImport}
          />
        </label>
        <span class="text-[11px] text-ink-muted">or paste below</span>
      </div>

      <textarea
        class="w-full resize-y font-mono text-[11px] leading-relaxed"
        bind:value={importText}
        placeholder={importTab === 'openapi'
          ? '{ "openapi": "3.0.0", ... }'
          : '{ "name": "My Collection", "folders": [...] }'}
        rows="10"
      />

      {#if importError}
        <div class="rounded-sm border border-err/20 bg-err/[.08] px-2.5 py-2 text-xs text-err">
          {importError}
        </div>
      {/if}

      <div class="flex justify-end gap-2">
        <button
          class="btn btn-primary"
          on:click={doImport}
          disabled={importLoading || !importText.trim()}
          >{importLoading ? 'Importing…' : 'Import'}</button
        >
        <button class="btn btn-ghost" on:click={() => (showImport = false)}>Cancel</button>
      </div>
    </div>
  </Modal>
{/if}

<!-- Edit collection modal -->
{#if editingCollection}
  <Modal title="Edit Collection" width="500px" on:close={() => (editingCollection = null)}>
    <div class="flex flex-col gap-3.5">
      <label class="flex flex-col gap-1.5">
        <span class="form-label">Name</span>
        <input bind:value={editingName} />
      </label>

      <label class="flex flex-col gap-1.5">
        <span class="form-label">Description</span>
        <input bind:value={editingDesc} placeholder="Optional description" />
      </label>

      <div class="flex flex-col gap-1.5">
        <div class="mb-0.5 flex items-baseline gap-2.5">
          <span class="form-label">Variables</span>
          <span class="text-[11px] text-ink-muted"
            >Use &#123;&#123;variable_name&#125;&#125; in paths and bodies</span
          >
        </div>
        {#each editingVars as v, i}
          <div class="flex items-center gap-1.5">
            <input class="flex-1 font-mono" placeholder="variable_name" bind:value={v.key} />
            <span class="shrink-0 font-mono text-sm text-ink-muted">=</span>
            <input class="flex-1 font-mono" placeholder="value" bind:value={v.value} />
            <button
              class="btn-icon"
              on:click={() => (editingVars = editingVars.filter((_, idx) => idx !== i))}>✕</button
            >
          </div>
        {/each}
        <button
          class="btn btn-ghost self-start px-2.5 py-1 text-[11px]"
          on:click={() => (editingVars = [...editingVars, { key: '', value: '' }])}
          >+ Add Variable</button
        >
      </div>

      <div class="flex justify-end gap-2">
        <button class="btn btn-primary" on:click={saveVars}>Save</button>
        <button class="btn btn-ghost" on:click={() => (editingCollection = null)}>Cancel</button>
      </div>
    </div>
  </Modal>
{/if}
