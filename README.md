# Sitedog CLI

Universal CLI for sites/apps config preview & management that renders them as cards:

![sitedog-cards](https://i.imgur.com/JKjTgqY.png)

## Examples

Below are source configs for two of the above cards.

```yaml
painlessrails.com:
  dns: route 53
  hosting: hetzner
  ip: 121.143.73.32
  ssh access given to: inem, kirill

  repository: https://gitlab.com/rbbr_io/painlessrails
  registry: https://gitlab.com/rbbr_io/painlessrails/cloud/container_registry
  ci: https://gitlab.com/rbbr_io/painlessrails/ci
  issues: https://gitlab.com/rbbr_io/painlessrails/issues
  monitoring: https://appsignal.com/rbbr
  metabase: https://metabase.com
  i18n: https://poeditor.com
  analythics: https://plausible.io/
```

```yaml
inem.at:
  registrar: gandi
  dns: Route 53
  hosting: https://carrd.com
  mail: zoho
```

---

## 🛠️ How-to

1. Install the CLI (see below).
2. Run `sitedog init` to create a basic sitedog.yml.
3. Use `sitedog live` to preview your card locally while editing.
4. Generate an HTML card with `sitedog render`.

---

## 🚀 Installation

```sh
curl -sL https://get.sitedog.io | sh
```


## 📦 Usage

### 1. Create sitedog.yml

```sh
sitedog init [--config my-config.yml] # ./sitedog.yml by default
```
Alternativaly you can just create sitedog.yml manually.

### 2. Live card preview on localhost

Useful when editing config.
```sh
sitedog live [--port 3030] # http://localhost:8081 by default
```

<!-- ### 3. Add your config to cloud dashboard
```sh
sitedog push [--title another_config_title]
``` -->

### 3. Render html file with card locally
```sh
sitedog render [--output another_file.html] # ./sitedog.html by default
```

### 4 CLI help

```sh
sitedog help
Usage: sitedog <command>

Commands:
  init    Create sitedog.yml configuration file
  live    Start live server with preview
  push    Push configuration to cloud
  render  Render template to HTML
  version Print version
  help    Show this help message

Options for init:
  --config PATH    Path to config file (default: ./sitedog.yml)

Options for live:
  --config PATH    Path to config file (default: ./sitedog.yml)
  --port PORT      Port to run server on (default: 8081)
```

<!-- Options for push:
  --config PATH    Path to config file (default: ./sitedog.yml) -->

## 🚀 Uninstallation

```sh
curl -sL https://get.sitedog.io/uninstall | sh
```

If you want to upgrade to the latest version, just uninstall and install again.

---

## 📖 Read more and try it live

- [Sitedog Website](https://sitedog.io/)
- [Sitedog Live Editor](https://sitedog.io/editor#demo)

---

## 🐶 We eat our own dog food!

We use sitedog for our own projects and documentation. Here's our card deck:

[![sitedog-cli](https://i.imgur.com/DNhwj6T.png)](https://sitedog.io/live-editor.html#pako:eJytVVtv2zYUfvevIPLQN4m-zbENBGmG2m2GJVtvQfskUOKxxJYiFZLyJdiP3yElW8qQFfUwwICBc_3Odz4eWeGA6zwWejkgBJg9VEyBXJLCucouKTVpamJbsQxoZfQ3yJyltkmirKqoZIoLlVOuSyaUxSKKlZAVwKquCKvi3W4XnzxxpstjRvufaeWMlqH7sQGiar1YtpKstiKV0JVVcxWdzPElzyWD-DRDXGjreqV8DbQUOu8qQB23tgCpnZDO5tPLEUVwIi-cnyljSh26tGPRYPYoGS8DSAOV7sJy4SRLQ2VPY4KBR-pa2gaY00R5-n1-Lqwzh3NqUE8dcgQmOaaHWk5_B7UkPjrXedJGjzsfgX0lDNgluRgPx7NoiL_phYfkcwiruXC2wVWxHBJbAfDEj2hcB9C7gifeQRpz2FKmmDxYYWkIiVq4ESLXU1eJaivd03Uhr0BFb399tdGmTDYsc9pclToVEkLHXOtcArHATFYQHNHq_uobe9xEBXYaS9RGXuNcujYZJIJf2SxqZLTsqcHLfe8MI9bVm82y1zTxkLpWXGe238h7UbV09PTLw816u_4z2bwb8d8uv07L8t7dJTd2z9a72-_vP6zGw6f848Nwe0cBuQwtcBOiRMq6-v5lbEResuZVAIpOUWFuYfLHt0e5Wte32Wf1sJ6n727oZ-VQD8CvleYQ4WijaDZ85a4mXz-t8_tHab-8nd3fr1Q0DM1Yqmv3TEtF3WjpIxLxpllKCKKp1Cn1HNEPq5s3d6u45IMBPvDYPrsQJmwhk7rm58i9SzhD3yHpZXVbMFswROCNGU3m8WixiMeTSbyYo88_exQ6tgD3hHloKlGSOWo3xcan-zDwz4DrZvMRwVk980yGeCFb4Ru9FRzMkuTe2LAqBbOJw8kvXofzKPRFcHCFb6nc45SZNryJ9YeB_OOmBmNfUiw7XcMG6mCQg3uReoxxTEocvp3znDVg0XOXgCkvr-A_fCt8-9PB-8F96gf8-5HySPojWBf3BP4IsioOqd7TDYwv53wyGS3mQzbKYLjIJnwxnaXpZpjOJ-lgcLxQmRQ9jesSd8GJxMERipYvEX3s1rtxvWL_L0pyhNTe5FYJiLQ2kkT2967LM-2Qv4gtQkatfjaHnkKb7L8BYybMKw)

Even for such a small project there are already a lot of information that easily can be lost.

---

## Cloud access (Beta)

Sitedog Cloud is a platform to manage your sites/apps configs.

[Request access](https://docs.google.com/forms/d/1z5VAFvFP_fH1dJ7Y4mmNtM_AsxaFwIkQRE20zgSV0vM/edit)

---

# Development process

```sh
# Show all created versions
make show-versions

# bump cli version and create a git tag
make bump-version v=v0.1.1

# if everything is okay, push tag to main triggering release pipeline
make push-version
```

