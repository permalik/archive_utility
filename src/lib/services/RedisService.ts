import { UPSTASH_REDIS_REST_URL, UPSTASH_REDIS_REST_TOKEN } from '$env/static/private';
import { Redis } from '@upstash/redis';
import type { Repo } from '../schema';

const redis = new Redis({
	url: UPSTASH_REDIS_REST_URL,
	token: UPSTASH_REDIS_REST_TOKEN
});

export class RedisService {
	constructor() {
		// Create instance of RedisService
	}
	public async SetRepos(repos: Repo[]) {
		for (const repo of repos) {
			const jsonRepo: JSON = <JSON>(<unknown>{
				repoID: repo.repoID,
				name: repo.name,
				description: repo.description || 'n/a',
				htmlURL: repo.htmlURL || 'n/a',
				homepage: repo.homepage || 'n/a',
				tag: repo.tag || [],
				createdAt: repo.createdAt,
				updatedAt: repo.updatedAt
			});
			await redis.set(repo.name, jsonRepo);
		}
	}
	public async SetHashes(repos: Repo[]) {
		const repoNames = repos.map((repo) => repo.name);
		await redis.hset(`currentRepos`, { repos: JSON.stringify(repoNames) });
	}
}
