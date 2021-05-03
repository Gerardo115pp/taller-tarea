<script>
    export let input_id = "";
    export let input_label;
    export let input_type = "text";
    export let input_placeholder;
    export let onKeypressed;
    export let onBlur;
    
    export let initial_value;
    export let min = 0;
    export let max;




    const handleOutsideClickDetected = e => {
        if(e.currentTarget === e.target) {
            e.target.getElementsByTagName("input")[0].focus()
        }
    }
</script>

<style>
    .input-container {
        cursor: text;
        display: flex;
        border: 2px solid var(--theme-color);
        border-radius: 9999px;
        padding: 1.5vh 4.2vh;
    }   

    .input-container input {
        font-family: var(--main-font);
        display: flex;
        background: none;
        border: none;
        color: var(--theme-color);
        font-size: 1.2rem;
        padding: 0;
        align-items: center;
        outline: none;
    }

    .input-container input[type='number'] {
        width: 3rem;
        text-align: center;
    }

    .input-container input::placeholder {
        font-family: 'Rhodium Libre';
        color: var(--dimtheme-color);
        text-transform: capitalize;
    }

</style>

<div on:click={handleOutsideClickDetected} class="input-container">
    {#if input_label !== undefined}
        <label for={input_id}>{input_label}:</label>
    {/if}
    {#if input_type!== "number"}
        <input id={input_id} type={input_type}
            placeholder={input_placeholder}
            on:keydown={onKeypressed}
            on:blur={onBlur}
        />
    {:else}
        <input id={input_id} type={input_type}
            placeholder={input_placeholder}
            on:keydown={onKeypressed}
            on:blur={onBlur}
            value={initial_value}
            min={min}
            max={max}
        />
    {/if}

</div>