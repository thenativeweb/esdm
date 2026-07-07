# Installing ESDM

This guide shows how to **install ESDM** and verify that it's working correctly. It covers **macOS**, **Linux**, and **Windows**. At the end, you will have a **working `esdm` binary**, ready to lint and view your models.

ESDM is distributed as **pre-built binaries** for your operating system and CPU architecture. Download the binary, place it in your project, and you're done.

## Downloading the Latest Version

To get the latest version of ESDM, select your operating system and CPU architecture:

=== "macOS"

    | Architecture | Version | Download | Size | SHA256 |
    |--------------|---------|----------|------|--------|
    | ARM64 | 0.14.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.14.0/esdm-darwin-arm64)** | 8.5 MB | `b124e609df2c7066fb153b5ac294ef4bcf8fba2807eb3cbe2879c6debd958bf6` |
    | x86 | 0.14.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.14.0/esdm-darwin-amd64)** | 9.1 MB | `f9ca09548b783b214b5455cc7bc2bb8aba373b84ac2670b209657898cc75c548` |

=== "Linux"

    | Architecture | Version | Download | Size | SHA256 |
    |--------------|---------|----------|------|--------|
    | ARM64 | 0.14.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.14.0/esdm-linux-arm64)** | 8.3 MB | `891fe719fa2ea24eef9ebd229b5032d769791ac926e091a2cacfd91520edbfb4` |
    | x86 | 0.14.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.14.0/esdm-linux-amd64)** | 9.0 MB | `c0a786972300f6f7e71e645f009b8e7b8b7967c8837daf9e51f968e756e1716e` |

=== "Windows"

    | Architecture | Version | Download | Size | SHA256 |
    |--------------|---------|----------|------|--------|
    | ARM64 | 0.14.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.14.0/esdm-windows-arm64.exe)** | 8.5 MB | `b4c9b7744d3036977eec83f11ead009a0052028cd1eec53cb8e16db68b276909` |
    | x86 | 0.14.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.14.0/esdm-windows-amd64.exe)** | 9.3 MB | `95d6140572209fc3a52449180cd538da01475e97781580b1c1ed5933d0d2e2fb` |

<!--
Previous Versions block: add once the first 1.0.0 release ships, mirroring
EventSourcingDB's installing page.

??? note "Previous Versions"

    === "macOS"

        | Architecture | Version | Download | Size | SHA256 |
        |--------------|---------|----------|------|--------|
        | ARM64 | <version> | **[Download](https://esdm.s3.fr-par.scw.cloud/<version>/esdm-darwin-arm64)** | <size> | `<sha256>` |
        | x86 | <version> | **[Download](https://esdm.s3.fr-par.scw.cloud/<version>/esdm-darwin-amd64)** | <size> | `<sha256>` |

    === "Linux"

        | Architecture | Version | Download | Size | SHA256 |
        |--------------|---------|----------|------|--------|
        | ARM64 | <version> | **[Download](https://esdm.s3.fr-par.scw.cloud/<version>/esdm-linux-arm64)** | <size> | `<sha256>` |
        | x86 | <version> | **[Download](https://esdm.s3.fr-par.scw.cloud/<version>/esdm-linux-amd64)** | <size> | `<sha256>` |

    === "Windows"

        | Architecture | Version | Download | Size | SHA256 |
        |--------------|---------|----------|------|--------|
        | ARM64 | <version> | **[Download](https://esdm.s3.fr-par.scw.cloud/<version>/esdm-windows-arm64.exe)** | <size> | `<sha256>` |
        | x86 | <version> | **[Download](https://esdm.s3.fr-par.scw.cloud/<version>/esdm-windows-amd64.exe)** | <size> | `<sha256>` |
-->

## Post-Download Steps

### Renaming the Binary

Rename the binary for simpler usage:

=== "macOS"

    ```shell
    mv esdm-darwin-arm64 esdm
    ```

=== "Linux"

    ```shell
    mv esdm-linux-arm64 esdm
    ```

=== "Windows"

    ```shell
    ren esdm-windows-arm64.exe esdm.exe
    ```

!!! info "Note for x86 Users"

    **Replace `arm64` with `amd64`** in the file name if you are using an x86 architecture.

### Making the Binary Executable

The ESDM binaries are **not signed**. The steps below get an unsigned binary past your operating system's default protections.

=== "macOS"

    Files downloaded from the internet are marked with a quarantine attribute by macOS, which prevents them from being run. **Remove the quarantine attribute**:

    ```shell
    xattr -d com.apple.quarantine esdm
    ```

    **Make the binary executable**:

    ```shell
    chmod a+x esdm
    ```

=== "Linux"

    **Make the binary executable**:

    ```shell
    chmod a+x esdm
    ```

=== "Windows"

    Files downloaded from the internet are marked with a "Mark of the Web", which makes SmartScreen block the binary as coming from an unknown publisher. **Remove that mark**:

    ```shell
    Unblock-File .\esdm.exe
    ```

    If SmartScreen still shows a **"Windows protected your PC"** dialog on first run, choose **More info** and then **Run anyway**.

### Verifying the Installation

After renaming and adjusting permissions, verify the installation by checking the version:

=== "macOS"

    ```shell
    ./esdm version
    ```

=== "Linux"

    ```shell
    ./esdm version
    ```

=== "Windows"

    ```shell
    esdm version
    ```

This command will display the installed version of ESDM. If the version number matches your expectation, the installation was successful.

## Where to Go Next

- **[Your First Model with AI](/getting-started/your-first-model-with-ai.md)** drafts a tiny ESDM model through a short conversation with a coding agent that reads the schemas for you.
- **[Your First Model by Hand](/getting-started/your-first-model.md)** walks through writing and linting a tiny ESDM model from scratch, artifact by artifact.
- **[Editor Support](/getting-started/editor-support.md)** sets up your editor so it offers autocomplete and validation against the ESDM schemas.
- **[Concepts](/concepts/overview.md)** is the canonical introduction to the ESDM vocabulary – Aggregates, Events, Commands, and the rest.
