# Security Policy

## Reporting a vulnerability

If you discover a security vulnerability in this SDK, please report it
privately using [GitHub's private vulnerability reporting](https://github.com/jerryisuwamakeri/Paystack-SDK/security/advisories/new)
rather than opening a public issue.

Please include:

- A description of the vulnerability and its potential impact.
- Steps to reproduce, or a minimal proof of concept.
- The SDK version affected.

## Scope

This policy covers the SDK code in this repository: the HTTP client,
retry and backoff logic, webhook signature verification, and request/
response handling.

It does not cover the Paystack API or dashboard themselves; for those,
contact Paystack support directly.

## Supported versions

Security fixes are made against the latest minor release on the current
major version line. See the compatibility matrix in [README.md](README.md).
