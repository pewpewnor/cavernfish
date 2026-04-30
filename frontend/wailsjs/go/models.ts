export namespace backend {
	
	export class MockCookie {
	    name: string;
	    value: string;
	    path: string;
	    domain: string;
	    maxAge: number;
	    httpOnly: boolean;
	    secure: boolean;
	    sameSite: string;
	
	    static createFrom(source: any = {}) {
	        return new MockCookie(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	        this.path = source["path"];
	        this.domain = source["domain"];
	        this.maxAge = source["maxAge"];
	        this.httpOnly = source["httpOnly"];
	        this.secure = source["secure"];
	        this.sameSite = source["sameSite"];
	    }
	}
	export class KVPair {
	    key: string;
	    value: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new KVPair(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
	        this.enabled = source["enabled"];
	    }
	}
	export class MockResponse {
	    id: string;
	    name: string;
	    statusCode: number;
	    body: string;
	    bodyType: string;
	    headers: KVPair[];
	    cookies: MockCookie[];
	
	    static createFrom(source: any = {}) {
	        return new MockResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.statusCode = source["statusCode"];
	        this.body = source["body"];
	        this.bodyType = source["bodyType"];
	        this.headers = this.convertValues(source["headers"], KVPair);
	        this.cookies = this.convertValues(source["cookies"], MockCookie);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Endpoint {
	    id: string;
	    name: string;
	    method: string;
	    path: string;
	    responses: MockResponse[];
	    strategy: string;
	    activeIdx: number;
	    delayMs: number;
	    breakpoint: boolean;
	    proxyUrl?: string;
	    wsEnabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Endpoint(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.method = source["method"];
	        this.path = source["path"];
	        this.responses = this.convertValues(source["responses"], MockResponse);
	        this.strategy = source["strategy"];
	        this.activeIdx = source["activeIdx"];
	        this.delayMs = source["delayMs"];
	        this.breakpoint = source["breakpoint"];
	        this.proxyUrl = source["proxyUrl"];
	        this.wsEnabled = source["wsEnabled"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class BreakpointHit {
	    id: string;
	    serverId: string;
	    // Go type: time
	    timestamp: any;
	    method: string;
	    path: string;
	    query: string;
	    reqHeaders: Record<string, string>;
	    reqBody: string;
	    endpoint: Endpoint;
	    response: MockResponse;
	
	    static createFrom(source: any = {}) {
	        return new BreakpointHit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.serverId = source["serverId"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.method = source["method"];
	        this.path = source["path"];
	        this.query = source["query"];
	        this.reqHeaders = source["reqHeaders"];
	        this.reqBody = source["reqBody"];
	        this.endpoint = this.convertValues(source["endpoint"], Endpoint);
	        this.response = this.convertValues(source["response"], MockResponse);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Folder {
	    id: string;
	    name: string;
	    endpoints: Endpoint[];
	
	    static createFrom(source: any = {}) {
	        return new Folder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.endpoints = this.convertValues(source["endpoints"], Endpoint);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Collection {
	    id: string;
	    name: string;
	    description: string;
	    variables: Record<string, string>;
	    folders: Folder[];
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Collection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.variables = source["variables"];
	        this.folders = this.convertValues(source["folders"], Folder);
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CollectionUpdate {
	    name: string;
	    description: string;
	    variables: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new CollectionUpdate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.variables = source["variables"];
	    }
	}
	
	
	
	
	
	export class RequestLogEntry {
	    id: string;
	    serverId: string;
	    // Go type: time
	    timestamp: any;
	    method: string;
	    path: string;
	    query: string;
	    reqHeaders: Record<string, string>;
	    reqBody: string;
	    statusCode: number;
	    latencyMs: number;
	    breakpointed: boolean;
	    breakpointId?: string;
	    isProxy: boolean;
	    isWs: boolean;
	    matched: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RequestLogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.serverId = source["serverId"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.method = source["method"];
	        this.path = source["path"];
	        this.query = source["query"];
	        this.reqHeaders = source["reqHeaders"];
	        this.reqBody = source["reqBody"];
	        this.statusCode = source["statusCode"];
	        this.latencyMs = source["latencyMs"];
	        this.breakpointed = source["breakpointed"];
	        this.breakpointId = source["breakpointId"];
	        this.isProxy = source["isProxy"];
	        this.isWs = source["isWs"];
	        this.matched = source["matched"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ServerConfig {
	    id: string;
	    name: string;
	    port: number;
	    collectionIds: string[];
	    https: boolean;
	    disabledEndpointIds: string[];
	
	    static createFrom(source: any = {}) {
	        return new ServerConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.port = source["port"];
	        this.collectionIds = source["collectionIds"];
	        this.https = source["https"];
	        this.disabledEndpointIds = source["disabledEndpointIds"];
	    }
	}
	export class ServerInfo {
	    id: string;
	    name: string;
	    port: number;
	    collectionIds: string[];
	    https: boolean;
	    status: string;
	    errorMsg?: string;
	    requestCount: number;
	    disabledEndpointIds: string[];
	
	    static createFrom(source: any = {}) {
	        return new ServerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.port = source["port"];
	        this.collectionIds = source["collectionIds"];
	        this.https = source["https"];
	        this.status = source["status"];
	        this.errorMsg = source["errorMsg"];
	        this.requestCount = source["requestCount"];
	        this.disabledEndpointIds = source["disabledEndpointIds"];
	    }
	}

}

