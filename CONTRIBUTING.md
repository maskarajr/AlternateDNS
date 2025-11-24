# Contributing to AlternateDNS (Community Edition)

Thank you for your interest in contributing to AlternateDNS! This document provides guidelines and information for contributors.

## About This Fork

This is the **actively maintained fork** of the original AlternateDNS project. All future issues, features, and pull requests should target **this repository**. We continue to credit and respect the original author while evolving the project.

## Original Author & Attribution

**Original Author:** [MaxIsJoe](https://github.com/MaxIsJoe)  
**Original Repository:** https://github.com/MaxIsJoe/AlternateDNS

We are building on top of MaxIsJoe’s original work. Keep attribution intact in source files and changelog entries, and mention the original project when describing major feature additions.

## How to Contribute

### 1. Fork and Clone

1. Fork **this repository** (maskarajr/AlternateDNS) on GitHub.
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/AlternateDNS.git
   cd AlternateDNS
   ```

### 2. (Optional) Track the Original Project

If you need to compare with the upstream project, add a read-only remote:

```bash
git remote add upstream https://github.com/MaxIsJoe/AlternateDNS.git
```

`origin` should always point to your fork of this community edition. Never open pull requests against the original upstream repository unless coordinated with MaxIsJoe.

### 3. Create a Branch

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

Open a Pull Request against **maskarajr/AlternateDNS**. Clearly describe the motivation, implementation details, and testing performed. PRs to the original repo will be closed and redirected here.

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
- Keep original copyright notices and references to MaxIsJoe intact.
- Mention the original project in release notes or large feature PRs.
- Add yourself to AUTHORS.md for meaningful contributions.
- Follow the established style so the codebase stays approachable.

Thank you for contributing!

