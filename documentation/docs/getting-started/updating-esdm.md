# Updating ESDM

This guide shows how to **update ESDM** when a new release is available, and explains how the CLI lets you know that an update exists in the first place. By the end you'll know which steps to take to upgrade, what `esdm` does in the background, and how to turn it off.

ESDM doesn't ship with an automatic self-updater. The CLI only **tells** you that a newer version is out – swapping the binary stays in your hands.

## Performing the Update

Updating is the **same procedure** as a fresh install. Download the latest binary, replace the old one, and run `esdm version` to confirm the new version. The detailed per-platform steps are on the **[Installing ESDM](/getting-started/installing-esdm.md)** page – binary names, permission flags, and verification commands are identical between install and update.

If your binary lives in a directory that requires elevated permissions (such as `/usr/local/bin` on macOS and Linux), use `sudo` when overwriting it.

## How the Update Notification Works

!!! info "Available from ESDM 0.10.0"

    The notification mechanism ships in **esdm 0.10.0 and later**. If you're running 0.9.0 or earlier, the CLI doesn't yet check for new releases on its own – update manually once to 0.10.0 (or a later version) using the steps above, and the notification starts working from then on.

When you run any `esdm` command in an interactive terminal, the CLI checks **once a day** whether a newer release has been published. The check is a single HTTP `GET` request to **[www.esdm.io/version.json](https://www.esdm.io/version.json)** that returns nothing more than the latest known version number. If your installed binary is older, you'll see a hint right after the command's output:

```text
⚡ A new version of esdm is available: <current> → <latest>
  See https://www.esdm.io/getting-started/updating-esdm/ for upgrade instructions.
```

`<current>` is the version your binary reports; `<latest>` is the version returned by the endpoint. The first time the hint shows up on a machine, an extra line tells you how to switch the check off.

### What Gets Sent

The check fetches a **small static file** and nothing else. There's no analytics call, no machine identifier, no schema or model contents – the request looks like any other HTTP `GET` against a public URL. The response payload is a JSON document of the form `{"version": "0.14.0"}`.

### When the Check Stays Quiet

The notification is **deliberately suppressed** in situations where it would be unwelcome:

- **Non-interactive runs** – if stderr isn't attached to a terminal (pipes, redirects, scripts), nothing is printed.
- **CI environments** – if the `CI` environment variable is set, the check doesn't run. Most CI providers set this automatically.
- **Development builds** – if you're running an unreleased build from source, no comparison is possible and the check skips.

### Disabling the Check

To switch the check off entirely, set the environment variable below:

```shell
export ESDM_DISABLE_UPDATE_CHECK=true
```

Add the line to your shell's profile (`.zshrc`, `.bashrc`, …) to make the setting permanent. With the variable in place, the CLI doesn't make any network request as part of update detection.

## Where to Go Next

- **[Installing ESDM](/getting-started/installing-esdm.md)** carries the per-platform download tables and the renaming and permission steps that the update process reuses.
- **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)** has a coding agent draft a small model for you, useful right after an install or update to verify everything is wired up correctly.
- **[Your First Model by Hand](/getting-started/your-first-model.md)** does the same artifact by artifact, if you'd rather write the YAML yourself.
