# Setting Up Your GitHub Repository

This guide will help you set up your GitHub repository for the AlternateDNS enhanced fork.

## Step 1: Create Repository on GitHub

1. Go to https://github.com/new
2. Repository name: `AlternateDNS` (or `AlternateDNS-Enhanced` if you prefer)
3. Description: "Enhanced fork of AlternateDNS with GUI, smart DNS switching, and DNS testing features"
4. Choose **Public** (recommended for open source) or **Private**
5. **DO NOT** initialize with README, .gitignore, or license (we already have these)
6. Click **Create repository**

## Step 2: Update Git Remote

After creating the repository, GitHub will show you the repository URL. It will be:
`https://github.com/maskarajr/AlternateDNS.git`

Run these commands in your terminal:

```bash
# Remove the old remote (pointing to MaxIsJoe's repo)
git remote remove origin

# Add your repository as origin
git remote add origin https://github.com/maskarajr/AlternateDNS.git

# Verify it's set correctly
git remote -v
```

You should see:
```
origin  https://github.com/maskarajr/AlternateDNS.git (fetch)
origin  https://github.com/maskarajr/AlternateDNS.git (push)
```

## Step 3: Add and Commit All Changes

```bash
# Add all new and modified files
git add .

# Commit with a descriptive message
git commit -m "feat: Add GUI, smart DNS switching, and DNS testing features

- Added modern desktop GUI using Fyne framework
- Implemented smart DNS switching with latency comparison
- Added DNS Tester tab for benchmarking servers
- Added portable build scripts and embedded resources
- Enhanced logging and error handling
- Added comprehensive documentation and attribution

Original work by MaxIsJoe: https://github.com/MaxIsJoe/AlternateDNS"
```

## Step 4: Push to Your Repository

```bash
# Push to your repository
git push -u origin master
```

If you get an error about the branch name, GitHub might use `main` instead of `master`:

```bash
# Check your default branch name
git branch

# If GitHub uses 'main', rename your branch or push to main:
git push -u origin master:main
# OR rename your branch first:
# git branch -M main
# git push -u origin main
```

## Step 5: Verify on GitHub

1. Go to https://github.com/maskarajr/AlternateDNS
2. Verify all files are there
3. Check that README.md displays correctly with credits
4. Verify CONTRIBUTING.md, AUTHORS.md, and CHANGELOG.md are visible

## Step 6: Set Up Repository Settings (Optional but Recommended)

1. Go to **Settings** → **General**
2. Scroll to **Features**:
   - ✅ Enable **Issues** (for bug reports and feature requests)
   - ✅ Enable **Discussions** (optional, for community discussions)
   - ✅ Enable **Projects** (optional, for project management)
   - ✅ Enable **Wiki** (optional)

3. Go to **Settings** → **Branches**
   - Add branch protection rule for `master`/`main` (optional but recommended)
   - Require pull request reviews before merging

4. Go to **Settings** → **Actions** → **General**
   - Enable GitHub Actions (if you want CI/CD later)

## Step 7: Add Repository Topics (Optional)

Go to your repository → Click the gear icon (⚙️) next to "About"
Add topics like:
- `dns`
- `golang`
- `gui`
- `network-tools`
- `fyne`
- `dns-switcher`

## Step 8: Create Initial Release (Optional)

1. Go to **Releases** → **Create a new release**
2. Tag: `v1.0.0-enhanced`
3. Title: `Enhanced Fork v1.0.0 - GUI and Smart DNS Features`
4. Description:
   ```
   ## Enhanced Features
   - Modern GUI with Fyne framework
   - Smart DNS switching with latency comparison
   - Built-in DNS tester
   - Portable builds
   
   ## Credits
   Original work by [MaxIsJoe](https://github.com/MaxIsJoe/AlternateDNS)
   ```
5. Upload `AlternateDNS.exe` as a release asset (if you have a build)

## Future Contributions Workflow

### For You (Maintainer):

```bash
# Make changes
git checkout -b feature/new-feature
# ... make changes ...
git add .
git commit -m "feat: Description"
git push origin feature/new-feature
# Create PR on GitHub, review, merge
```

### For Contributors:

They will:
1. Fork your repository
2. Clone their fork
3. Make changes
4. Push to their fork
5. Create Pull Request to your repository

You'll review and merge PRs into your `master`/`main` branch.

## Keeping Track of Original Repository (Optional)

If you want to keep an eye on the original repository:

```bash
# Add original as upstream (for reference only)
git remote add upstream https://github.com/MaxIsJoe/AlternateDNS.git

# View all remotes
git remote -v
```

This allows you to:
- See what's happening in the original repo
- Compare your fork with the original
- But you won't push to it

---

**That's it!** Your repository is now set up and ready for contributions.

