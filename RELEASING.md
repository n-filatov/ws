# Releasing

## One-time setup

### 1. Create the Homebrew tap repo

Go to GitHub and create a new **public** repo named `homebrew-tap` under `n-filatov`.
It can be completely empty — GoReleaser will push the formula on first release.

### 2. Create a Personal Access Token for the tap

Go to **GitHub → Settings → Developer settings → Fine-grained tokens → Generate new token**.

- **Resource owner**: `n-filatov`
- **Repository access**: Only `homebrew-tap`
- **Permissions**: Contents → **Read and write**

Copy the token and add it as a secret in the `ws` repo:

```
Settings → Secrets and variables → Actions → New repository secret
Name:  HOMEBREW_TAP_TOKEN
Value: <your token>
```

### 3. Generate a GPG key for APT signing

```bash
gpg --batch --gen-key <<EOF
Key-Type: RSA
Key-Length: 4096
Name-Real: ws releases
Name-Email: releases@n-filatov.github.io
Expire-Date: 0
%no-protection
EOF
```

Export the private key and add it as a secret:

```bash
gpg --armor --export-secret-keys releases@n-filatov.github.io
```

```
Settings → Secrets and variables → Actions → New repository secret
Name:  GPG_PRIVATE_KEY
Value: <paste the armored private key, including -----BEGIN/END lines>
```

Add a second secret for the passphrase (leave empty if you used `%no-protection`):

```
Name:  GPG_PASSPHRASE
Value: (empty or your passphrase)
```

### 4. Initialize the `gh-pages` branch

```bash
git checkout --orphan gh-pages
git rm -rf .
echo "# ws APT Repository" > README.md
git add README.md
git commit -m "init gh-pages"
git push origin gh-pages
git checkout main
```

### 5. Enable GitHub Pages

In the `ws` repo: **Settings → Pages → Source: Deploy from branch → Branch: `gh-pages` → / (root)**.

The APT repo will be served at `https://n-filatov.github.io/ws/`.

---

## Cutting a release

```bash
make release VERSION=v0.1.0
```

This tags the commit and pushes the tag, which triggers two workflows:

1. **`release.yml`** — GoReleaser builds binaries for macOS + Linux (amd64/arm64),
   creates a GitHub Release with tarballs + `.deb` files, and updates the
   Homebrew formula in `n-filatov/homebrew-tap`.

2. **`apt-repo.yml`** — Downloads the `.deb` files from the release, regenerates
   the APT index, GPG-signs it, and pushes to `gh-pages`.

---

## User installation

### Homebrew (macOS / Linux)

```bash
brew tap n-filatov/tap
brew install ws
```

### APT (Debian / Ubuntu)

```bash
curl -fsSL https://n-filatov.github.io/ws/gpg.key \
  | sudo gpg --dearmor -o /usr/share/keyrings/ws.gpg

echo "deb [signed-by=/usr/share/keyrings/ws.gpg] https://n-filatov.github.io/ws ./" \
  | sudo tee /etc/apt/sources.list.d/ws.list

sudo apt update && sudo apt install ws
```

### Manual / direct download

Grab the tarball or `.deb` from the [Releases page](https://github.com/n-filatov/ws/releases).
