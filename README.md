# Privacy-Preserving Gateway

A client-server capstone project providing secure access to a privacy-preserving cloud storage platform, featuring a **Go backend** with native self-bootstrapping cryptographic keys/TLS and a **SvelteKit / TypeScript frontend**.

---

## Prerequisites

Ensure you have installed:
- **Go** (v1.26+)
- **Node.js** (v18+ & npm)
- **Git** & **Docker** (optional for containerized deployment)

---

## Getting Started (Fresh Clone)

Clone the repository and build the entire stack with a single command:

```bash
git clone [https://github.com/0x03ff/privacy-preserving-gateway.git](https://github.com/0x03ff/privacy-preserving-gateway.git)
cd privacy-preserving-gateway

# Installs Node modules, compiles WASM, builds SvelteKit, and compiles Go binary
make
```

## Running the Application

### 1. Run the Backend Server
Execute the Go application directly from the **project root directory** so that relative paths targeting the `var/` directory resolve correctly:

```bash
go run cmd/main.go
```

Note on drop_flag (cmd/main.go):
- `false (Default):` The SQLite database (var/db/app.db) persists records, files, and user data across restarts.
- `true:` Wipes and reinitializes the database schema on every startup.

### 2. Frontend Development Server (SvelteKit)

To run the SvelteKit frontend in live-reload development mode:
```bash
# Install dependencies (if not already done via make)
npm install

# Start Vite dev server
npm run dev

# Or start and open in browser automatically
npm run dev -- --open
```
#### Available NPM Scripts
- `npm run dev` — Start Vite development server.
- `npm run build` — Create a production static build.
- `npm run preview` — Locally preview the production build.
- `npm run check` — Run SvelteKit type checking (svelte-check).
- `npm run lint` — Lint code with ESLint and Prettier.
- `npm run format` — Automatically format code with Prettier.