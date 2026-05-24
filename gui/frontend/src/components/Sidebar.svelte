<script>
  export let entries    = []
  export let stats      = null
  export let activeDate = ''
  export let onSelect   = (_date) => {}
  export let onRefresh  = () => {}
  export let onToday    = () => {}

  let refreshing = false

  async function handleRefresh() {
    if (refreshing) return
    refreshing = true
    // Run the actual refresh and a minimum spin duration in parallel so the
    // animation is always visible even when the reload completes instantly.
    await Promise.all([onRefresh(), new Promise(r => setTimeout(r, 600))])
    refreshing = false
  }

  const now          = new Date()
  const currentYear  = now.getFullYear()
  const currentMonth = now.getMonth() + 1   // 1-indexed

  // ── Grouping ──────────────────────────────────────────────────────────────

  $: grouped = buildTree(entries)

  function buildTree(entries) {
    const map = {}
    for (const entry of entries) {
      const date       = entry.Date.slice(0, 10)
      const [y, m]     = date.split('-').map(Number)
      if (!map[y])     map[y] = {}
      if (!map[y][m])  map[y][m] = []
      map[y][m].push(date)
    }

    // Years descending → months descending → dates descending (already sorted
    // newest-first from the backend, but we enforce it here too).
    return Object.keys(map)
      .map(Number)
      .sort((a, b) => b - a)
      .map(year => ({
        year,
        months: Object.keys(map[year])
          .map(Number)
          .sort((a, b) => b - a)
          .map(month => ({
            month,
            dates: [...map[year][month]].sort().reverse(),
          }))
      }))
  }

  // ── Toggle state ──────────────────────────────────────────────────────────
  // Current year and current month start open; everything else collapsed.

  let openYears  = new Set([currentYear])
  let openMonths = new Set([`${currentYear}-${currentMonth}`])

  function toggleYear(year) {
    openYears = toggle(openYears, year)
  }

  function toggleMonth(year, month) {
    openMonths = toggle(openMonths, `${year}-${month}`)
  }

  // Returns a new Set with key added or removed.
  function toggle(set, key) {
    const next = new Set(set)
    next.has(key) ? next.delete(key) : next.add(key)
    return next
  }

  // ── Formatting ────────────────────────────────────────────────────────────

  function monthName(m) {
    return new Date(2000, m - 1, 1).toLocaleString('en-US', { month: 'long' })
  }

  // Within a month group the year+month are already shown, so just show
  // the weekday + day number: "Sun 24".
  function formatDay(dateStr) {
    const d = new Date(dateStr + 'T00:00:00')
    return d.toLocaleString('en-US', { weekday: 'short', day: 'numeric' })
  }

  function isToday(dateStr) {
    return dateStr === new Date().toISOString().slice(0, 10)
  }
</script>

