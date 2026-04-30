export interface KVPair {
  key: string
  value: string
  enabled: boolean
}

export interface MockCookie {
  name: string
  value: string
  path: string
  domain: string
  maxAge: number
  httpOnly: boolean
  secure: boolean
  sameSite: 'Strict' | 'Lax' | 'None'
}

export interface MockResponse {
  id: string
  name: string
  statusCode: number
  body: string
  bodyType: 'json' | 'text' | 'html' | 'xml' | 'none'
  headers: KVPair[]
  cookies: MockCookie[]
}

export interface Endpoint {
  id: string
  name: string
  method: string
  path: string
  responses: MockResponse[]
  strategy: 'fixed' | 'cycle' | 'random'
  activeIdx: number
  delayMs: number
  breakpoint: boolean
  proxyUrl?: string
  wsEnabled: boolean
}

export interface Folder {
  id: string
  name: string
  endpoints: Endpoint[]
}

export interface Collection {
  id: string
  name: string
  description: string
  variables: Record<string, string>
  folders: Folder[]
  createdAt: string
  updatedAt: string
}

export interface ServerInfo {
  id: string
  name: string
  port: number
  collectionIds: string[]
  https: boolean
  status: 'running' | 'stopped' | 'error'
  errorMsg?: string
  requestCount: number
  disabledEndpointIds: string[]
}

export interface ServerConfig {
  id: string
  name: string
  port: number
  collectionIds: string[]
  https: boolean
  disabledEndpointIds: string[]
}

export interface RequestLogEntry {
  id: string
  serverId: string
  timestamp: string
  method: string
  path: string
  query: string
  reqHeaders: Record<string, string>
  reqBody: string
  statusCode: number
  latencyMs: number
  breakpointed: boolean
  breakpointId?: string
  isProxy: boolean
  isWs: boolean
  matched: boolean
}

export interface BreakpointHit {
  id: string
  serverId: string
  timestamp: string
  method: string
  path: string
  query: string
  reqHeaders: Record<string, string>
  reqBody: string
  endpoint: Endpoint
  response: MockResponse
}

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE' | 'HEAD' | 'OPTIONS' | '*'

export const METHOD_COLORS: Record<string, string> = {
  GET: '#22d3a0',
  POST: '#3b82f6',
  PUT: '#f59e0b',
  PATCH: '#a78bfa',
  DELETE: '#f43f5e',
  HEAD: '#64748b',
  OPTIONS: '#64748b',
  '*': '#64748b',
}

export const STATUS_COLORS: Record<number, string> = {}

export function statusColor(code: number): string {
  if (code >= 500) return '#f43f5e'
  if (code >= 400) return '#f59e0b'
  if (code >= 300) return '#3b82f6'
  if (code >= 200) return '#22d3a0'
  return '#64748b'
}

export function newMockResponse(): MockResponse {
  return {
    id: crypto.randomUUID(),
    name: '200 OK',
    statusCode: 200,
    body: '{\n  "message": "OK"\n}',
    bodyType: 'json',
    headers: [],
    cookies: [],
  }
}

export function newKVPair(): KVPair {
  return { key: '', value: '', enabled: true }
}

export function newCookie(): MockCookie {
  return {
    name: '',
    value: '',
    path: '/',
    domain: '',
    maxAge: 0,
    httpOnly: false,
    secure: false,
    sameSite: 'Lax',
  }
}
