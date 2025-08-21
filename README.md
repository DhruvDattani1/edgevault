# EdgeVault

**EdgeVault** is a lightweight, edge-friendly object storage tool written in Go.
It focuses on **secure, efficient, and streaming-safe file storage** using modern encryption (`ChaCha20-Poly1305`).

The goal is to make it practical for resource-constrained devices (like drones, IoT sensors, or edge servers) to store logs, images, or binary data **securely at rest**.

---

## 🔑 What’s Built So Far

### 1. **Encryption at Rest**

* Uses `ChaCha20-Poly1305` for authenticated encryption.
* Supports both **small files (inline)** and **large files (chunked streaming)**.
* Large files are split into chunks with independent nonces for safe parallelism.

### 2. **`put` Utility**

* `edgevault put <source_file>`
* Encrypts a file and saves it into the vault (`crypta/` directory).
* Handles both small and large files automatically:

  * Small files → encrypt whole buffer at once.
  * Large files → chunked encryption with a custom header (`EV1\x00`).

### 3. **`get` Utility**

* `edgevault get <object_name> <dest_file>`
* Retrieves and decrypts stored objects.
* Detects format automatically:

  * Chunked (`EV1\x00`) → stream decryption.
  * Small file → inline decryption.

### 4. **Atomic Storage Safety**

* Writes go to `.partial` files first.
* Only renamed to final object after successful encryption, avoiding corruption.

---

## 🛠️ Work in Progress

The project is still an **MVP in active development**.
Here’s what’s next on the roadmap:

* [ ] **Master Key Handling**

  * Currently hardcoded (`main.go`) as a 32-byte value.
  * Needs a more secure approach: env vars, config file, hardware secure module, or key derivation.

* [ ] **Object Management**

  * `list` — show stored objects.
  * `delete` — remove objects safely.

* [ ] **Metadata Support**

  * Store metadata (timestamp, size, type, checksums).

* [ ] **Integration**

  * Expose as a service/daemon or embed in edge apps (cameras, drones, sensors).
  * Potential for gRPC/HTTP API later.

---

## 🚧 Current Status

Right now, EdgeVault can:

* **Store files securely** (`put`)
* **Retrieve and decrypt files** (`get`)

Everything else (key management, metadata, cleanup utilities, integrations) is still **in progress**.
Think of this repo as a **secure storage skeleton** that will grow into a more complete edge-first storage layer.

---

