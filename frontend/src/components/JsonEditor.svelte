<script lang="ts">
    import { onDestroy, onMount, untrack } from "svelte";
    import { EditorView, basicSetup } from "codemirror";
    import { EditorState } from "@codemirror/state";
    import { json } from "@codemirror/lang-json";
    import { oneDark } from "@codemirror/theme-one-dark";

    interface Props {
        value: string;
        readOnly?: boolean;
        onChange?: (v: string) => void;
    }

    let { value, readOnly = false, onChange }: Props = $props();

    let host: HTMLDivElement;
    let view: EditorView | null = null;
    // svelte-ignore state_referenced_locally
    let lastEmitted = value;
    let io: IntersectionObserver | null = null;

    onMount(() => {
        const heightTheme = EditorView.theme({
            "&": { height: "100%", fontSize: "0.8125rem" },
            ".cm-scroller": { fontSize: "0.8125rem" },
            ".cm-gutters": { fontSize: "0.8125rem" },
        });

        view = new EditorView({
            doc: value,
            extensions: [
                basicSetup,
                json(),
                oneDark,
                heightTheme,
                EditorState.readOnly.of(readOnly),
                EditorView.updateListener.of((update) => {
                    if (update.docChanged) {
                        const next = update.state.doc.toString();
                        lastEmitted = next;
                        onChange?.(next);
                    }
                }),
            ],
            parent: host,
        });

        io = new IntersectionObserver((entries) => {
            for (const e of entries) {
                if (e.isIntersecting) view?.requestMeasure();
            }
        });
        io.observe(host);
    });

    onDestroy(() => {
        io?.disconnect();
        view?.destroy();
    });

    $effect(() => {
        const v = value;
        untrack(() => {
            if (!view) return;
            if (v === lastEmitted) return;
            if (v === view.state.doc.toString()) return;
            lastEmitted = v;
            view.dispatch({
                changes: { from: 0, to: view.state.doc.length, insert: v },
            });
        });
    });
</script>

<div bind:this={host} class="h-full w-full overflow-hidden rounded-sm border border-wire"></div>
