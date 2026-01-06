import { NodePgDatabase } from "drizzle-orm/node-postgres";
import { Pool } from "pg";
import { drizzle } from "drizzle-orm/node-postgres";
import { ResponseError, createError, whereAmI } from "../../internal/infrastructure/custom_error";

export class TxManager {
  private pool: Pool;

  constructor(pool: Pool) {
    this.pool = pool;
  }

  async withTransaction<T>(
    fn: (db: NodePgDatabase) => Promise<{ result: T | null; errors: ResponseError[] }>
  ): Promise<{ result: T | null; errors: ResponseError[] }> {
    const client = await this.pool.connect();

    try {
      await client.query("BEGIN");

      const txDb = drizzle(client);

      const { result, errors } = await fn(txDb);

      if (errors.length > 0) {
        await client.query("ROLLBACK");
        return { result: null, errors };
      }

      await client.query("COMMIT");
      return { result, errors: [] };
    } catch (error) {
      await client.query("ROLLBACK");
      const message = error instanceof Error ? error.message : "Transaction failed";
      return {
        result: null,
        errors: [createError(500, whereAmI(), "Transaction Failed", message)],
      };
    } finally {
      client.release();
    }
  }
}

export function newTxManager(pool: Pool): TxManager {
  return new TxManager(pool);
}
