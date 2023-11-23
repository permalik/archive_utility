import { RepoService } from '../../../lib/services';

export async function POST() {
	const repoService = new RepoService();
	await repoService.refreshRepos();
	return new Response();
}
