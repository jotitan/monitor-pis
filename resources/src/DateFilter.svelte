<script>
    import {current_date, refresh_date} from './store.js';

</script>

<style>
    .options {
        float:right;
    }
    input,button {
        background-color: #262a33;
        color:white;
    }
    p {
        position: absolute;
        right: 0;
        top: 0;
    }
</style>

<p class="options">
    <input type="date" on:change={e=>current_date.set(e.target.value)}/>
    <button on:click={()=>refresh_date.update(b=>!b)}>Refresh</button>
</p>
