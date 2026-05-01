import type {
    Collection,
    ServerInfo,
    ServerConfig,
    RequestLogEntry,
    BreakpointHit,
    Endpoint,
    MockResponse,
} from "./types";

import * as App from "../../wailsjs/go/main/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { mock } from "./mockApi";

const cast = <T>(p: Promise<unknown>): Promise<T> => p as Promise<T>;

export function isWailsReady(): boolean {
    return (
        typeof (window as unknown as { go?: { main?: { App?: unknown } } }).go?.main?.App !==
        "undefined"
    );
}

export async function waitForWails(timeoutMs = 2000): Promise<boolean> {
    if (isWailsReady()) return true;
    const start = Date.now();
    while (Date.now() - start < timeoutMs) {
        await new Promise((r) => setTimeout(r, 50));
        if (isWailsReady()) return true;
    }
    return false;
}

export function isMockMode(): boolean {
    return !isWailsReady();
}

export const getCollections = (): Promise<Collection[]> =>
    isMockMode() ? mock.getCollections() : cast(App.GetCollections());

export const createCollection = (name: string, description: string): Promise<Collection> =>
    isMockMode()
        ? mock.createCollection(name, description)
        : cast(App.CreateCollection(name, description));

export const updateCollection = (
    id: string,
    upd: { name: string; description: string; variables: Record<string, string> },
): Promise<Collection> =>
    isMockMode() ? mock.updateCollection(id, upd) : cast(App.UpdateCollection(id, upd));

export const deleteCollection = (id: string): Promise<void> =>
    isMockMode() ? mock.deleteCollection(id) : App.DeleteCollection(id);

export const exportCollection = (id: string): Promise<string> =>
    isMockMode() ? mock.exportCollection(id) : App.ExportCollection(id);

export const importCollection = (data: string): Promise<Collection> =>
    isMockMode() ? mock.importCollection(data) : cast(App.ImportCollection(data));

export const importFromOpenAPI = (data: string): Promise<Collection> =>
    isMockMode() ? mock.importFromOpenAPI(data) : cast(App.ImportFromOpenAPI(data));

export const createFolder = (collectionID: string, name: string): Promise<Collection> =>
    isMockMode()
        ? mock.createFolder(collectionID, name)
        : cast(App.CreateFolder(collectionID, name));

export const renameFolder = (
    collectionID: string,
    folderID: string,
    name: string,
): Promise<Collection> =>
    isMockMode()
        ? mock.renameFolder(collectionID, folderID, name)
        : cast(App.RenameFolder(collectionID, folderID, name));

export const deleteFolder = (collectionID: string, folderID: string): Promise<Collection> =>
    isMockMode()
        ? mock.deleteFolder(collectionID, folderID)
        : cast(App.DeleteFolder(collectionID, folderID));

export const createEndpoint = (
    collectionID: string,
    folderID: string,
    ep: Endpoint,
): Promise<Collection> =>
    isMockMode()
        ? mock.createEndpoint(collectionID, folderID, ep)
        : cast(App.CreateEndpoint(collectionID, folderID, ep as never));

export const updateEndpoint = (
    collectionID: string,
    folderID: string,
    endpointID: string,
    ep: Endpoint,
): Promise<Collection> =>
    isMockMode()
        ? mock.updateEndpoint(collectionID, folderID, endpointID, ep)
        : cast(App.UpdateEndpoint(collectionID, folderID, endpointID, ep as never));

export const deleteEndpoint = (
    collectionID: string,
    folderID: string,
    endpointID: string,
): Promise<Collection> =>
    isMockMode()
        ? mock.deleteEndpoint(collectionID, folderID, endpointID)
        : cast(App.DeleteEndpoint(collectionID, folderID, endpointID));

export const getServers = (): Promise<ServerInfo[]> =>
    isMockMode() ? mock.getServers() : cast(App.GetServers());

export const createServer = (cfg: ServerConfig): Promise<ServerInfo> =>
    isMockMode() ? mock.createServer(cfg) : cast(App.CreateServer(cfg));

export const updateServer = (cfg: ServerConfig): Promise<ServerInfo> =>
    isMockMode() ? mock.updateServer(cfg) : cast(App.UpdateServer(cfg));

export const isServerRunning = (id: string): Promise<boolean> =>
    isMockMode() ? mock.isServerRunning(id) : App.IsServerRunning(id);

export const startServer = (id: string): Promise<void> =>
    isMockMode() ? mock.startServer(id) : App.StartServer(id);

export const stopServer = (id: string): Promise<void> =>
    isMockMode() ? mock.stopServer(id) : App.StopServer(id);

export const restartServer = (id: string): Promise<void> =>
    isMockMode() ? mock.restartServer(id) : App.RestartServer(id);

export const deleteServer = (id: string): Promise<void> =>
    isMockMode() ? mock.deleteServer(id) : App.DeleteServer(id);

export const getRequestLog = (serverID: string): Promise<RequestLogEntry[]> =>
    isMockMode() ? mock.getRequestLog(serverID) : cast(App.GetRequestLog(serverID));

export const clearRequestLog = (serverID: string): Promise<void> =>
    isMockMode() ? mock.clearRequestLog(serverID) : App.ClearRequestLog(serverID);

export const getPendingBreakpoints = (): Promise<BreakpointHit[]> =>
    isMockMode() ? mock.getPendingBreakpoints() : cast(App.GetPendingBreakpoints());

export const releaseBreakpoint = (id: string, resp: MockResponse): Promise<boolean> =>
    isMockMode() ? mock.releaseBreakpoint(id, resp) : App.ReleaseBreakpoint(id, resp as never);

export const discardBreakpoint = (id: string): Promise<boolean> =>
    isMockMode() ? mock.discardBreakpoint(id) : App.DiscardBreakpoint(id);

export const onServerStatus = (cb: (servers: ServerInfo[]) => void) =>
    isMockMode() ? mock.onServerStatus(cb) : EventsOn("server:status", cb as (s: unknown) => void);

export const onRequestLogged = (cb: (entry: RequestLogEntry) => void) =>
    isMockMode()
        ? mock.onRequestLogged(cb)
        : EventsOn("request:logged", cb as (e: unknown) => void);

export const onBreakpointHit = (cb: (hit: BreakpointHit) => void) =>
    isMockMode()
        ? mock.onBreakpointHit(cb)
        : EventsOn("breakpoint:hit", cb as (h: unknown) => void);

export const onBreakpointReleased = (cb: (id: string) => void) =>
    isMockMode() ? mock.onBreakpointReleased(cb) : EventsOn("breakpoint:released", cb);
