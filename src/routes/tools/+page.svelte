<script lang="ts">
	export let data;

	async function refreshRepos() {
		await fetch('/api/repos', {
			method: 'POST'
		});
	}
</script>

<section>
	<button on:click={refreshRepos}>Refresh Repos</button>
	<div>
		<a href="/work">View {data.result.length} Projects</a>
	</div>
</section>

<style>
	button {
		margin: 0 0 1rem 0;

		background-color: transparent;
		color: var(--primary-light);
		border: 2px solid var(--primary-light);
		border-radius: 5px;

		font-weight: 500;
	}
	button:hover {
		background-color: var(--primary-light);
		color: var(--primary-dark);

		cursor: pointer;
	}
</style>
