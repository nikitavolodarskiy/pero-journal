/**
 * Live Preview extension for CodeMirror 6.
 *
 * When the cursor is NOT on a line, markdown syntax is hidden and the line
 * is rendered as formatted text (headings, bold, italic, code).
 * When the cursor moves onto a line, decorations are removed and raw
 * markdown source is shown — exactly like Obsidian's Live Preview mode.
 */

import { Decoration, EditorView, ViewPlugin, WidgetType } from '@codemirror/view'
import { RangeSetBuilder } from '@codemirror/state'
import { syntaxTree } from '@codemirror/language'

// ─── Widget types ────────────────────────────────────────────────────────────

class HiddenWidget extends WidgetType {
  toDOM() {
    const span = document.createElement('span')
    span.style.display = 'none'
    return span
  }
  ignoreEvent() { return false }
}

const hidden = Decoration.replace({ widget: new HiddenWidget() })

// ─── Helpers ─────────────────────────────────────────────────────────────────

/** Returns the set of line numbers the cursor(s) are on. */
function cursorLines(view) {
  const lines = new Set()
  for (const range of view.state.selection.ranges) {
    const from = view.state.doc.lineAt(range.from).number
    const to   = view.state.doc.lineAt(range.to).number
    for (let n = from; n <= to; n++) lines.add(n)
  }
  return lines
}

/** True when any part of [from, to) overlaps a cursor line. */
function onCursorLine(from, to, view, curLines) {
  const startLine = view.state.doc.lineAt(from).number
  const endLine   = view.state.doc.lineAt(to).number
  for (let n = startLine; n <= endLine; n++) {
    if (curLines.has(n)) return true
  }
  return false
}

// ─── Decoration builder ──────────────────────────────────────────────────────

function buildDecorations(view) {
  const builder  = new RangeSetBuilder()
  const curLines = cursorLines(view)
  const tree     = syntaxTree(view.state)
  const doc      = view.state.doc

  // Walk every visible range in the viewport.
  for (const { from, to } of view.visibleRanges) {
    tree.iterate({
      from, to,
      enter(node) {
        const { name, from: nFrom, to: nTo } = node

        // Skip anything on the cursor line — show raw source.
        if (onCursorLine(nFrom, nTo, view, curLines)) return

        // ── Headings ── hide the `# ` prefix ─────────────────────────────
        if (name === 'HeaderMark') {
          // HeaderMark covers the `#` symbols + trailing space.
          builder.add(nFrom, nTo, hidden)
        }

        // ── Bold markers ── hide `**` or `__` ────────────────────────────
        if (name === 'StrongEmphasis') {
          // StrongEmphasis wraps the whole `**text**`. Its first and last
          // two characters are the markers — hide them.
          builder.add(nFrom, nFrom + 2, hidden)
          builder.add(nTo - 2, nTo, hidden)
        }

        // ── Italic markers ── hide `*` or `_` ────────────────────────────
        if (name === 'Emphasis') {
          // Only process if it is NOT inside StrongEmphasis.
          const parent = node.node.parent
          if (parent && parent.name !== 'StrongEmphasis') {
            builder.add(nFrom, nFrom + 1, hidden)
            builder.add(nTo - 1, nTo, hidden)
          }
        }

        // ── Inline code backticks ── hide `` ` `` ────────────────────────
        if (name === 'InlineCode') {
          builder.add(nFrom, nFrom + 1, hidden)
          builder.add(nTo - 1, nTo, hidden)
        }

        // ── Link syntax ── hide [text](url) → show only text ─────────────
        if (name === 'Link') {
          // Find [ ] and ( ) delimiters inside the Link node.
          tree.iterate({
            from: nFrom, to: nTo,
            enter(child) {
              if (child.name === 'LinkMark' || child.name === 'URL') {
                builder.add(child.from, child.to, hidden)
              }
            }
          })
        }
      }
    })
  }

  return builder.finish()
}

// ─── ViewPlugin ──────────────────────────────────────────────────────────────

const livePreviewPlugin = ViewPlugin.fromClass(
  class {
    constructor(view) {
      this.decorations = buildDecorations(view)
    }

    update(update) {
      if (
        update.docChanged      ||
        update.selectionSet    ||
        update.viewportChanged ||
        update.geometryChanged
      ) {
        this.decorations = buildDecorations(update.view)
      }
    }
  },
  { decorations: v => v.decorations }
)

// ─── Inline code background ──────────────────────────────────────────────────
// HighlightStyle handles colours and font; we add a subtle background via
// baseTheme because HighlightStyle doesn't support background on inline spans
// as reliably across all CM6 versions.
export const livePreviewTheme = EditorView.baseTheme({
  // Target any span whose inline style contains a monospace font family.
  // Using "monospace" (the generic family, always last in our stack) is more
  // robust than matching a specific font name like "ui-monospace" or "Fira Code".
  '.cm-line span[style*="monospace"]': {
    background:   'rgba(255,255,255,0.07)',
    borderRadius: '3px',
    padding:      '0 3px',
  },
})

export { livePreviewPlugin }
