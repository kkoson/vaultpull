# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files with profile support

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

```bash
# Authenticate and pull secrets into a .env file
vaultpull pull --profile staging --output .env

# List available profiles
vaultpull profiles list

# Use a specific Vault address and secret path
vaultpull pull --addr https://vault.example.com --path secret/myapp --output .env.local
```

### Configuration

Create a `vaultpull.yaml` in your project root to define profiles:

```yaml
profiles:
  staging:
    addr: https://vault.example.com
    path: secret/data/myapp/staging
  production:
    addr: https://vault.example.com
    path: secret/data/myapp/production
```

Then run:

```bash
vaultpull pull --profile staging
```

Secrets are written to the specified output file in standard `KEY=VALUE` format.

---

## Requirements

- Go 1.21+
- A running [HashiCorp Vault](https://www.vaultproject.io/) instance
- A valid `VAULT_TOKEN` set in your environment

---

## License

MIT © [yourusername](https://github.com/yourusername)