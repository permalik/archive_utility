import { sql } from 'drizzle-orm';
import type { Repo } from '../schema/types/github.types';
import { ReposTable } from '../schema/models/ReposTable.model';
import { OctokitService } from './Octokit.service';
import { RedisService } from './RedisService';
import { DrizzleService } from './Drizzle.service';
const drizzleService = new DrizzleService();
const db = await drizzleService.ConnectDB();

export class RepoService {
	constructor() {
		// Create instance of RepoService
	}
	public async Seed() {
		await db.execute(sql`
	      CREATE TABLE IF NOT EXISTS repos (
	        id SERIAL PRIMARY KEY
	        , repo_id INTEGER NOT NULL
	        , name VARCHAR(150) NOT NULL
	        , description TEXT
	        , html_url VARCHAR(255) NOT NULL
	        , homepage VARCHAR(255)
	        , tag JSON
	        , created_at VARCHAR(30)
	        , updated_at VARCHAR(30)
	      )
	    `);
		console.log(`Created "repos" table`);
		const octokitService = new OctokitService();
		const updatedRepos: Repo[] = await octokitService.FetchRepos();
		const newRepos = await db
			.insert(ReposTable)
			.values(updatedRepos)
			.returning({
				repoID: ReposTable.repoID,
				name: ReposTable.name,
				description: ReposTable.description || 'n/a',
				htmlURL: ReposTable.htmlURL || 'n/a',
				homepage: ReposTable.homepage || 'n/a',
				tag: ReposTable.tag || [],
				createdAt: ReposTable.createdAt,
				updatedAt: ReposTable.updatedAt
			});
		const redisService = new RedisService();
		await redisService.SetRepos(newRepos);
		await redisService.SetHashes(newRepos);
	}
	public async FetchRepos() {
		let repos;
		try {
			repos = await db.select().from(ReposTable);
			return repos;
		} catch (error: any) {
			if (error.message === `relation "repos" does not exist`) {
				console.log('Table does not exist. Creating and seeding now.');
				await this.Seed();
				repos = await db.select().from(ReposTable);
				return repos;
			} else {
				throw error;
			}
		}
	}
	public async refreshRepos() {
		await db.delete(ReposTable);
		await this.Seed();
	}
}
