<script>

    import {current_instance} from './store.js';

    import FusionCharts from 'fusioncharts';
    import Timeseries from 'fusioncharts/fusioncharts.timeseries';

    import SvelteFC, { fcRoot } from 'svelte-fusioncharts';

    // Add dependencies
    fcRoot(FusionCharts, Timeseries);

    let promise = null;
    let instance = "";

    const formatKey = key => key.replaceAll("_"," ")

    const formatValue = value => {
        return instance === "heartbeat" ? value *100 :value;
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
    let updateChart = ci => {
        instance = ci;
        return new Promise(resolve=>{
            fetch(`http://localhost:9000/search?instance=${ci}&date=${currentDate}`)
                .then(r=>r.json().then(data=>{
                    let series = [];
                    let schema = [];
                    schema.push({name:"Timestamp",type:"date",format:"%s"});
                    let tempData = {}
                    let size = Object.keys(data).length + 1;
                    let axis = [];
                    Object.keys(data).forEach((key,pos)=>{
                        schema.push({name:key,type:"number"});
                        axis.push( {plot: {value: key},title: formatKey(key),min:0,max:100,format:{suffix:formatSuffix(key)}});
                        data[key].forEach(d=>{
                            if(tempData[d.Timestamp] == null){
                                tempData[d.Timestamp] = new Array(size);
                            }
                            tempData[d.Timestamp][0] = d.Timestamp;
                            tempData[d.Timestamp][pos+1] = formatValue(d.Value);
                        });
                    });
                    series = Object.keys(tempData).sort((a,b)=>a-b).map(key=>tempData[key]);
                    resolve(createChart(series,schema,axis));
                }));


        });
    }

    current_instance.subscribe(value=>{
        if(value !== ""){
            promise = updateChart(value);
        }
    });

    const refresh = ()=> {
        promise = updateChart(instance);
    }

    const createChart = (data,schema,axis)=> {
        const fusionDataStore = new FusionCharts.DataStore(),
            fusionTable = fusionDataStore.createDataTable(data, schema);
        let dataSource = {
            caption: "Metrics history",
            showValues: "1",
            showPercentInTooltip: "0",
            enableMultiSlicing: "1",
            theme: "fusion",
            yAxis:axis,
            extensions: {
                standardRangeSelector: {
                    enabled: "0"
                }
            },
            tooltip: {
                outputTimeFormat: {
                    minute: "%d/%m/%y %H:%M:%S"
                }
            },
            navigator:{enabled:false},
            data: fusionTable
        };

        return {
            type: 'timeseries',
            width: '100%',
            height: (schema.length-1) * 160,
            renderAt: 'chart-container',
            dataSource,
            yAxis:axis
        };
    }

    let currentDate = "";

    const setDate = date => {
        currentDate = date;
    }
</script>

<style>

</style>

<input type="date" on:change={e=>setDate(e.target.value)}/>
<button on:click={refresh}>Refresh</button>

{#if promise != null}
    {#await promise}
        Loading chart...
    {:then config}
        <div id="chart-container" >
            <SvelteFC {...config} />
        </div>
    {/await}

{/if}



