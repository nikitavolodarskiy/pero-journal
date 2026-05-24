/**
 * Custom highlight style for pero's markdown editor.
 *
 * Replaces defaultHighlightStyle entirely. Key differences:
 *  - Headings: sized properly (1.8em → 1.1em), bold, NO underline
 *  - Markdown markers (##, **, *) shown in a subtle grey when visible
 *  - Bold, italic, code, links all styled correctly
 */

import { HighlightStyle } from '@codemirror/language'
import { tags } from '@lezer/highlight'

export const markdownHighlightStyle = HighlightStyle.define([
  // ── Headings ── bigger, bold, never underlined ──────────────────────────
  { tag: tags.heading1, fontSize: '1.8em', fontWeight: '700', lineHeight: '1.3', color: '#e8e8e8' },
  { tag: tags.heading2, fontSize: '1.5em', fontWeight: '700', lineHeight: '1.3', color: '#e8e8e8' },
  { tag: tags.heading3, fontSize: '1.3em', fontWeight: '600', lineHeight: '1.3', color: '#e0e0e0' },
  { tag: tags.heading4, fontSize: '1.1em', fontWeight: '600', color: '#e0e0e0' },
  { tag: tags.heading5, fontSize: '1em',   fontWeight: '600', color: '#d8d8d8' },
  { tag: tags.heading6, fontSize: '1em',   fontWeight: '600', color: '#d0d0d0' },

  // ── Inline formatting ────────────────────────────────────────────────────
  { tag: tags.strong,   fontWeight: '700' },
  { tag: tags.emphasis, fontStyle: 'italic' },

  // ── Inline code ──────────────────────────────────────────────────────────
  { tag: tags.monospace, fontFamily: 'ui-monospace, "Fira Code", monospace', fontSize: '0.88em' },

  // ── Links ────────────────────────────────────────────────────────────────
  { tag: tags.link, color: '#7f9cf5', textDecoration: 'underline' },
  { tag: tags.url,  color: '#4a9eff' },

  // ── Markdown punctuation (##, **, *, `, [, ], etc.) ─────────────────────
  // These are visible when the cursor is on the line (live preview off).
  // Keep them subtle so they don't fight with the content.
  { tag: tags.processingInstruction, color: '#6b7280' },
  { tag: tags.punctuation,           color: '#6b7280' },

  // ── Blockquotes ──────────────────────────────────────────────────────────
  { tag: tags.quote,   fontStyle: 'italic', color: '#9ca3af' },

  // ── Code blocks ──────────────────────────────────────────────────────────
  { tag: tags.meta, color: '#6b7280' },
])
