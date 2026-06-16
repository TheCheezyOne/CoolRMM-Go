# CoolRMM-Go

## Project Overview

CoolRMM-Go (eventually just "CoolRMM") is a lightweight, self-built RMM tool
replacing Pulseway/Kaseya at work. Written in Go. Designed to monitor ~180
Windows 11 endpoints across three groups: CorpA, CorpB, and CorpC (internal).

Build order: Agent → Server → Dashboard → DB.

### Customers / Groups
- CorpA
- CorpB
- CorpC (internal group; parent company: CorpD)
- Each customer will eventually have sub-groups and tags.

### Agent (v0.x.x)
- Single self-contained .exe, no installer
- Phones home every minute to the server
- First feature: is it online? Who is logged in?
- APIs for security stack tools added one at a time

### Server (v0.x.x, after agent)
- Runs on primary sandbox device (Windows 11 Pro) at work
- Receives check-ins from agents, stores data in DB

### Dashboard (v0.x.x, after server)
- Dark theme
- Devices listed vertically, line by line
- Columns: display name, power status dot (green/yellow/red), logged user
- Top-left: "CoolRMM" name
- Below name: status bar — total devices, green count, red count, yellow count
- No cards
- Easy to read, sort, and scan

### DB (v0.x.x, after dashboard)
- DB choice TBD — will evaluate and suggest when we get there (lightweight preferred)

---

## Development Rules

These rules govern all development on this project. The agent must follow them at all times.

### Comments
- Every file must have a **multiline header** explaining what the file is.
- Every file must have a **multiline footer** explaining how to run it or other pertinent info.
- All other comments: **one line if possible**. Choose clarity over a single line only when necessary.
- Add a new comment for every new block of code.

### Code Quality
- Follow **DRY** (Don't Repeat Yourself) — no duplicated logic.
- Follow **WORM** (Write Once, Read Many) — write for readability first.
- Keep code clean and easy to follow for a novice programmer.

### Dependencies
- **Always ask before adding any new dependency or package.**

### Refactoring
- **Never refactor code that was not explicitly requested to be changed.**

### Error Handling
- Always handle errors explicitly. **No silent failures.**

### Naming Conventions
- **snake_case only** for all identifiers.

### Version Control
- Using **Git and GitHub**.
- Always summarize changes before committing.

### Testing
- **Write tests for every new function.**

### Data
- **No placeholder or mock data in production code.** Be explicit when something is not real data.

### Type Checking
- **Always typecheck before declaring something done.**

### Language
- This is a **Go project**.

### Platform
- The **agent runs on Windows 11**. Always consider Windows paths and Windows service behavior.

### File Size
- **Keep files under 300 lines.** Split into multiple files if larger.

### Binary
- **Keep the agent binary self-contained — no installer, single .exe.**

### Pace
- **Start small — one feature at a time.** Confirm before moving to the next.
- "Small" means very small: an empty window is a feature. A color change is a feature.
- **Always show the plan before writing any code.**

### Versioning (v0.0.0)
- **Right digit (patch):** increments after each corrective iteration.
- **Middle digit (minor):** increments when a feature block is complete and we move on.
- **Left digit (major):** moves to `1` when the project is ready for deployment.
- Current version: **v0.0.0**

---

## Repository & Environment

- **GitHub URL:** https://github.com/TheCheezyOne/CoolRMM-Go
- **Local path:** `/home/cheez/Public/CoolRMM/CoolRMM-Go/`
- **Auth method:** SSH (ED25519 key already configured)
- **Dev environment:** Linux (Replit at home, Linux at home)
- **Work environment:** Windows 11 (continuing development on-site)
- **Target deployment:** Windows 11 x64

### Cross-Platform Rules
- Always use Go's `filepath` package for paths — never hardcode `/` or `\`
- Cross-compile with `GOOS=windows GOARCH=amd64`
- Flag any code that cannot be tested on Linux

### Git Sync Note
Replit's sandbox blocks all `.git/config` modifications — remotes cannot be added programmatically.
Sync is manual: copy changed files from Replit to local machine, then commit and push via SSH from there.
If you edit `replit.md` locally and push, paste the relevant section into chat and it will be updated here.

---

## User Preferences
- Write code that can be followed easily.
- User is a mid-level developer and sysadmin — will proofread every line.
- Keep code tight, clean, and not over-engineered.
- Profanity is not only allowed but encouraged — keeps the vibe right.
- User tends to over-explain — do not let that leak into the code.

---

## Notes

### Agent Deployment (per machine)
1. Copy `coolrmm.conf.example` → `coolrmm.conf` in the same folder as `coolrmm-agent.exe`
2. Edit `coolrmm.conf` — set `server_url` to the real server IP and port
3. Run `coolrmm-agent.exe`

### Wishlist (future features, no timeline)
- **Agent self-healing via Windows Service** — register agent as a Windows Service using `golang.org/x/sys/windows/svc` so the OS auto-restarts it on failure. Requires Windows-specific code (untestable on Linux) and a new dependency. For now, manual restart via Pulseway remote terminal is sufficient.

### Milestones
- First successful live end-to-end test on Windows 11 — agent phoning home to server, check-ins logging every 60 seconds ✓ `061426@1409`
- CPU usage % confirmed reporting live from Windows 11 agent to server console and database ✓ `061626@1053`
- Web dashboard live — dark theme, device list with status dots, cpu%, auto-refresh ✓ `061626@1101`
- Multi-device confirmed — two endpoints reporting live to dashboard simultaneously ✓ `061626@1151`
- RAM % and Disk (C:) % confirmed live on dashboard — both endpoints reporting ✓ `061626@1203`

### Version Log

**Agent**
- agent v0.0.0: Initial scaffold — compiles, runs, prints startup message
- agent v0.1.0: Reads hostname and logged-in user, prints to console
- agent v0.2.0: Defines check-in payload struct (hostname, logged user, UTC timestamp)
- agent v0.3.0: POSTs check-in payload as JSON to server /checkin endpoint over HTTP
- agent v0.4.0: Runs a 60-second check-in loop — sends immediately on startup, logs errors without dying
- agent v0.5.0: server_url loaded from coolrmm.conf at startup — no more hardcoded localhost
- agent v0.6.0: CPU usage % added to payload via gopsutil (500ms sample, combined across all cores)
- agent v0.7.0: RAM usage % and disk usage % (C:\) added to payload
- agent v0.8.0: Uptime (seconds since last boot) added to payload

**Server**
- server v0.0.0: (no standalone scaffold — server started at first feature)
- server v0.1.0: Minimal HTTP server — POST /checkin receives agent payload, prints to console, responds 200 OK
- server v0.2.0: Persists check-ins to SQLite DB — open_db(), create_schema(), insert_checkin()
- server v0.3.0: cpu_percent added to payload, DB schema, and console output — existing DBs auto-migrated via ALTER TABLE
- server v0.4.0: Web dashboard — GET / serves dark theme HTML page, GET /devices returns device JSON; status dots green/yellow/red by check-in age
- server v0.5.0: RAM and Disk (C:\) columns added to dashboard and database
- server v0.6.0: Uptime column added to dashboard and database
