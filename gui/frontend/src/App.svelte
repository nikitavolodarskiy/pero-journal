<script>
  import { onMount } from 'svelte'
  import Sidebar  from './components/Sidebar.svelte'
  import Editor   from './components/Editor.svelte'
  import Settings from './components/Settings.svelte'

  import { GetToday, GetEntries, ReadEntry, WriteEntry, GetStats, GetConfig, SaveConfig } from '../wailsjs/go/main/App.js'
  import { WindowSetTitle } from '../wailsjs/runtime/runtime.js'

  let entries      = []
  let stats        = null
  let activeDate   = ''
  let content      = ''
  let saveStatus   = ''   // '' | 'saving' | 'saved' | 'error'
  let saveTimer    = null
  let isLoading    = false

  let showSettings = false
  let appConfig    = null

  // ── Startup ──────────────────────────────────────────────────────────────

  onMount(async () => {
    const today = await GetToday()   // creates today's file first
    await loadEntries()              // now today's entry is in the list
    await loadStats()
    await openEntry(today.Date.slice(0, 10))
    appConfig = await GetConfig()
  })

  // ── Data loading ─────────────────────────────────────────────────────────

  async function loadEntries() {
    try { entries = await GetEntries() }
    catch (e) { console.error('loadEntries:', e) }
  }

  async function loadStats() {
    try { stats = await GetStats() }
    catch (e) { console.error('loadStats:', e) }
  }

  // Reload entries + stats from disk (used by the refresh button).
  async function refreshEntries() {
    await loadEntries()
    await loadStats()
  }

  // Open (or create) today's entry, then refresh the list so it appears in
  // the sidebar if this is the first entry of the day.
  async function todayEntry() {
    const today = await GetToday()
    await loadEntries()
    await openEntry(today.Date.slice(0, 10))
  }

  async function openEntry(date) {
    if (date === activeDate) return
    if (activeDate) await saveNow()
    activeDate = date
    isLoading = true
    content   = ''          // prevent stale-content save during the await below
    // Show the current entry date in the native title bar
    WindowSetTitle(formatTitle(date))
    try   { content = await ReadEntry(date) }
    catch (e) { console.error('openEntry:', e); content = '' }
    finally   { isLoading = false }
  }

  // ── Save ─────────────────────────────────────────────────────────────────

  $: if (content !== undefined && activeDate && !isLoading) scheduleSave()

  function scheduleSave() {
    clearTimeout(saveTimer)
    saveTimer = setTimeout(saveNow, 1000)
  }

  async function saveNow() {
    if (!activeDate) return
    clearTimeout(saveTimer)
    saveStatus = 'saving'
    try {
      await WriteEntry(activeDate, content)
      saveStatus = 'saved'
      await loadStats()
      setTimeout(() => { saveStatus = '' }, 2000)
    } catch (e) {
      console.error('saveNow:', e)
      saveStatus = 'error'
      setTimeout(() => { saveStatus = '' }, 3000)
    }
  }

  async function saveSettings(journalDir, editor) {
    await SaveConfig(journalDir, editor)
    appConfig = await GetConfig()
    // Cancel any pending autosave before the directory changes under it.
    clearTimeout(saveTimer)
    activeDate = ''
    content    = ''
    // Same boot sequence as onMount: create today's file, load list, open it.
    const today = await GetToday()
    await loadEntries()
    await loadStats()
    await openEntry(today.Date.slice(0, 10))
  }

  function formatTitle(dateStr) {
    if (!dateStr) return 'Pero'
    const d = new Date(dateStr + 'T00:00:00')
    return d.toLocaleDateString('en-US', {
      weekday: 'long', month: 'long', day: 'numeric', year: 'numeric'
    })
  }
</script>

<div class="app">

  <Sidebar
    {entries}
    {stats}
    {activeDate}
    onSelect={openEntry}
    onRefresh={refreshEntries}
    onToday={todayEntry}
    onOpenSettings={() => showSettings = true}
  />

  {#if showSettings && appConfig}
    <Settings
      config={appConfig}
      onSave={saveSettings}
      onClose={() => showSettings = false}
    />
  {/if}

  <main class="main">
    {#if saveStatus === 'saving'}
      <div class="save-banner">Saving…</div>
    {:else if saveStatus === 'saved'}
      <div class="save-banner saved">Saved</div>
    {:else if saveStatus === 'error'}
      <div class="save-banner error">Save failed</div>
    {/if}

    <div class="editor-area">
      <Editor bind:content onSave={saveNow} />
    </div>
  </main>

</div>

<style>
  :global(*, *::before, *::after) { box-sizing: border-box; }

  :global(body) {
    margin:      0;
    padding:     0;
    background:  #1e1e1e;
    color:       #d4d4d4;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
    overflow:    hidden;
  }

  .app {
    display:  flex;
    height:   100vh;
    width:    100vw;
    overflow: hidden;
  }

  /* ── Main area ── */
  .main {
    display:        flex;
    flex-direction: column;
    flex:           1;
    overflow:       hidden;
    position:       relative;
  }

  /* ── Save indicator ── */
  .save-banner {
    position:   absolute;
    top:        8px;
    right:      16px;
    font-size:  11px;
    color:      #555;
    z-index:    10;
    transition: color 0.3s;
    pointer-events: none;
  }
  .save-banner.saved  { color: #4ec9b0; }
  .save-banner.error  { color: #f87171; }

  /* ── Editor area ── */
  .editor-area {
    flex:     1;
    overflow: hidden;
  }
</style>
