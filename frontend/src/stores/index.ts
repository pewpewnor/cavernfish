import { writable, derived } from 'svelte/store'
import type { Collection, ServerInfo, RequestLogEntry, BreakpointHit, Endpoint } from '../lib/types'

export const collections = writable<Collection[]>([])
export const servers = writable<ServerInfo[]>([])
export const requestLogs = writable<Record<string, RequestLogEntry[]>>({})
export const pendingBreakpoints = writable<BreakpointHit[]>([])

export type View = 'servers' | 'collections' | 'breakpoints' | 'settings'
export const activeView = writable<View>('servers')

export const expandedCollections = writable<Set<string>>(new Set())
export const expandedFolders = writable<Set<string>>(new Set())

export interface EndpointTab {
  tabId: string
  endpointId: string
  collectionId: string
  folderId: string
}

export const openTabs = writable<EndpointTab[]>([])
export const activeTabId = writable<string | null>(null)
export const dirtyTabs = writable<Set<string>>(new Set())
export const tabEdits = writable<Record<string, Endpoint>>({})

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
