<script>
	// Load instances from server
	import Instance from './Instance.svelte';
	
	const loadInstancesList = async ()=>{
		return new Promise(resolve=>{
			fetch('http://localhost:9000/instances')
				.then(r=>r.json()
					.then(data=>
						resolve(Object.keys(data).map(name=>{
							let instance = {name:name};
							Object.keys(data[name]).map(metric=>instance[metric] = data[name][metric]);
							return instance;
						}))));								
		});
	}
	
	let promise = loadInstancesList();


</script>

<style>
	
	
</style>

{#await promise}
loading...
{:then instances}
	{#each instances as instance}
		<Instance instance={instance}/>
	{/each}
{/await}