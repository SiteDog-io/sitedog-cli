# SiteDog CLI

Universal CLI for static site config preview & management.

```sh
# Install latest version of sitedog
curl -sL https://get.sitedog.io | sh

# Init config
sitedog init [--config my-config.yml]

# Live card preview on localhost
sitedog live [--port 3030] # http://localhost:8081 by default

# Add your config to cloud dashboard
sitedog push [--name another_config_name]

# Render html file with card locally
sitedog render [--output another_file.html] # ./sitedog.html by default

# Uninstall sidedog cli
curl -sL https://get.sitedog.io/uninstall | sh
```

Docs: [sitedog.io](https://sitedog.io/) · [gist](https://gist.github.com/qelphybox/fe278d331980a1ce09c3d946bbf0b83b)

# Release

```sh
# Show all created versions
make show-versions

# bump cli version and create a git tag
make bump-version v=v0.1.1

# if everything is okay, push tag to main triggering release pipeline
make push-version
```

