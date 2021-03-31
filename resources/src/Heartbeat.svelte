<script>
	export let heartbeat = {};
	const formatName = name => {
		return name.replaceAll("_"," ");
	}
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

	.name::first-letter {
		text-transform: capitalize;
	}

</style>

<div class="heartbeat">
	<div style="width:200px" class="name">{formatName(heartbeat.Name)}</div>
	<div class="{heartbeat.IsUp ? 'up':'down'}" style="margin-left:20px;">{heartbeat.IsUp ? 'UP':'DOWN'}</div>
</div>