import { Octokit } from 'octokit';
import type { Repo } from '../schema/types';

const octokit = new Octokit();

export class OctokitService {
	constructor() {
		// Create new OctokitService instance
	}
	public async FetchRepos() {
		const { data } = await octokit.request('GET /users/{username}/repos', {
			username: 'permalik',
			sort: 'created'
		});

		const repos = data.map((repo) => {
			const newRepo: Repo = {
				repoID: repo.id,
				name: repo.name,
				description: repo.description,
				createdAt: repo.created_at,
				updatedAt: repo.updated_at,
				htmlURL: repo.html_url,
				homepage: repo.homepage,
				tag: repo.topics
			};
			return newRepo;
		});
		return repos;
	}
}
