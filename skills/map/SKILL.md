---
name: map
description: Find all files relevant to a feature, flow, or area of the codebase and add them to your ws working set. Use when asked to "map", "track", or "find files for" a feature.
disable-model-invocation: true
allowed-tools: Read, Grep, Glob, Bash(ws *)
argument-hint: [feature or area, e.g. "user auth flow"]
---

If "$ARGUMENTS" is empty or blank: respond only with "What feature or area do you want to map? (e.g. `/ws:map user auth flow`)" and stop — do not search or add any files.

Otherwise, map files relevant to: "$ARGUMENTS"

## Steps

1. **Check existing working set**:
   ```bash
   ws list
   ```
   If `ws list` returns files, evaluate each one: does it look related to "$ARGUMENTS"?
   - If **all files are related** (or the list is empty): proceed directly to step 2.
   - If **some files are unrelated**: before searching, ask the user exactly once:
     > "Your working set has files that don't look related to **$ARGUMENTS**:
     > - `path/to/file1` — [why it seems unrelated]
     > - `path/to/file2` — [why it seems unrelated]
     >
     > Clear them and start fresh, or keep everything?"
     >
     > **A) Clear unrelated files** · **B) Keep everything**

     Wait for the user's answer:
     - **A** → run `ws clear`, then continue to step 2
     - **B** → continue to step 2 without clearing

2. **Search systematically**:
   - Glob for filenames matching keywords in the topic
   - Grep for function names, type names, route paths, and identifiers related to the topic
   - Read key files to follow imports and discover related files

3. **Add each file as you find it**:
   ```bash
   ws add <absolute-path>
   ```
   Add files one by one or in small groups. Keep searching until you've covered the relevant surface area.

4. **Scope**: Target 5–15 directly relevant files. Include:
   - Core implementation files
   - Type/interface definitions used by those files
   - Entry points and key callers

   Skip: test files (unless asked), vendor/node_modules, generated files, config files that are only tangentially related.

5. **Summarize**: After all `ws add` calls, list what you added with a one-line reason for each file. Tell the user to run `ws` in a split terminal to navigate the working set.
