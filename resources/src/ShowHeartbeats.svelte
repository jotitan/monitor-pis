<script>
	// Load instances from server
	import Heartbeat from './Heartbeat.svelte'
	import {current_instance} from './store.js';

	const loadHeartbeatsList = async ()=>{
		return new Promise((resolve,reject)=> {
			fetch("/heartbeats")
					.then(d=>d.json().then(resolve))
		});
	}

	let promise = loadHeartbeatsList();
	let isSelected = false;
	const showChart = ()=> {
		current_instance.update(()=>"heartbeat")
	}

	current_instance.subscribe(ci=>isSelected = ci === "heartbeat")

</script>

<style>
	.title {
		color: #737373;
		font-size:22px;
		font-weight:bold;
		text-align:right;
		padding-bottom:10px;
		cursor:pointer;
	}
	.block.selected {
		border-color: #ff8c00;
		border-width: 2px;
	}
	.block {
		border:solid 1px black;
		padding:10px;
	}
</style>

<div class='block {isSelected ? "selected":""}'>
	<div on:click="{showChart}" class="title">Monitoring services</div>

	{#await promise}
		loading...
	{:then heartbeats}
		{#each heartbeats as hb}
			<Heartbeat heartbeat={hb}/>
		{/each}
	{/await}

</div>

