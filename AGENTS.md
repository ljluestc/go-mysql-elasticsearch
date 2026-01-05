# go-mysql-elasticsearch — Agent Guide
 
## Project overview
 
`go-mysql-elasticsearch` is a service that syncs data from MySQL into Elasticsearch.
 
From `README.md`:
 
- It uses `mysqldump` for the initial full sync (when available).
- It then follows MySQL binlog for incremental updates.
 
## Tech stack
 
- **Language**: Go (`go.mod` module `github.com/siddontang/go-mysql-elasticsearch`).
- **Config format**: TOML (example at `etc/river.toml`).
- **External dependencies at runtime** (per README / Dockerfile):
  - MySQL/MariaDB (binlog replication; `row` format required)
  - Elasticsearch
  - `mysqldump` (optional but recommended for initial snapshot)
 
## Repo layout
 
- `cmd/go-mysql-elasticsearch/main.go`: CLI entrypoint; parses flags, loads config, runs the river.
- `river/`: core syncing logic and configuration types.
- `etc/river.toml`: example configuration.
- `elastic/`, `river/`, `cmd/`: main implementation packages.
- `Dockerfile`: container build (installs `tini` and `mariadb-client`).
 
## Build and test commands
 
The `Makefile` is the canonical entry point.
 
- **Build**:
 
```bash
make
# or
make build
```
 
This produces `bin/go-mysql-elasticsearch`.
 
- **Test**:
 
```bash
make test
```
 
## Running
 
From `README.md`, the typical workflow is:
 
- Create tables in MySQL.
- Create Elasticsearch indices/mappings.
- Configure `etc/river.toml`.
- Start the service:
 
```bash
./bin/go-mysql-elasticsearch -config=./etc/river.toml
```
 
Useful CLI flags (from `cmd/go-mysql-elasticsearch/main.go`):
 
- `-config` (defaults to `./etc/river.toml`)
- `-my_addr`, `-my_user`, `-my_pass`
- `-es_addr`
- `-data_dir`
- `-server_id`
- `-flavor` (mysql/mariadb)
- `-exec` (mysqldump path)
- `-log_level`
 
## Runtime architecture (high level)
 
- `cmd/.../main.go` loads TOML config via `river.NewConfigWithFile`.
- It creates a `river.River` via `river.NewRiver(cfg)` and runs it in a goroutine.
- The river performs initial loading (optionally using `mysqldump`) and then follows binlog for incremental sync.
- The process runs until it receives a signal or the river context ends.
 
## Development conventions
 
- Configuration examples live in `etc/river.toml`.
- The project’s `README.md` documents supported MySQL/ES versions and required binlog settings.
 
## Security considerations
 
- The config contains database credentials (`my_user`, `my_pass`) and Elasticsearch credentials (`es_user`, `es_pass`). Do not commit real secrets.
- The tool uses replication privileges; restrict MySQL permissions to the minimum needed.

## AI Agent Workflow

### 1. Requirements Discovery
- **Primary Source**: `PRD.md` (Always prioritize this if present).
- **Secondary**: `requirements.txt`, `README.md`, or specific task files.
- **Goal**: Understand the full scope before writing code.

### 2. Implementation Protocol
- **Branching**: Work on a dedicated feature branch (e.g., `feat/implementation-details`).
- **Development**:
  - Analyze code structure.
  - Implement changes in `src/` or relevant directories.
  - Adhere to existing code style.
- **Verification**:
  - Run build commands (see above).
  - Run test suite (see above).
  - Ensure no regressions.

### 3. Delivery
- **Commit**: Use conventional commits (e.g., `feat: ...`, `fix: ...`).
- **PR Creation**:
  - Push branch: `git push -u origin <branch-name>`
  - Create a Pull Request against the main branch.
  - Summary: Link to `PRD.md` requirements solved.

## Task Implementation
1. **Analyze Requirements**: Refer to `README.md` for detailed feature specifications and system design.
2. **Implementation**: Modify source code in the respective directories (e.g., `src/`, `internal/`).
3. **Verification**: Run provided build and test commands (see above) to ensure correctness.
4. **Push Changes**:
   - Commit changes: `git commit -m "feat: implement <feature>"`
   - Push to remote: `git push origin <branch-name>`
