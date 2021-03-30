<script>
	export let heartbeat = {};
</script>

<style>
	.heartbeat > div {
		display:inline-block;
	}
	.heartbeat {
		text-align:right;
	}
	.up {
		color:green;
		font-weight:bold;
	}
	
	.down {
		color:red;
	}
</style>

<div class="heartbeat">
	<div style="width:200px">{heartbeat.Name}</div>
	<div class="{heartbeat.IsUp ? 'up':'down'}" style="margin-left:20px;">{heartbeat.IsUp ? 'UP':'DOWN'}</div>
</div>