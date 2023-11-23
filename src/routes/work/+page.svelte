<script lang='ts'>
  export let data;
</script>

<section>
  <h2>Work</h2>
  <ul>
    {#each data.result as repo}
      <li>
        <a href={repo.htmlURL}>{repo.name}</a>
      </li>
    {/each}
  </ul>
</section>

<style lang='scss'>
  ul {
    list-style: none;
    & li:not(:last-of-type) {
      margin: 0 0 1rem 0;
    }
  }
</style>