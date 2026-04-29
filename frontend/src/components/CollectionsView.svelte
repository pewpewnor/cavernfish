<script lang="ts">
  import type { Endpoint } from '../lib/types'
  import { patchCollection, selectedEndpoint } from '../stores/index'
  import * as api from '../lib/api'
  import EndpointEditor from './EndpointEditor.svelte'

  async function handleSave(
    e: CustomEvent<{ collectionId: string; folderId: string; endpoint: Endpoint }>,
  ) {
    const { collectionId, folderId, endpoint } = e.detail
    const updated = await api.updateEndpoint(collectionId, folderId, endpoint.id, endpoint)
    patchCollection(updated)
  }
</script>

<div class="flex min-w-0 flex-1 flex-col overflow-hidden">
  {#if $selectedEndpoint}
    <EndpointEditor
      endpoint={$selectedEndpoint.endpoint}
      collectionId={$selectedEndpoint.collectionId}
      folderId={$selectedEndpoint.folderId}
      on:save={handleSave}
    />
  {:else}
    <div class="flex flex-1 flex-col items-center justify-center gap-3 p-10 text-center">
      <div class="relative mb-2 flex items-center justify-center">
        <span class="relative z-10 text-[52px] opacity-20">⬡</span>
        <div class="absolute h-24 w-24 animate-breathe rounded-full bg-accent/10" />
      </div>
      <h3 class="text-base font-semibold text-ink">Select an endpoint</h3>
      <p class="max-w-sm text-[13px] leading-relaxed text-ink-muted">
        Choose an endpoint from the sidebar to edit its mock response, or create a new one.
      </p>
      <div class="mt-3 flex flex-col items-center gap-2">
        {#each [['click folder +', 'to add a new endpoint'], ['⚡ breakpoint', 'to intercept live requests and edit on the fly'], ['cycle strategy', 'to rotate through multiple responses']] as [key, hint]}
          <div class="flex items-center gap-2 text-xs text-ink-muted">
            <kbd
              class="rounded-sm border border-wire-hi bg-cave-raised px-2 py-0.5 font-mono text-[11px] text-ink-mid"
              >{key}</kbd
            >
            <span>{hint}</span>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>
