export interface Repo {
	id?: number;
	repoID: number;
	name: string;
	description: string | null;
	htmlURL: string;
	homepage: string | null | undefined;
	tag: string[] | null | undefined | unknown;
	createdAt: string | null | undefined;
	updatedAt: string | null | undefined;
}
