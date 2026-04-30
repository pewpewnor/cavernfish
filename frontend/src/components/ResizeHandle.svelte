<script lang="ts">
  import { createEventDispatcher } from 'svelte'

  export let min = 160
  export let max = 520
  export let currentWidth = 240

  const dispatch = createEventDispatcher()
  let active = false

  function onMousedown(e: MouseEvent) {
    const startX = e.clientX
    const startW = currentWidth
    active = true
    e.preventDefault()

    function onMove(ev: MouseEvent) {
      dispatch('resize', Math.max(min, Math.min(max, startW + ev.clientX - startX)))
    }

    function onUp() {
      active = false
      window.removeEventListener('mousemove', onMove)
      window.removeEventListener('mouseup', onUp)
    }

    window.addEventListener('mousemove', onMove)
    window.addEventListener('mouseup', onUp)
  }
</script>

<div
  class="absolute right-0 top-0 z-10 h-full w-1.5 cursor-col-resize select-none transition-colors
    {active ? 'bg-accent/50' : 'hover:bg-accent/30'}"
  on:mousedown={onMousedown}
/>
