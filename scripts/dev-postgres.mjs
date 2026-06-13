import { existsSync } from "node:fs";
import { join } from "node:path";
import { pathToFileURL } from "node:url";

const runtimeRoot =
  process.env.INFRACORE_RUNTIME ??
  join(process.env.LOCALAPPDATA, "InfraCore", "runtime");
const modulePath = join(
  runtimeRoot,
  "node_modules",
  "embedded-postgres",
  "dist",
  "index.js",
);
const dataDir = join(runtimeRoot, "postgres-data");
const log = (message) => process.stdout.write(String(message));
const logError = (message) => process.stderr.write(String(message));

const { default: EmbeddedPostgres } = await import(
  pathToFileURL(modulePath).href
);

const postgres = new EmbeddedPostgres({
  databaseDir: dataDir,
  user: "infracore",
  password: "infracore_secret",
  port: 5432,
  persistent: true,
  initdbFlags: ["--encoding=UTF8", "--locale=C"],
  onLog: log,
  onError: logError,
});

// A persisted cluster is initialized only once; subsequent runs reuse its data.
if (!existsSync(join(dataDir, "PG_VERSION"))) {
  await postgres.initialise();
}

await postgres.start();

const client = postgres.getPgClient();
await client.connect();
const result = await client.query(
  "SELECT 1 FROM pg_database WHERE datname = $1",
  ["infracore"],
);
await client.end();

if (result.rowCount === 0) {
  await postgres.createDatabase("infracore");
}

console.log("InfraCore PostgreSQL is ready on localhost:5432");

// Keep the manager alive so its exit hook can stop PostgreSQL cleanly.
setInterval(() => {}, 60_000);
