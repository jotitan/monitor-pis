<script>
	import {current_instance} from './store.js';
	export let instance = {};
    let isSelected = false;
	const update = ()=> {
		current_instance.update(()=>instance.name)
	}

	current_instance.subscribe(ci=>{
	   isSelected = false;
	    if(ci === instance.name) {
            isSelected = true;
       }
    });

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
		color: #737373;
		font-size:22px;
		text-align:center;
		font-weight:bold;
	}
    .metric::first-letter {
        text-transform: capitalize;
    }

    .selected {
        border-color: #ff8c00;
        border-width: 2px;
    }
</style>

<div class='instance {isSelected ? "selected":""}' on:click="{update}">
<div class="title">Instance {instance.name}</div>
    {#each Object.keys(instance) as id}
        <div class="metric">{showMetric(id,instance[id])}</div>
    {/each}
</div>