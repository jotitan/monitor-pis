<script>
	import ShowInstances from './ShowInstances.svelte'
	import ShowHeartbeats from './ShowHeartbeats.svelte'
	import Chart from './Chart.svelte'
</script>

<style>
	h1 {
		margin-top:0px;
		font-size:1.6em;
	}
	.body {
		color:white;
	}
</style>
<div class="body">
<h1>Raspberry home monitoring</h1>

<div style="float:left;width:70%">
	<ShowInstances/>
</div>

<div style="float:left;width:30%">
	<ShowHeartbeats/>
</div>

<Chart/>
</div>



