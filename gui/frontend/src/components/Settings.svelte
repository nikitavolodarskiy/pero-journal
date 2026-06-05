<script>
  export let config    = { JournalDir: '', Editor: '' }
  export let onSave    = async (_journalDir, _editor) => {}
  export let onClose   = () => {}

  let journalDir = config.JournalDir
  let editor     = config.Editor
  let saving     = false
  let errorMsg   = ''

  async function save() {
    saving   = true
    errorMsg = ''
    try {
      await onSave(journalDir, editor)
      onClose()
    } catch (e) {
      errorMsg = String(e)
    } finally {
      saving = false
    }
  }

  function handleKeydown(e) {
    if (e.key === 'Escape') onClose()
    if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) save()
  }
</script>

<svelte:window on:keydown={handleKeydown} />

<!-- clicking the backdrop closes the panel -->
<div class="overlay" on:click|self={onClose} role="dialog" aria-modal="true" aria-label="Settings">
  <div class="panel">
    <h2 class="panel-title">Settings</h2>

    <div class="field">
      <label for="journal-dir">Journal directory</label>
      <input
        id="journal-dir"
        bind:value={journalDir}
        placeholder="~/pero"
        spellcheck="false"
        autocomplete="off"
      />
      <p class="hint">Where your .md entry files are stored.</p>
    </div>

    <div class="field">
      <label for="editor">CLI editor</label>
      <input
        id="editor"
        bind:value={editor}
        placeholder="nvim"
        spellcheck="false"
        autocomplete="off"
      />
      <p class="hint">Used by <code>pero</code> / <code>pero open</code>.</p>
    </div>

    {#if errorMsg}
      <p class="error-msg">{errorMsg}</p>
    {/if}

    <div class="actions">
      <button class="btn-cancel" on:click={onClose} disabled={saving}>Cancel</button>
      <button class="btn-save"   on:click={save}    disabled={saving}>
        {saving ? 'Saving…' : 'Save'}
      </button>
    </div>
  </div>
</div>

<style>
  .overlay {
    position:        fixed;
    inset:           0;
    background:      rgba(0, 0, 0, 0.55);
    display:         flex;
    align-items:     center;
    justify-content: center;
    z-index:         100;
  }

  .panel {
    background:    #252526;
    border:        1px solid #3c3c3c;
    border-radius: 8px;
    padding:       24px 28px;
    width:         380px;
    max-width:     calc(100vw - 48px);
    box-shadow:    0 8px 32px rgba(0, 0, 0, 0.5);
  }

  .panel-title {
    margin:      0 0 20px;
    font-size:   14px;
    font-weight: 600;
    color:       #ccc;
    letter-spacing: 0.02em;
  }

  .field {
    margin-bottom: 16px;
  }

  label {
    display:       block;
    font-size:     11px;
    font-weight:   600;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color:         #666;
    margin-bottom: 6px;
  }

  input {
    width:         100%;
    box-sizing:    border-box;
    background:    #1e1e1e;
    border:        1px solid #3c3c3c;
    border-radius: 4px;
    color:         #d4d4d4;
    font-size:     13px;
    font-family:   ui-monospace, 'Cascadia Code', 'Menlo', monospace;
    padding:       7px 10px;
    outline:       none;
    transition:    border-color 0.15s;
  }

  input:focus {
    border-color: #5a7db5;
  }

  .hint {
    margin:    5px 0 0;
    font-size: 11px;
    color:     #555;
    line-height: 1.4;
  }

  .hint code {
    font-family: ui-monospace, 'Cascadia Code', 'Menlo', monospace;
    color:       #7f9cf5;
    font-size:   11px;
  }

  .error-msg {
    margin:      0 0 14px;
    font-size:   12px;
    color:       #f87171;
    background:  rgba(248, 113, 113, 0.08);
    border-radius: 4px;
    padding:     8px 10px;
  }

  .actions {
    display:         flex;
    justify-content: flex-end;
    gap:             8px;
    margin-top:      20px;
  }

  button {
    border:        none;
    border-radius: 4px;
    font-size:     12px;
    font-weight:   500;
    padding:       6px 16px;
    cursor:        pointer;
    transition:    background 0.15s, color 0.15s;
  }

  button:disabled {
    opacity: 0.5;
    cursor:  default;
  }

  .btn-cancel {
    background: transparent;
    color:      #888;
    border:     1px solid #3c3c3c;
  }

  .btn-cancel:hover:not(:disabled) {
    background: #2d2d2d;
    color:      #ccc;
  }

  .btn-save {
    background: #5a7db5;
    color:      #fff;
  }

  .btn-save:hover:not(:disabled) {
    background: #6a8dc5;
  }
</style>
