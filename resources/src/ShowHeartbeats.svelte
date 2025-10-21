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
	const showChart = ()=> current_instance.update(()=>"heartbeat")

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
		border: solid 2px #e5e5e5;
	}
	.block {
		background-color:#262a33;
		padding:10px;
	}
</style>

<div class='block {isSelected ? "selected":""}' on:click="{showChart}" >
	<div class="title">Monitoring services</div>

	{#await promise}
		loading...
	{:then heartbeats}
		{#each heartbeats as hb}
			<Heartbeat heartbeat={hb}/>
		{/each}
	{/await}
</div>

