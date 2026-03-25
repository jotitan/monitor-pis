<script>
	export let heartbeat = {};
	const formatName = name => name.replaceAll("_"," ")
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
	.heartbeat > div.name {
		width:100px;
	}
	.name::first-letter {
		text-transform: capitalize;
	}
	.status {

	}

</style>

<div class="heartbeat">
	<div class="name">{formatName(heartbeat.Name)}</div>
	<div class="status">{heartbeat.IsUp?'🟢':'🔴'}</div>
</div>