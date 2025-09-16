# Setting up homebrew-syno-vm Tap Repository

This directory contains the setup files for creating the Homebrew tap repository.

## Steps to create the tap repository:

1. **Create the GitHub repository**:
   ```bash
   # Create a new repository at: https://github.com/scttfrdmn/homebrew-syno-vm
   ```

2. **Initialize the repository structure**:
   ```bash
   git clone https://github.com/scttfrdmn/homebrew-syno-vm.git
   cd homebrew-syno-vm
   mkdir Formula
   cp ../homebrew-tap-setup/README-tap.md ./README.md
   cp ../homebrew-tap-setup/.github-workflows-update-formula.yml ./.github/workflows/update-formula.yml
   git add .
   git commit -m "Initial tap repository setup"
   git push origin main
   ```

3. **Generate a GitHub Personal Access Token**:
   - Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
   - Generate new token with scopes:
     - `public_repo` (for public repositories)
     - `repo` (if using private repositories)
   - Name it something like "GoReleaser Homebrew Tap Token"

4. **Add the token to syno-vm repository secrets**:
   - Go to https://github.com/scttfrdmn/syno-vm/settings/secrets/actions
   - Add new repository secret:
     - Name: `HOMEBREW_TAP_GITHUB_TOKEN`
     - Value: [your personal access token]

5. **Test the setup**:
   - Create a tag on syno-vm: `git tag v0.1.0 && git push origin v0.1.0`
   - This should trigger GoReleaser and create the formula in homebrew-syno-vm

## Directory Structure After Setup:
```
homebrew-syno-vm/
├── README.md
├── Formula/
│   └── syno-vm.rb (created by GoReleaser)
└── .github/
    └── workflows/
        └── update-formula.yml
```