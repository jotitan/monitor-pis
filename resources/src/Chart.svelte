<script>

    import {current_instance} from './store.js';

    import FusionCharts from 'fusioncharts';
    import Timeseries from 'fusioncharts/fusioncharts.timeseries';
    import Charts from 'fusioncharts/fusioncharts.charts';
    import Widgets from 'fusioncharts/fusioncharts.widgets';

    import CandyTheme from "fusioncharts/themes/fusioncharts.theme.candy";

    import SvelteFC, { fcRoot } from 'svelte-fusioncharts';

    // Add dependencies
    fcRoot(FusionCharts, Charts, Timeseries, Widgets, CandyTheme);

    let promise = null;
    let instance = "";

    const colorRange = {
        "color": [{
            "minValue": "0",
            "maxValue": "50",
            "code": "#F2726F"
        }, {
            "minValue": "50",
            "maxValue": "75",
            "code": "#FFC533"
        }, {
            "minValue": "75",
            "maxValue": "100",
            "code": "#62B58F"
        }]
    };

    const formatKey = key => key.replaceAll("_"," ")

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
            fetch(`/search?instance=${ci}&date=${currentDate}`)
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
                            tempData[d.Timestamp][pos+1] = d.Value;
                        });
                    });
                    series = Object.keys(tempData).sort((a,b)=>a-b).map(key=>tempData[key]);
                    resolve(createCharts(series,schema,axis));
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

    const countAvailability = (data,pos)=>{
        let count = data.map(d=>d[pos] != null ? d[pos]:1).reduce((a,b)=>a+b,0)
        return (count/data.length)*100;
    }

    // return an array with at least one chart
    const createCharts = (data,schema,axis) => {
        return instance === "heartbeat" ? createGauges(data,schema,axis) : createChart(data,schema,axis)
    }

    const createGauge = (data,pos,title)=> {
        const dials = {
            "dial": [{
                value: countAvailability(data,pos)
            }]
        };

        const dataSource = {
            "chart": {
                caption: "Availability",
                subcaption: formatKey(title),
                numberSuffix: "%",
                theme: "candy",
                lowerLimit: "0",
                upperLimit: "100",
                showValue: "1",
            },
            colorRange: colorRange,
            dials: dials
        };

        const chartConfigs = {
            type: 'angulargauge',
            width: 400,
            height: 200,
            dataFormat: 'json',
            dataSource
        };
        return chartConfigs
    }

    const createGauges = (data,schema)=> {
        let gauges = [];
        for(let pos = 1 ; pos < schema.length ; pos++){
            gauges.push(createGauge(data,pos,schema[pos].name));
        }
        return gauges;
    }

    const createChart = (data,schema,axis)=> {
        const fusionDataStore = new FusionCharts.DataStore(),
            fusionTable = fusionDataStore.createDataTable(data, schema);
        let dataSource = {
            caption: "Metrics history",
            showValues: "1",
            showPercentInTooltip: "0",
            enableMultiSlicing: "1",
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
            chart:{
                theme:"candy",
                multiCanvas: false
            },
            data: fusionTable,
        };

        return [{
            type: 'timeseries',
            width: '100%',
            height: 500,
            renderAt: 'chart-container',
            dataSource,
            chart:{
                multiCanvas: false
            },
            yAxis:axis
        }];
    }

    let currentDate = "";

    const setDate = date => {
        currentDate = date;
    }
</script>

<style>
    input,button {
        background-color: #262a33;
        color:white;
    }
</style>

<p>
    <input type="date" on:change={e=>setDate(e.target.value)}/>
    <button on:click={refresh}>Refresh</button>
</p>

{#if promise != null}
    {#await promise}
        Loading chart...
    {:then config}
        {#each config as chart}
            <div id="chart-container" style="display: inline-block;margin-left: 10px">
                <SvelteFC {...chart} />
            </div>
        {/each}
    {/await}

{/if}



