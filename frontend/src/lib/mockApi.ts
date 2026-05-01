import type {
    Collection,
    Endpoint,
    ServerInfo,
    ServerConfig,
    RequestLogEntry,
    BreakpointHit,
    MockResponse,
} from "./types";

const COLLECTIONS_KEY = "cavernfish:mock:collections";
const SERVERS_KEY = "cavernfish:mock:servers";

function loadCollections(): Collection[] {
    try {
        return JSON.parse(localStorage.getItem(COLLECTIONS_KEY) || "[]");
    } catch {
        return [];
    }
}

function saveCollections(cs: Collection[]): void {
    localStorage.setItem(COLLECTIONS_KEY, JSON.stringify(cs));
}

function loadServers(): ServerInfo[] {
    try {
        return JSON.parse(localStorage.getItem(SERVERS_KEY) || "[]");
    } catch {
        return [];
    }
}

function saveServers(ss: ServerInfo[]): void {
    localStorage.setItem(SERVERS_KEY, JSON.stringify(ss));
    emitServerStatus();
}

const serverStatusListeners = new Set<(servers: ServerInfo[]) => void>();

function emitServerStatus(): void {
    const list = loadServers();
    queueMicrotask(() => {
        for (const cb of serverStatusListeners) cb(list);
    });
}

function now(): string {
    return new Date().toISOString();
}

function findCollection(cs: Collection[], id: string): number {
    return cs.findIndex((c) => c.id === id);
}

