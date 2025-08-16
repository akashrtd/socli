# Contributing to SOCLI

Thank you for your interest in contributing to SOCLI! This document provides guidelines and information to help you get started with contributing to the project.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Fork and Clone](#fork-and-clone)
  - [Set up Development Environment](#set-up-development-environment)
- [Making Changes](#making-changes)
  - [Branching](#branching)
  - [Coding Standards](#coding-standards)
  - [Testing](#testing)
- [Submitting Changes](#submitting-changes)
  - [Pull Requests](#pull-requests)
- [Community and Communication](#community-and-communication)
- [Code of Conduct](#code-of-conduct)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.21 or higher)
- A terminal or command prompt
- A text editor or IDE (e.g., VS Code with Go extension)

### Fork and Clone

1. Fork the SOCLI repository on GitHub.
2. Clone your forked repository to your local machine:
   ```bash
   git clone https://github.com/YOUR_USERNAME/socli.git
   cd socli
   ```

### Set up Development Environment

1. Navigate to the project directory.
2. Ensure you have Go modules enabled (Go 1.11+ has this by default).
3. Run the application to verify your setup:
   ```bash
   go run .
   ```
   You might encounter a TTY error if not running from a proper terminal, but successful compilation indicates a correct setup.

## Making Changes

### Branching

1. Create a new branch for your feature or bugfix from `main`:
   ```bash
   git checkout -b feature/your-new-feature
   # or
   git checkout -b bugfix/issue-description
   ```
2. Make your changes in this branch.

### Coding Standards

- **Language:** SOCLI is written in [Go](https://golang.org/).
- **Style:** Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html) guides.
- **Formatting:** Use `gofmt` to format your code. Most editors can be configured to do this automatically on save.
- **Naming:** Use clear, descriptive names for variables, functions, and types.
- **Comments:** Comment your code, especially for complex logic or non-obvious decisions. Exported functions should have Godoc-style comments.
- **Errors:** Handle errors explicitly. Do not ignore them with `_`.
- **Logging:** Use the standard `log` package for logging in non-TUI code. For TUI-related logging or status updates, prefer sending `tui/types.StatusMsg` to the model.

### Testing

Writing tests is highly encouraged.

- SOCLI uses Go's built-in testing framework (`testing`).
- Test files are named `*_test.go` and reside in the same package as the code they test.
- Run all tests using:
  ```bash
  go test ./...
  ```
- Run tests for a specific package with verbose output:
  ```bash
  cd package_name && go test -v
  ```
- Strive for good test coverage, especially for core logic in packages like `crypto`, `internal`, `storage`, and `messaging`.

## Submitting Changes

### Pull Requests

1. Ensure your code adheres to the coding standards and passes all tests.
2. Commit your changes with a clear and concise commit message following the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification (e.g., `feat: Add new TUI feature`, `fix: Resolve crash on startup`).
3. Push your branch to your fork on GitHub:
   ```bash
   git push origin feature/your-new-feature
   ```
4. Open a Pull Request (PR) against the `main` branch of the main SOCLI repository.
5. Provide a detailed description of your changes, including the problem being solved and the solution implemented.
6. Link any relevant issues (e.g., `Closes #123`).
7. Be prepared to address feedback during the review process.

## Community and Communication

- For general discussion, feature requests, or questions, please open an issue on GitHub.
- Be respectful and constructive in all interactions.

## Code of Conduct

This project adheres to the Contributor Covenant [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.