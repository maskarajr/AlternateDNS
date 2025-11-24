# Contributing to AlternateDNS (Enhanced Fork)

Thank you for your interest in contributing to AlternateDNS! This document provides guidelines and information for contributors.

## About This Fork

This is a **community-maintained fork** of the original AlternateDNS project. We maintain this fork to add features, accept contributions, and provide ongoing updates while properly crediting the original author.

## Original Author & Attribution

**Original Author:** [MaxIsJoe](https://github.com/MaxIsJoe)  
**Original Repository:** https://github.com/MaxIsJoe/AlternateDNS

All contributions to this fork must maintain proper attribution to the original author. When making significant changes, please acknowledge the original work in your commit messages and pull requests. The original author's work and code structure should always be credited.

## How to Contribute

### 1. Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/AlternateDNS.git
   cd AlternateDNS
   ```

### 2. Set Up Remotes (Optional)

If you want to track the original repository for reference:

```bash
git remote add upstream https://github.com/MaxIsJoe/AlternateDNS.git
```

Verify your remotes:
```bash
git remote -v
```

You should see:
- `origin` - this fork's repository
- `upstream` - original repository (optional, for reference only)

### 3. Create a Branch

Create a feature branch for your changes:

```bash
git checkout -b feature/your-feature-name
```

### 4. Make Your Changes

- Follow the existing code style
- Add comments for complex logic
- Test your changes thoroughly
- Update documentation as needed

### 5. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git add .
git commit -m "Add: Description of your changes"
```

### 6. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub to this repository (the enhanced fork).

**Note:** Pull requests should be made to this fork's repository, not the original repository.

## Code Style

- Follow Go conventions (gofmt, golint)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and small

## Testing

Before submitting a PR, ensure:
- The code compiles without errors
- All existing functionality still works
- New features are tested
- No console errors appear

## Pull Request Guidelines

- Use clear, descriptive titles
- Describe what changes were made and why
- Reference any related issues
- Include screenshots for UI changes
- Ensure all checks pass

## Questions?

If you have questions, please open an issue on this repository.

## Attribution Requirements

When contributing, please:
- Maintain all existing attribution to the original author
- Add your name to commit messages for significant contributions
- Acknowledge the original work when making major changes
- Keep the original author's code structure and style where possible

Thank you for contributing!

