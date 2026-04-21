terraform {
  required_providers {
    flagd = {
      source  = "gappylul/flagd"
      version = "~> 0.1"
    }
  }
}

provider "flagd" {
  url = "https://flagd.gappy.hu"
}

resource "flagd_flag" "dark_mode" {
  name        = "dark-mode"
  enabled     = true
  description = "Dark mode UI"
}

resource "flagd_flag" "new_checkout" {
  name        = "new-checkout"
  enabled     = true
  description = "New checkout flow"
}