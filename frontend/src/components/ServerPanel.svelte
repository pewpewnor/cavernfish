<script lang="ts">
  import { onMount } from 'svelte'
  import type { ServerInfo, ServerConfig } from '../lib/types'
  import { METHOD_COLORS, statusColor } from '../lib/types'
  import { servers, collections, requestLogs, addLog } from '../stores/index'
  import * as api from '../lib/api'
  import Modal from './Modal.svelte'

  let showNewServer = false
  let editingServer: ServerInfo | null = null
  let expandedLogServer: string | null = null
  let form: Partial<ServerConfig> = {}

  function resetForm(base?: ServerInfo) {
    form = base
      ? {
          id: base.id,
          name: base.name,
          port: base.port,
          collectionId: base.collectionId,
          https: base.https,
        }
      : { name: '', port: 8000, collectionId: '', https: false }
  }

  async function saveServer() {
    if (!form.port || form.port < 1 || form.port > 65535) {
      alert('Invalid port')
      return
    }
    if (editingServer) {
      const updated = await api.updateServer(form as ServerConfig)
      servers.update((list) => list.map((s) => (s.id === updated.id ? updated : s)))
    } else {
      const created = await api.createServer({ ...form, id: crypto.randomUUID() } as ServerConfig)
      servers.update((list) => [...list, created])
    }
    showNewServer = false
    editingServer = null
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
  }

  function openEdit(s: ServerInfo) {
    editingServer = s
    resetForm(s)
    showNewServer = true
  }

  function collectionName(id: string): string {
    return $collections.find((c) => c.id === id)?.name ?? '(no collection)'
  }

  function toggleLog(id: string) {
    if (expandedLogServer === id) {
      expandedLogServer = null
    } else {
      expandedLogServer = id
      api.getRequestLog(id).then((log) => {
        requestLogs.update((r) => ({ ...r, [id]: log }))
      })
    }
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

<div class="flex flex-1 flex-col gap-4 overflow-y-auto p-5">
  <!-- Header -->
  <div class="flex shrink-0 items-center justify-between">
    <h2 class="text-lg font-semibold text-ink">Servers</h2>
    <button
      class="btn btn-primary"
      on:click={() => {
        resetForm()
        editingServer = null
        showNewServer = true
      }}>+ New Server</button
    >
  </div>

  {#if $servers.length === 0}
    <div class="flex flex-1 flex-col items-center justify-center gap-2.5 text-center">
      <div class="text-5xl opacity-[0.15]">◎</div>
      <p class="text-[13px] text-ink-mid">No servers configured.</p>
      <p class="max-w-xs text-xs text-ink-muted">
        Create a server and bind it to a collection to start mocking endpoints.
      </p>
      <button
        class="btn btn-primary"
        on:click={() => {
          resetForm()
          showNewServer = true
        }}>+ New Server</button
      >
    </div>
  {:else}
    <div class="flex flex-col gap-2.5">
      {#each $servers as server}
        <div
          class="overflow-hidden rounded-lg border border-wire bg-cave-surface transition-colors hover:border-wire-hi"
        >
          <!-- Main row -->
          <div class="flex items-center gap-4 px-4 py-3.5">
            <!-- Left: name + meta -->
            <div class="min-w-0 flex-1">
              <div class="text-sm font-semibold text-ink">{server.name || 'Unnamed Server'}</div>
              <div class="mt-0.5 flex flex-wrap items-center gap-2">
                <span class="font-mono text-base font-bold text-accent">:{server.port}</span>
                {#if server.https}
                  <span
                    class="rounded-[3px] bg-accent/[.12] px-1.5 py-0.5 font-mono text-[10px] font-bold tracking-[0.06em] text-accent"
                    >HTTPS</span
                  >
                {/if}
                <span class="text-xs text-ink-muted">{collectionName(server.collectionId)}</span>
              </div>
            </div>

            <!-- Center: status + error -->
            <div class="flex min-w-[100px] flex-col items-start gap-1">
              <span class="status-pill status-{server.status}">{server.status}</span>
              {#if server.errorMsg}
                <span
                  class="max-w-[200px] overflow-hidden text-ellipsis whitespace-nowrap text-[11px] text-err"
                  title={server.errorMsg}>⚠ {server.errorMsg}</span
                >
              {/if}
            </div>

            <!-- Right: request count + actions -->
            <div class="flex shrink-0 items-center gap-3">
              <div class="flex flex-col items-center">
                <span class="font-mono text-lg font-bold leading-none text-ink"
                  >{server.requestCount}</span
                >
                <span class="text-[10px] tracking-[0.05em] text-ink-muted">requests</span>
              </div>

              <div class="flex items-center gap-1.5">
                <button
                  class="btn {server.status === 'running'
                    ? 'btn-ghost'
                    : 'btn-primary'} min-w-[80px]"
                  on:click={() => toggleServer(server)}
                  disabled={server.status === 'error'}
                >
                  {server.status === 'running' ? '⏹ Stop' : '▶ Start'}
                </button>

                <button
                  class="btn btn-ghost {expandedLogServer === server.id
                    ? 'border-accent bg-accent/[.12] text-accent'
                    : ''}"
                  on:click={() => toggleLog(server.id)}>📋 Log</button
                >

                <button class="btn-icon" title="Edit" on:click={() => openEdit(server)}>✏</button>
                <button
                  class="btn-icon hover:text-err"
                  title="Delete"
                  on:click={() => deleteServer(server.id)}>🗑</button
                >
              </div>
            </div>
          </div>

          <!-- Log drawer -->
          {#if expandedLogServer === server.id}
            <div class="border-t border-wire bg-cave-deep">
              <div class="flex items-center justify-between border-b border-wire px-3 py-1.5">
                <span class="text-[10px] font-semibold uppercase tracking-[0.08em] text-ink-muted"
                  >Request Log</span
                >
                <button
                  class="btn btn-ghost px-2 py-0.5 text-[11px]"
                  on:click={() => clearLog(server.id)}>Clear</button
                >
              </div>

              <div class="max-h-[280px] overflow-y-auto py-1">
                {#if !$requestLogs[server.id] || $requestLogs[server.id].length === 0}
                  <div class="py-4 text-center text-xs text-ink-muted">No requests yet.</div>
                {:else}
                  {#each [...($requestLogs[server.id] ?? [])].reverse() as entry}
                    <div
                      class="flex items-center gap-2 border-b border-transparent px-3 py-1 text-[11px] transition-colors hover:bg-cave-raised
                        {entry.breakpointed ? 'border-l-2 border-l-accent pl-2.5' : ''}
                        {!entry.matched ? 'opacity-60' : ''}"
                    >
                      <span class="w-[72px] shrink-0 font-mono text-ink-muted"
                        >{formatTime(entry.timestamp)}</span
                      >
                      <span
                        class="tag method-badge shrink-0"
                        style="background:{mc(entry.method)};font-size:9px;padding:1px 5px"
                        >{entry.method}</span
                      >
                      <span class="min-w-0 flex-1 truncate font-mono text-ink-mid">
                        {entry.path}{entry.query ? '?' + entry.query : ''}
                      </span>
                      <span
                        class="w-8 shrink-0 text-right font-mono font-bold"
                        style="color:{statusColor(entry.statusCode)}">{entry.statusCode}</span
                      >
                      <span class="w-11 shrink-0 text-right font-mono text-ink-muted"
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
                      {#if !entry.matched}
                        <span
                          class="shrink-0 rounded-[3px] bg-err/10 px-1 py-0.5 font-mono text-[9px] text-err"
                          >404</span
                        >
                      {/if}
                    </div>
                  {/each}
                {/if}
              </div>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</div>

{#if showNewServer}
  <Modal
    title={editingServer ? 'Edit Server' : 'New Server'}
    on:close={() => {
      showNewServer = false
      editingServer = null
    }}
  >
    <div class="flex flex-col gap-3.5">
      <label class="flex flex-col gap-1.5">
        <span class="form-label">Name</span>
        <input bind:value={form.name} placeholder="My API Server" />
      </label>

      <label class="flex flex-col gap-1.5">
        <span class="form-label">Port</span>
        <input
          type="number"
          min="1"
          max="65535"
          bind:value={form.port}
          placeholder="8000"
          class="font-mono"
        />
      </label>

      <label class="flex flex-col gap-1.5">
        <span class="form-label">Collection</span>
        <select bind:value={form.collectionId}>
          <option value="">— No collection —</option>
          {#each $collections as c}
            <option value={c.id}>{c.name}</option>
          {/each}
        </select>
      </label>

      <div class="flex items-center justify-between">
        <span class="form-label">HTTPS (self-signed)</span>
        <input type="checkbox" bind:checked={form.https} />
      </div>

      <div class="mt-1 flex justify-end gap-2">
        <button class="btn btn-primary" on:click={saveServer}
          >{editingServer ? 'Save' : 'Create'}</button
        >
        <button
          class="btn btn-ghost"
          on:click={() => {
            showNewServer = false
            editingServer = null
          }}>Cancel</button
        >
      </div>
    </div>
  </Modal>
{/if}