export const mock = {
    async getCollections(): Promise<Collection[]> {
        return loadCollections();
    },

    async createCollection(name: string, description: string): Promise<Collection> {
        const c: Collection = {
            id: crypto.randomUUID(),
            name,
            description,
            variables: {},
            folders: [],
            createdAt: now(),
            updatedAt: now(),
        };
        const cs = loadCollections();
        cs.push(c);
        saveCollections(cs);
        return c;
    },

    async updateCollection(
        id: string,
        upd: { name: string; description: string; variables: Record<string, string> },
    ): Promise<Collection> {
        const cs = loadCollections();
        const i = findCollection(cs, id);
        if (i === -1) throw new Error("collection not found");
        if (upd.name) cs[i].name = upd.name;
        cs[i].description = upd.description;
        if (upd.variables) cs[i].variables = upd.variables;
        cs[i].updatedAt = now();
        saveCollections(cs);
        return cs[i];
    },

    async deleteCollection(id: string): Promise<void> {
        saveCollections(loadCollections().filter((c) => c.id !== id));
    },

    async exportCollection(id: string): Promise<string> {
        const c = loadCollections().find((c) => c.id === id);
        if (!c) throw new Error("not found");
        return JSON.stringify(c, null, 2);
    },

    async importCollection(data: string): Promise<Collection> {
        const c = JSON.parse(data) as Collection;
        c.id = crypto.randomUUID();
        c.updatedAt = now();
        const cs = loadCollections();
        cs.push(c);
        saveCollections(cs);
        return c;
    },

    async importFromOpenAPI(_data: string): Promise<Collection> {
        throw new Error("OpenAPI import is unavailable in browser mode (requires the Go backend)");
    },

    async createFolder(collectionId: string, name: string): Promise<Collection> {
        const cs = loadCollections();
        const i = findCollection(cs, collectionId);
        if (i === -1) throw new Error("collection not found");
        cs[i].folders.push({
            id: crypto.randomUUID(),
            name,
            endpoints: [],
        });
        cs[i].updatedAt = now();
        saveCollections(cs);
        return cs[i];
    },

    async renameFolder(collectionId: string, folderId: string, name: string): Promise<Collection> {
        const cs = loadCollections();
        const i = findCollection(cs, collectionId);
        if (i === -1) throw new Error("collection not found");
        const f = cs[i].folders.find((f) => f.id === folderId);
        if (f) f.name = name;
        cs[i].updatedAt = now();
        saveCollections(cs);
        return cs[i];
    },

    async deleteFolder(collectionId: string, folderId: string): Promise<Collection> {
        const cs = loadCollections();
        const i = findCollection(cs, collectionId);
        if (i === -1) throw new Error("collection not found");
        cs[i].folders = cs[i].folders.filter((f) => f.id !== folderId);
        cs[i].updatedAt = now();
        saveCollections(cs);
        return cs[i];
    },

    async createEndpoint(
        collectionId: string,
        folderId: string,
        ep: Endpoint,
    ): Promise<Collection> {
        const cs = loadCollections();
        const i = findCollection(cs, collectionId);
        if (i === -1) throw new Error("collection not found");
        const f = cs[i].folders.find((f) => f.id === folderId);
        if (f) {
            const ne: Endpoint = { ...ep, id: ep.id || crypto.randomUUID() };
            f.endpoints.push(ne);
        }
        cs[i].updatedAt = now();
        saveCollections(cs);
        return cs[i];
    },

    async updateEndpoint(
        collectionId: string,
        folderId: string,
        endpointId: string,
        ep: Endpoint,
    ): Promise<Collection> {
        const cs = loadCollections();
        const i = findCollection(cs, collectionId);
        if (i === -1) throw new Error("collection not found");
        const f = cs[i].folders.find((f) => f.id === folderId);
        if (f) {
            const j = f.endpoints.findIndex((e) => e.id === endpointId);
            if (j !== -1) f.endpoints[j] = { ...ep, id: endpointId };
        }
        cs[i].updatedAt = now();
        saveCollections(cs);
        return cs[i];
    },

    async deleteEndpoint(
        collectionId: string,
        folderId: string,
        endpointId: string,
    ): Promise<Collection> {
        const cs = loadCollections();
        const i = findCollection(cs, collectionId);
        if (i === -1) throw new Error("collection not found");
        const f = cs[i].folders.find((f) => f.id === folderId);
        if (f) f.endpoints = f.endpoints.filter((e) => e.id !== endpointId);
        cs[i].updatedAt = now();
        saveCollections(cs);
        return cs[i];
    },

    async getServers(): Promise<ServerInfo[]> {
        return loadServers();
    },

    async createServer(cfg: ServerConfig): Promise<ServerInfo> {
        const s: ServerInfo = {
            id: cfg.id || crypto.randomUUID(),
            name: cfg.name,
            port: cfg.port,
            collectionIds: cfg.collectionIds,
            https: cfg.https,
            disabledEndpointIds: cfg.disabledEndpointIds,
            status: "stopped",
            requestCount: 0,
        };
        const ss = loadServers();
        ss.push(s);
        saveServers(ss);
        return s;
    },

    async updateServer(cfg: ServerConfig): Promise<ServerInfo> {
        const ss = loadServers();
        const i = ss.findIndex((s) => s.id === cfg.id);
        if (i === -1) throw new Error("server not found");
        ss[i] = { ...ss[i], ...cfg };
        saveServers(ss);
        return ss[i];
    },

    async isServerRunning(id: string): Promise<boolean> {
        return loadServers().find((s) => s.id === id)?.status === "running";
    },

    async startServer(id: string): Promise<void> {
        const ss = loadServers();
        const i = ss.findIndex((s) => s.id === id);
        if (i !== -1) {
            ss[i].status = "running";
            saveServers(ss);
        }
    },

    async stopServer(id: string): Promise<void> {
        const ss = loadServers();
        const i = ss.findIndex((s) => s.id === id);
        if (i !== -1) {
            ss[i].status = "stopped";
            saveServers(ss);
        }
    },

    async restartServer(id: string): Promise<void> {
        await mock.stopServer(id);
        await mock.startServer(id);
    },

    async deleteServer(id: string): Promise<void> {
        saveServers(loadServers().filter((s) => s.id !== id));
    },

    async getRequestLog(_id: string): Promise<RequestLogEntry[]> {
        return [];
    },

    async clearRequestLog(_id: string): Promise<void> {},

    async getPendingBreakpoints(): Promise<BreakpointHit[]> {
        return [];
    },

    async releaseBreakpoint(_id: string, _resp: MockResponse): Promise<boolean> {
        return false;
    },

    async discardBreakpoint(_id: string): Promise<boolean> {
        return false;
    },

    onServerStatus(cb: (servers: ServerInfo[]) => void): () => void {
        serverStatusListeners.add(cb);
        return () => serverStatusListeners.delete(cb);
    },
    onRequestLogged(_cb: (e: RequestLogEntry) => void): () => void {
        return () => {};
    },
    onBreakpointHit(_cb: (h: BreakpointHit) => void): () => void {
        return () => {};
    },
    onBreakpointReleased(_cb: (id: string) => void): () => void {
        return () => {};
    },
};
