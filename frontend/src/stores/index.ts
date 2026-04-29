import { writable, derived } from 'svelte/store'
import type {
  Collection,
  ServerInfo,
  RequestLogEntry,
  BreakpointHit,
  Endpoint,
  Folder,
} from '../lib/types'

export const collections = writable<Collection[]>([])
export const servers = writable<ServerInfo[]>([])
export const requestLogs = writable<Record<string, RequestLogEntry[]>>({})
export const pendingBreakpoints = writable<BreakpointHit[]>([])

export type View = 'servers' | 'collections' | 'settings'
export const activeView = writable<View>('servers')

export const selectedCollectionId = writable<string | null>(null)
export const selectedFolderId = writable<string | null>(null)
export const selectedEndpointId = writable<string | null>(null)

export const expandedCollections = writable<Set<string>>(new Set())
export const expandedFolders = writable<Set<string>>(new Set())

export const selectedEndpoint = derived(
  [collections, selectedCollectionId, selectedFolderId, selectedEndpointId],
  ([$cols, $cid, $fid, $eid]) => {
    if (!$cid || !$fid || !$eid) return null
    const col = $cols.find((c) => c.id === $cid)
    if (!col) return null
    const folder = col.folders.find((f) => f.id === $fid)
    if (!folder) return null
    const ep = folder.endpoints.find((e) => e.id === $eid)
    return ep ? { collectionId: $cid, folderId: $fid, endpoint: ep } : null
  },
)

export const selectedCollection = derived(
  [collections, selectedCollectionId],
  ([$cols, $cid]) => $cols.find((c) => c.id === $cid) ?? null,
)

// Patch a collection in the store after any mutation
export function patchCollection(updated: Collection) {
  collections.update((list) => {
    const idx = list.findIndex((c) => c.id === updated.id)
    if (idx === -1) return [...list, updated]
    const next = [...list]
    next[idx] = updated
    return next
  })
}

export function addLog(entry: RequestLogEntry) {
  requestLogs.update((logs) => {
    const prev = logs[entry.serverId] ?? []
    return { ...logs, [entry.serverId]: [...prev, entry].slice(-500) }
  })
}
