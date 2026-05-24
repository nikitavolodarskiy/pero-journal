# pero

A minimal journalling app. Open it, write, close it. That's all.

Entries are plain Markdown files stored locally — one file per day, no cloud, no account.

---

## Download

Get the latest build from the [Releases page](https://github.com/nikitavolodarskiy/pero-journal/releases).

| Platform | File |
|---|---|
| macOS (Apple Silicon) | `pero-journal-macos-arm64.zip` |
| Windows | `pero-journal-windows-amd64.zip` |
| Linux | `pero-journal-linux-amd64.tar.gz` |

> **Intel Mac** (pre-2020) is not supported. If you need it, build from source.

### macOS

Unzip and drag `pero.app` to your Applications folder.

**First launch:** macOS will say the app "cannot be opened because the developer cannot be verified."
Right-click the app → **Open** → **Open**. You only need to do this once.

### Windows

Unzip and run `pero.exe`.

**First launch:** Windows SmartScreen may warn about an unknown publisher.
Click **More info** → **Run anyway**.

### Linux

```bash
tar -xzf pero-journal-linux-amd64.tar.gz
./pero
```

Requires WebKit2GTK at runtime:

```bash
# Ubuntu / Debian
sudo apt install libwebkit2gtk-4.0-dev

# Fedora
sudo dnf install webkit2gtk4.0
```

---

## Journal structure

Entries are stored as plain Markdown files:

```
~/pero/
  2026/
    05/
      2026-05-24.md
    04/
      2026-04-15.md
  2025/
    ...
```

You can open, edit, and search them with any text editor or grep.

---

## Configuration

Config file: `~/.config/pero/config.toml`

```toml
journal_dir = "~/pero"   # where entries are stored
editor      = "nvim"     # editor used by the CLI
```

Run `pero init` to create the config file with defaults.

Set `PERO_JOURNAL_DIR` to override the journal directory without touching the config — useful for running a dev/test instance alongside the production app:

```bash
PERO_JOURNAL_DIR=~/pero-test wails dev
```

---

## CLI

```bash
pero            # open today's entry in $EDITOR
pero list       # list all entries
pero open DATE  # open a specific date (yyyy-mm-dd)
pero stats      # word count, streaks, averages
pero init       # create default config
```

---

## Privacy & security

Pero stores your journal as **plain, unencrypted text files** on your local disk. There is no password protection, no encryption at rest, and no access control beyond what your operating system provides.

This means:

- Anyone with access to your filesystem can read your entries directly.
- If your laptop is lost, stolen, or accessed without your permission, your journal is readable by whoever gets in.
- Cloud backup and sync tools (iCloud, Dropbox, Google Drive, etc.) will upload your entries in plain text — their own privacy policies and security practices apply.
- Pero makes no guarantees about the confidentiality of your data.

If the contents of your journal are sensitive, consider enabling full-disk encryption on your machine (FileVault on macOS, BitLocker on Windows) and reviewing what you sync to the cloud. Pero is a writing tool — securing the device and the files is your responsibility.

---

## Build from source

Requirements: [Go 1.22+](https://go.dev/dl/), [Node.js 20+](https://nodejs.org), [Wails v2](https://wails.io/docs/gettingstarted/installation)

```bash
git clone https://github.com/nikitavolodarskiy/pero-journal
cd pero-journal/gui
wails build          # → gui/build/bin/pero.app (or .exe / binary)
wails dev            # hot-reload development mode
```

---

## Releasing

Tag a commit and push — GitHub Actions builds all four platforms and publishes the release automatically:

```bash
git tag v1.0.0
git push --tags
```
