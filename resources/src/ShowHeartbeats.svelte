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
	
	const showChart = ()=> {
		current_instance.update(()=>"heartbeat")
	}
	
	
	
</script>

<style>
	.title {
		font-size:22px;
		font-weight:bold;
		text-align:right;
		padding-bottom:10px;
		cursor:pointer;
	}
</style>

<div on:click="{showChart}" class="title">Monitoring services</div>

{#await promise}
loading...
{:then heartbeats}
	{#each heartbeats as hb}
		<Heartbeat heartbeat={hb}/>
	{/each}
{/await}



