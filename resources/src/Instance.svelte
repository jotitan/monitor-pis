<script>
	import {current_instance} from './store.js';
	export let instance = {};

	const update = ()=> {
		current_instance.update(()=>instance.name)
	}

    const formatSuffix = key => {
        switch(key){
            case 'temperature' :return ' °C';
            case 'memory':
            case 'cpu':
            case 'disk':return ' %';
            default :return '';
        }
    }

    const showMetric = (name,value)=>{
	    if( name === "name"){
	        return "";
        }
	    return `${name} : ${parseFloat(value).toFixed(2)} ${formatSuffix(name)}`
    }
</script>

<style>
	.instance {
		border:solid 1px darkgrey;
		width:300px;
		display:inline-block;
		margin-left:20px;
		padding:10px;
		cursor:pointer;
	}
	.title {
		color:darkgray;
		font-size:22px;
		text-align:center;
		padding:10px;
		font-weight:bold;
	}
</style>

<div class="instance" on:click="{update}">
<div class="title">Instance {instance.name}</div>
    {#each Object.keys(instance) as id}
        <div>{showMetric(id,instance[id])}</div>
    {/each}
</div>