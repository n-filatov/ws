---
name: install
description: Install the ws CLI tool using the recommended method for the current OS.
disable-model-invocation: true
allowed-tools: Bash
argument-hint: (no arguments needed)
---

Install the `ws` CLI tool on this machine.

## Steps

1. **Check if already installed**:
   ```bash
   which ws && ws version
   ```
   If `ws` is found, report the version and stop — nothing to do.

2. **Detect OS and available package managers**:
   ```bash
   uname -s
   ```
   Then check what's available:
   ```bash
   command -v brew
   command -v apt-get
   ```

3. **Install using the right method**:

   **macOS or Linux with Homebrew** (`brew` is available):
   ```bash
   brew tap n-filatov/tap && brew install ws
   ```

   **Debian / Ubuntu** (`apt-get` is available, no `brew`):
   ```bash
   curl -fsSL https://n-filatov.github.io/ws/gpg.key \
     | sudo gpg --dearmor -o /usr/share/keyrings/ws.gpg
   echo "deb [signed-by=/usr/share/keyrings/ws.gpg] https://n-filatov.github.io/ws ./" \
     | sudo tee /etc/apt/sources.list.d/ws.list
   sudo apt update && sudo apt install ws
   ```

   **Other Linux** (no `brew`, no `apt-get`):
   ```bash
   git clone https://github.com/n-filatov/ws /tmp/ws-install
   cd /tmp/ws-install && make install
   ```
   This installs to `~/.local/bin/ws`. Remind the user to ensure `~/.local/bin` is in their `$PATH`.

4. **Verify**:
   ```bash
   ws version
   ```

5. **Report** what was installed and suggest running `/ws:map <feature>` to get started.
