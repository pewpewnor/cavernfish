<script lang="ts">
  import { createEventDispatcher, onDestroy, onMount } from 'svelte'
  import { EditorView, basicSetup } from 'codemirror'
  import { json } from '@codemirror/lang-json'
  import { oneDark } from '@codemirror/theme-one-dark'

  export let value: string = ''
  export let readOnly: boolean = false
  export let minHeight: string = '200px'

  const dispatch = createEventDispatcher()

  let el: HTMLDivElement
  let view: EditorView
  let internalUpdate = false

  const cavernTheme = EditorView.theme({
    '&': { background: '#1e2236', width: '100%', flex: '1' },
    '.cm-scroller': {
      fontFamily: 'JetBrains Mono, Fira Code, Cascadia Code, ui-monospace, monospace',
      fontSize: '13px',
      minHeight,
      overflow: 'auto',
      width: '100%',
    },
    '.cm-content': { padding: '10px 0', width: '100%' },
    '.cm-line': { padding: '0 12px' },
    '.cm-focused': { outline: 'none' },
    '.cm-gutters': { background: '#191c2a', borderRight: '1px solid #2c3048', color: '#636f8f' },
    '&.cm-focused .cm-cursor': { borderLeftColor: '#3b6ef0' },
    '.cm-selectionBackground': { background: 'rgba(59,110,240,0.18) !important' },
    '&.cm-focused .cm-selectionBackground': { background: 'rgba(59,110,240,0.12) !important' },
    '.cm-editor': { width: '100%' },
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
            dispatch('change', value)
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
  class="flex w-full flex-1 flex-col overflow-hidden rounded-sm border border-wire transition-colors focus-within:border-accent focus-within:ring-2 focus-within:ring-accent/20"
/>
