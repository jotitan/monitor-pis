<script>
	export let heartbeat = {};
	const formatName = name => {
		return name.replaceAll("_"," ");
	}
</script>

<style>
	.heartbeat > div {
		display:inline-block;
		font-size:0.8em;
	}
	.heartbeat {
		text-align:right;
		display:inline-block;
	}
	.name::first-letter {
		text-transform: capitalize;
	}

</style>

<div class="heartbeat">
	<div style="width:180px" class="name">{formatName(heartbeat.Name)}</div>
	<div style="margin-left:10px;">{heartbeat.IsUp?'🟢':'🔴'}</div>
</div>