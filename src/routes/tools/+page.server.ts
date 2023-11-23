import { RepoService } from '../../lib';

export async function load() {
	const repoService = new RepoService();
	const response = await repoService.FetchRepos();
	const data = { result: response };
	return data;
}
