<script>
	let name = 'world nnnn';
	import Boucle from './Instance.svelte'
	
	let list = [1,2,3];
</script>

<style>
	h1 {
		color:red
	}
</style>

<h1>Hello {name}!</h1>

{#each list as id}
<Boucle/> sd
{/each}