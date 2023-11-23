import { SUPABASE_PG_URI } from '$env/static/private';
import { drizzle } from 'drizzle-orm/postgres-js';
import postgres from 'postgres';

export class DrizzleService {
	constructor() {
		// Create instance of DrizzleService
	}
	public async ConnectDB() {
		const client = postgres(SUPABASE_PG_URI);
		const db = drizzle(client);
		return db;
	}
}
