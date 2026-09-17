# Installing ESDM

This guide shows how to **install ESDM** and verify that it's working correctly. It covers **macOS**, **Linux**, and **Windows**. At the end, you will have a **working `esdm` binary**, ready to lint and view your models.

ESDM is distributed as **pre-built binaries** for your operating system and CPU architecture. Download the binary, place it in your project, and you're done.

## Downloading the Latest Version

To get the latest version of ESDM, select your operating system and CPU architecture:

=== "macOS"

    | Architecture | Version | Download | Size | SHA256 |
    |--------------|---------|----------|------|--------|
    | ARM64 | 0.15.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.15.0/esdm-darwin-arm64)** | 9.1 MB | `ea11a2629447b54a8475b67d30d01215b74fc61fe563c4715119c73d6f823f88` |
    | x86 | 0.15.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.15.0/esdm-darwin-amd64)** | 9.8 MB | `b4e0009243fb5ed822fbade6e8f469395227feba95269b1e59afbe33fad14705` |

=== "Linux"

    | Architecture | Version | Download | Size | SHA256 |
    |--------------|---------|----------|------|--------|
    | ARM64 | 0.15.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.15.0/esdm-linux-arm64)** | 8.8 MB | `5076335000631e8701cf07fc610dcc86c806c7422fd86ad68eea27b212820019` |
    | x86 | 0.15.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.15.0/esdm-linux-amd64)** | 9.6 MB | `db76e45729e5318023120143356df461dbddfbedf2108111c0d8b7b595f42a84` |

=== "Windows"

    | Architecture | Version | Download | Size | SHA256 |
    |--------------|---------|----------|------|--------|
    | ARM64 | 0.15.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.15.0/esdm-windows-arm64.exe)** | 9.0 MB | `8cdd0a01f740e4d1833f839ad96865bf202b748c18550659b873ff77a03b39ca` |
    | x86 | 0.15.0 | **[Download](https://esdm.s3.fr-par.scw.cloud/0.15.0/esdm-windows-amd64.exe)** | 9.9 MB | `0f82f30bbe73f0c68ec6633851599f3a1cf2c1f1e663cb6815b7b13f6abda5e7` |

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
