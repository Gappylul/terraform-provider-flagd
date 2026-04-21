# terraform-provider-flagd

Terraform provider for [flagd](https://github.com/gappylul/flagd) - manage feature flags as code.

## Why

Without this provider, flags live outside your infrastructure - you create them manually via curl or a dashboard. With it, flags become declarative like everything else. Check them into git, review them in PRs, roll them back with `terraform destroy`.

## Requirements

- [flagd](https://github.com/gappylul/flagd) running and reachable
- Terraform 1.0+

## Installation

```hcl
terraform {
  required_providers {
    flagd = {
      source  = "gappylul/flagd"
      version = "~> 0.1"
    }
  }
}
```

```bash
terraform init
```

## Quick start

```hcl
provider "flagd" {
  url       = "https://flagd.yourdomain.com"
  admin_key = var.flagd_admin_key   # or set FLAGD_ADMIN_KEY env var
}

resource "flagd_flag" "dark_mode" {
  name        = "dark-mode"
  enabled     = false
  description = "Dark mode UI"
}

resource "flagd_flag" "new_checkout" {
  name        = "new-checkout"
  enabled     = true
  description = "New checkout flow"
}
```

```bash
terraform apply -var="flagd_admin_key=<your-key>"
```

## Provider configuration

| Argument    | Env var           | Description                                         |
|-------------|-------------------|-----------------------------------------------------|
| `url`       | `FLAGD_URL`       | flagd server URL. Defaults to http://localhost:8080 |
| `admin_key` | `FLAGD_ADMIN_KEY` | Bearer token for write access                       |

## Resources

### `flagd_flag`

Manages a feature flag. Renaming a flag (`name`) forces destroy and recreate - the old flag is deleted and a new one is created.

```hcl
resource "flagd_flag" "example" {
  name        = "my-feature"  # required, immutable
  enabled     = false         # optional, default false
  description = "My feature"  # optional
}
```

**Attributes**

| Name          | Type   | Description                                       |
|---------------|--------|---------------------------------------------------|
| `name`        | string | Unique flag name. Changing it forces replacement. |
| `enabled`     | bool   | Whether the flag is enabled. Default `false`.     |
| `description` | string | Human-readable description.                       |
| `created_at`  | string | Set by server. Read-only.                         |
| `updated_at`  | string | Set by server. Read-only.                         |

## Data sources

### `flagd_flag`

Read an existing flag without managing it. Useful for referencing flags created outside Terraform.

```hcl
data "flagd_flag" "existing" {
  name = "some-flag"
}

output "is_enabled" {
  value = data.flagd_flag.existing.enabled
}
```

## Drift detection

If someone toggles a flag manually via the flagd dashboard or API, `terraform plan` will detect the drift:

```bash
terraform plan -var="flagd_admin_key=<your-key>"

# ~ update in-place
# ~ enabled: true -> false
```

`terraform apply` brings it back in line with your declared state.

## License

MIT