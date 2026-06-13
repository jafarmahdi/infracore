import { readFile, readdir } from "node:fs/promises";
import { join, resolve } from "node:path";
import { pathToFileURL } from "node:url";

const runtimeRoot =
  process.env.INFRACORE_RUNTIME ??
  join(process.env.LOCALAPPDATA, "InfraCore", "runtime");
const pgModule = join(runtimeRoot, "node_modules", "pg", "lib", "index.js");
const migrationsDir = resolve("migrations", "postgres");
const { default: pg } = await import(pathToFileURL(pgModule).href);

const client = new pg.Client({
  host: "localhost",
  port: 5432,
  database: "infracore",
  user: "infracore",
  password: "infracore_secret",
});

// Embedded PostgreSQL has no TimescaleDB extension. Local development keeps
// the metric tables as regular PostgreSQL tables while production uses the
// original migrations unchanged.
function withoutTimescale(sql) {
  return sql
    .replace(
      /CREATE EXTENSION IF NOT EXISTS "timescaledb" CASCADE;[^\r\n]*/g,
      "",
    )
    .replace(/SELECT create_hypertable\([\s\S]*?\);\s*/g, "")
    .replace(/SELECT add_compression_policy\([\s\S]*?\);\s*/g, "")
    .replace(/SELECT add_retention_policy\([\s\S]*?\);\s*/g, "");
}

await client.connect();
await client.query(`
  CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
  )
`);

const files = (await readdir(migrationsDir))
  .filter((name) => name.endsWith(".up.sql"))
  .sort();

for (const file of files) {
  const applied = await client.query(
    "SELECT 1 FROM schema_migrations WHERE version = $1",
    [file],
  );
  if (applied.rowCount > 0) continue;

  const sql = withoutTimescale(
    await readFile(join(migrationsDir, file), "utf8"),
  );
  await client.query("BEGIN");
  try {
    await client.query(sql);
    await client.query(
      "INSERT INTO schema_migrations (version) VALUES ($1)",
      [file],
    );
    await client.query("COMMIT");
    console.log(`Applied ${file}`);
  } catch (error) {
    await client.query("ROLLBACK");
    throw new Error(`Migration ${file} failed`, { cause: error });
  }
}

await client.end();
console.log("InfraCore migrations are up to date");
