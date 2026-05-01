<script lang="ts">
    import { activeWorkspaceTab, pendingBreakpoints } from "../stores/index";
    import { TabStrip } from "../lib/ui";
    import Sidebar from "./Sidebar.svelte";
    import CollectionsPanel from "./CollectionsPanel.svelte";
    import ServersSidebar from "./ServersSidebar.svelte";
    import ServersPanel from "./ServersPanel.svelte";
    import BreakpointsPanel from "./BreakpointsPanel.svelte";
</script>

<div class="flex min-w-0 flex-1 flex-col overflow-hidden">
    <TabStrip
        bind:value={$activeWorkspaceTab}
        tabs={[
            { value: "collections", label: "Collections" },
            { value: "servers", label: "Servers" },
            {
                value: "breakpoints",
                label: "Breakpoints",
                badge: $pendingBreakpoints.length || null,
            },
        ]}
    />
    <div class="flex min-h-0 flex-1 overflow-hidden">
        {#if $activeWorkspaceTab === "collections"}
            <Sidebar />
            <CollectionsPanel />
        {:else if $activeWorkspaceTab === "servers"}
            <ServersSidebar />
            <ServersPanel />
        {:else if $activeWorkspaceTab === "breakpoints"}
            <BreakpointsPanel />
        {/if}
    </div>
</div>
