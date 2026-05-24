<script>
  import { onMount, onDestroy } from 'svelte'
  import { EditorView, keymap, lineNumbers, drawSelection } from '@codemirror/view'
  import { EditorState } from '@codemirror/state'
  import { defaultKeymap, historyKeymap, history } from '@codemirror/commands'
  import { markdown, markdownLanguage } from '@codemirror/lang-markdown'
  import { languages } from '@codemirror/language-data'
  import { syntaxHighlighting } from '@codemirror/language'
  import { livePreviewPlugin, livePreviewTheme } from '../lib/livePreview.js'
  import { markdownHighlightStyle } from '../lib/highlight.js'

  export let content = ''
  export let onSave  = () => {}

  let editorEl
  let view

  // Dark theme matching the app shell
  const darkTheme = EditorView.theme({
    '&': {
      height:     '100%',
      background: '#1e1e1e',
      color:      '#d4d4d4',
      fontSize:   '15px',
    },
    '.cm-scroller': {
      fontFamily:  '"iA Writer Quattro", "iA Writer Mono", ui-monospace, monospace',
      lineHeight:  '1.75',
      padding:     '0 0 40px 0',
      overflowY: 'auto',
    },
    '.cm-content': {
      maxWidth:  '680px',
      margin:    '0 auto',
      padding:   '40px 24px',
      caretColor: '#7f9cf5',
    },
    // Hide the cursor line highlight — too distracting for prose
    '.cm-activeLine':       { background: 'transparent' },
    '.cm-activeLineGutter': { background: 'transparent' },
    // Selection
    '.cm-selectionBackground, ::selection': { background: '#264f78 !important' },
    // No border on focused editor
    '&.cm-focused': { outline: 'none' },
    // Cursor
    '.cm-cursor': { borderLeftColor: '#7f9cf5' },
  }, { dark: true })

  // Save on Ctrl/Cmd+S
  const saveKeymap = keymap.of([{
    key: 'Mod-s',
    run() { onSave(view.state.doc.toString()); return true },
  }])

  onMount(() => {
    const state = EditorState.create({
      doc: content,
      extensions: [
        history(),
        keymap.of([...defaultKeymap, ...historyKeymap]),
        saveKeymap,
        markdown({ base: markdownLanguage, codeLanguages: languages }),
        syntaxHighlighting(markdownHighlightStyle),
        livePreviewPlugin,
        livePreviewTheme,
        darkTheme,
        EditorView.lineWrapping,
        // Emit content upward on every change (debounced by Svelte reactivity)
        EditorView.updateListener.of(update => {
          if (update.docChanged) {
            content = update.state.doc.toString()
          }
        }),
      ],
    })

    view = new EditorView({ state, parent: editorEl })
    view.focus()
  })

  // Sync external content changes into the editor (e.g. switching entries)
  $: if (view && content !== view.state.doc.toString()) {
    view.dispatch({
      changes: {
        from: 0,
        to:   view.state.doc.length,
        insert: content,
      },
    })
  }

  onDestroy(() => view?.destroy())
</script>

<div class="editor-wrap" bind:this={editorEl}></div>

<style>
  .editor-wrap {
    height: 100%;
    overflow: hidden;
  }

  /* Heading sizes applied via livePreview theme — nothing extra needed here */
</style>
