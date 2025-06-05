# SiteDog CLI

Universal CLI for static site config preview & management.

```sh
# Install
curl -sL https://get.sitedog.io | sh

# Init config
touch sitedog.yml && sitedog init

# Live preview
sitedog live  # http://localhost:8081

# Build all binaries (for devs, via Docker)
make build

# Update gist (binaries/scripts)
make push

# Uninstall
~/.sitedog/bin/sitedog uninstall
```

Docs: [sitedog.io](https://sitedog.io/) · [gist](https://gist.github.com/qelphybox/fe278d331980a1ce09c3d946bbf0b83b)

---

© rbbr.io, 2024
