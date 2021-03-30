<script>
	let name = 'world nnnn';
	import ShowInstances from './ShowInstances.svelte'
	import ShowHeartbeats from './ShowHeartbeats.svelte'
	import Chart from './Chart.svelte'
</script>

<style>

</style>

<h1>Raspberry home monitoring</h1>

<div style="float:left;width:70%">
	<ShowInstances/>
</div>

<div style="float:left;width:30%">
	<ShowHeartbeats/>
</div>

<Chart/>



