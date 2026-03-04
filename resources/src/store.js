import {writable} from 'svelte/store';


export let current_instance = writable('');
export let current_date = writable('');

export let refresh_date = writable(false);