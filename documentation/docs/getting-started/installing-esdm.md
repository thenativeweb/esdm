# Installing ESDM

This guide shows how to **install ESDM** and verify that it's working correctly. It covers **macOS**, **Linux**, and **Windows**. At the end, you will have a **working `esdm` binary**, ready to lint and view your models.

ESDM is distributed as **pre-built binaries** for your operating system and CPU architecture. Download the binary, place it in your project, and you're done.

## Downloading the Latest Version

To get the latest version of ESDM, select your operating system and CPU architecture:

=== "macOS"

    | Architecture | Version | Download | Size | SHA256 |
    |--------------|---------|----------|------|--------|
    | ARM64 | 0.13.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.13.0/esdm-darwin-arm64)** | 8.5 MB | `3677161511b7ec01a9828988d5d73bbd12ec3849f863a8d5b0c4dea74803f5b4` |
    | x86 | 0.13.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.13.0/esdm-darwin-amd64)** | 9.1 MB | `9409e3ab094d8e85eeb88220184c5c0433a3f080259830c6e99d1c7064f736db` |

=== "Linux"

    | Architecture | Version | Download | Size | SHA256 |
    |--------------|---------|----------|------|--------|
    | ARM64 | 0.13.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.13.0/esdm-linux-arm64)** | 8.3 MB | `770046e3e5cdf2973d1ceb41a0f064f065f92719a0e13f0ce6d08d272e561f64` |
    | x86 | 0.13.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.13.0/esdm-linux-amd64)** | 9.0 MB | `32a7a8974c59169b52e471b7ab9fc43b7e8b93b0fe4b2fe72eb55b94d80e09f4` |

=== "Windows"

    | Architecture | Version | Download | Size | SHA256 |
    |--------------|---------|----------|------|--------|
    | ARM64 | 0.13.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.13.0/esdm-windows-arm64.exe)** | 8.5 MB | `7ccbd46852782388e7056bcf549e95dd729bb399b6e12f4d34f7ff88a8321f22` |
    | x86 | 0.13.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.13.0/esdm-windows-amd64.exe)** | 9.3 MB | `93ecf1e832ad311cb83d9632544ff29e0c8f7a54798c3d6f51d9f4868afc0f68` |

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
