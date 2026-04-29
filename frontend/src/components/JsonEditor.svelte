<script lang="ts">
  import { onDestroy, onMount } from 'svelte'
  import { EditorView, basicSetup } from 'codemirror'
  import { json } from '@codemirror/lang-json'
  import { oneDark } from '@codemirror/theme-one-dark'

  export let value: string = ''
  export let readOnly: boolean = false
  export let minHeight: string = '120px'

  let el: HTMLDivElement
  let view: EditorView
  let internalUpdate = false

  const cavernTheme = EditorView.theme({
    '&': { background: '#1c2035', height: '100%' },
    '.cm-scroller': { fontFamily: 'var(--font-mono, monospace)', fontSize: '12px', minHeight },
    '.cm-content': { padding: '8px 0' },
    '.cm-line': { padding: '0 10px' },
    '.cm-focused': { outline: 'none' },
    '.cm-gutters': { background: '#141724', borderRight: '1px solid #252a40', color: '#4a5368' },
    '&.cm-focused .cm-cursor': { borderLeftColor: '#3b6ef0' },
    '.cm-selectionBackground': { background: 'rgba(59,110,240,0.18) !important' },
    '&.cm-focused .cm-selectionBackground': { background: 'rgba(59,110,240,0.12) !important' },
  })

  onMount(() => {
    view = new EditorView({
      doc: value,
      extensions: [
        basicSetup,
        json(),
        oneDark,
        cavernTheme,
        EditorView.editable.of(!readOnly),
        EditorView.updateListener.of((update) => {
          if (update.docChanged && !internalUpdate) {
            value = update.state.doc.toString()
          }
        }),
      ],
      parent: el,
    })
  })

  onDestroy(() => view?.destroy())

  $: if (view) {
    const current = view.state.doc.toString()
    if (value !== current) {
      internalUpdate = true
      view.dispatch({ changes: { from: 0, to: view.state.doc.length, insert: value } })
      internalUpdate = false
    }
  }
</script>

<div
  bind:this={el}
  class="flex-1 overflow-hidden rounded-sm border border-wire transition-colors focus-within:border-accent focus-within:ring-2 focus-within:ring-accent/20"
/>

<style>
  :global(div.cm-editor) {
    height: 100%;
  }
</style>
