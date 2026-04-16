# vaultpull

> CLI tool to sync HashiCorp Vault secrets into local `.env` files safely

---

## Installation

```bash
go install github.com/yourname/vaultpull@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/vaultpull.git
cd vaultpull
go build -o vaultpull .
```

---

## Usage

Authenticate with Vault and pull secrets into a `.env` file:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.yourtoken"

vaultpull --path secret/data/myapp --output .env
```

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path | _(required)_ |
| `--output` | Output file path | `.env` |
| `--overwrite` | Overwrite existing file | `false` |

**Example `.env` output:**

```env
DB_HOST=localhost
DB_PASSWORD=supersecret
API_KEY=abc123
```

> ⚠️ Never commit your `.env` file to version control. Add it to `.gitignore`.

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance
- A valid Vault token or auth method configured

---

## License

[MIT](LICENSE) © 2024 yourname