<aside class="sidebar">

  <section class="entries">
    <div class="section-header">
      <h2 class="section-title">Entries</h2>
      <div class="header-actions">
        <button class="icon-btn" on:click={onToday}   title="Open or create today's entry">+</button>
        <button class="icon-btn" class:spinning={refreshing} on:click={handleRefresh} title="Refresh entries">↻</button>
      </div>
    </div>

    {#if grouped.length === 0}
      <p class="empty">No entries yet</p>
    {:else}
      {#each grouped as { year, months }}

        <!-- Year row -->
        <button
          class="group-toggle year-toggle"
          on:click={() => toggleYear(year)}
          aria-expanded={openYears.has(year)}
        >
          <span class="chevron" class:open={openYears.has(year)}>›</span>
          <span>{year}</span>
        </button>

        {#if openYears.has(year)}
          {#each months as { month, dates }}

            <!-- Month row -->
            <button
              class="group-toggle month-toggle"
              on:click={() => toggleMonth(year, month)}
              aria-expanded={openMonths.has(`${year}-${month}`)}
            >
              <span class="chevron" class:open={openMonths.has(`${year}-${month}`)}>›</span>
              <span>{monthName(month)}</span>
              <span class="month-count">{dates.length}</span>
            </button>

            {#if openMonths.has(`${year}-${month}`)}
              <ul class="entry-list">
                {#each dates as date}
                  <li>
                    <button
                      class="entry-item"
                      class:active={date === activeDate}
                      on:click={() => onSelect(date)}
                    >
                      <span class="entry-date">{formatDay(date)}</span>
                      {#if isToday(date)}
                        <span class="today-dot" title="Today">●</span>
                      {/if}
                    </button>
                  </li>
                {/each}
              </ul>
            {/if}

          {/each}
        {/if}

      {/each}
    {/if}
  </section>

  {#if stats}
    <section class="stats">
      <h2 class="section-title">Stats</h2>
      <dl class="stat-list">
        <div class="stat-row">
          <dt>Streak</dt>
          <dd>{stats.CurrentStreak}d</dd>
        </div>
        <div class="stat-row">
          <dt>Longest</dt>
          <dd>{stats.LongestStreak}d</dd>
        </div>
        <div class="stat-row">
          <dt>Entries</dt>
          <dd>{stats.TotalEntries}</dd>
        </div>
        <div class="stat-row">
          <dt>Words</dt>
          <dd>{stats.TotalWords.toLocaleString()}</dd>
        </div>
      </dl>
    </section>
  {/if}

</aside>

<style>
  .sidebar {
    display:        flex;
    flex-direction: column;
    width:          200px;
    min-width:      200px;
    background:     #252526;
    border-right:   1px solid #333;
    overflow:       hidden;
    user-select:    none;
  }

  .section-header {
    display:     flex;
    align-items: center;
    padding:     16px 8px 8px 16px;
  }

  .section-title {
    flex:           1;
    font-size:      10px;
    font-weight:    600;
    letter-spacing: 0.08em;
    text-transform: uppercase;
    color:          #666;
    margin:         0;
    padding:        16px 16px 8px;
  }

  /* Inside the flex header row the padding is handled by the row itself. */
  .section-header .section-title { padding: 0; }

  .header-actions {
    display:  flex;
    gap:      2px;
  }

  .icon-btn {
    display:     flex;
    align-items: center;
    justify-content: center;
    width:       22px;
    height:      22px;
    border:      none;
    border-radius: 4px;
    background:  transparent;
    color:       #555;
    font-size:   14px;
    line-height: 1;
    cursor:      pointer;
    transition:  color 0.1s, background 0.1s;
  }

  .icon-btn:hover {
    color:      #ccc;
    background: #2d2d2d;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .icon-btn.spinning {
    animation:        spin 0.6s linear infinite;
    color:            #aaa;
    pointer-events:   none;   /* prevent double-click during spin */
  }

  .entries {
    flex:       1;
    overflow-y: auto;
    padding-bottom: 8px;
  }

  /* ── Year / Month toggle rows ── */

  .group-toggle {
    display:     flex;
    align-items: center;
    gap:         5px;
    width:       100%;
    border:      none;
    background:  transparent;
    cursor:      pointer;
    text-align:  left;
    color:       #888;
  }

  .group-toggle:hover { color: #bbb; }

  .year-toggle {
    padding:     5px 10px;
    font-size:   12px;
    font-weight: 600;
  }

  .month-toggle {
    padding:     3px 10px 3px 22px;
    font-size:   11px;
    font-weight: 500;
  }

  .chevron {
    display:          inline-block;
    font-size:        12px;
    line-height:      1;
    transform:        rotate(0deg);
    transition:       transform 0.15s ease;
    color:            #555;
    flex-shrink:      0;
  }
  .chevron.open { transform: rotate(90deg); }

  .month-count {
    margin-left: auto;
    font-size:   10px;
    color:       #444;
  }

  /* ── Entry list ── */

  .entry-list {
    list-style: none;
    margin:     0;
    padding:    0 8px 2px 28px;
  }

  .entry-item {
    display:         flex;
    align-items:     center;
    justify-content: space-between;
    width:           100%;
    padding:         4px 8px;
    border:          none;
    border-radius:   4px;
    background:      transparent;
    color:           #999;
    font-size:       12px;
    cursor:          pointer;
    text-align:      left;
    transition:      background 0.1s;
  }

  .entry-item:hover  { background: #2d2d2d; color: #ddd; }
  .entry-item.active { background: #37373d; color: #fff; }

  .today-dot {
    font-size: 7px;
    color:     #7f9cf5;
  }

  .empty {
    padding:   8px 16px;
    font-size: 12px;
    color:     #555;
  }

  /* ── Stats ── */

  .stats {
    border-top:     1px solid #333;
    padding-bottom: 16px;
  }

  .stat-list {
    margin:  0;
    padding: 0 16px;
  }

  .stat-row {
    display:         flex;
    justify-content: space-between;
    padding:         3px 0;
  }

  dt { font-size: 12px; color: #666; }

  dd {
    font-size:   12px;
    font-weight: 600;
    color:       #bbb;
    margin:      0;
  }
</style>
