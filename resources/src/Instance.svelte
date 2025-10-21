<script>
	import {current_instance} from './store.js';
	export let instance = {};
    let isSelected = false;
	const update = ()=> {
		current_instance.update(()=>instance.name)
	}

	current_instance.subscribe(ci=>isSelected = ci === instance.name);

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
        background-color:#262a33;
		min-width:160px;
		display:inline-block;
		margin-left:20px;
		padding:10px;
		cursor:pointer;
        vertical-align: top;
	}
	.title {
		color: #737373;
		font-size:22px;
		text-align:center;
		font-weight:bold;
	}
    .metric::first-letter {
        text-transform: capitalize;
    }

    .selected {
        border: solid 2px #e5e5e5;
    }
</style>

<div class='instance {isSelected ? "selected":""}' on:click="{update}">
<div class="title">Instance {instance.name}</div>
    {#each Object.keys(instance) as id}
        <div class="metric">{showMetric(id,instance[id])}</div>
    {/each}
</div>