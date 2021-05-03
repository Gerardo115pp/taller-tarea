<script>
    import { onMount } from 'svelte';
    import { server_name } from '../server_data';
    import SearchBar from './SearchBar.svelte';
    import Button from './Button.svelte';
    import Select from './Select.svelte';
    import Input from './Input.svelte';
    export let onClose;
    export let getSessionKey = () => {};
    export let type_data = {};
    export let type_name = "Nothing at all";



    let type_content = [];
    let item_id =  undefined;
    let is_form_ready = false;

    const updateItemsContent = async () => {
        const headers = new Headers();
        headers.set("X-sk", getSessionKey())
        const promise = await fetch(`${server_name}/${type_name}?id=*`, {method: "GET", headers:headers});
        const response_content = await promise.json()
        type_content = response_content;
    }

    onMount(updateItemsContent);

    const checkFormCompletness = () => {
        if (item_id != undefined) {
            console.log("doing nothing:", is_form_ready);
            return is_form_ready; // do nothing
        }

        let current_element;
        for(let field of type_data.fields) {
            current_element = document.getElementById(`ft-${field["field-name"]}`);
            if (current_element.value === "") {
                return false;
            }
        }

        return true;
    }

    const clear = () => {
        let current_element
        for(let field of type_data.fields) {
            current_element = document.getElementById(`ft-${field['field-name']}`);
            if (current_element === null) {
                console.warn(`no element with id 'ft-${field['field-name']}'`);
                continue;
            }
            current_element.value = field.type !== "select" ? "" : "0";
        }
        item_id = undefined;
        is_form_ready = false;
    }

    const deleteItem = async () => {
        if (item_id != undefined) {
            const headers = new Headers();
            headers.set("X-sk", getSessionKey());

            const request = new Request(`${server_name}/${type_name}?id=${item_id}`, {method:"DELETE", headers: headers});

            await fetch(request)
                .then(promise => {
                    if (promise.status === 403) {
                        alert("muppets like you dont have permission to do this action, now for trying you are soooo fired...");
                    }
                });
            clear();
            updateItemsContent();
        }
    }

    const getFormContent = () => {
        let current_element;
        const fields_content = [];
        for(let field of type_data.fields) {
            current_element = document.getElementById(`ft-${field["field-name"]}`);
            if (current_element === null) {
                console.log(`ft-${field["field-name"]}`)
            }
            fields_content.push({
                name: field["field-name"],
                content: current_element.value
            });
        }

        if(type_data.extras !== undefined) {
            const { extras } = type_data;
            extras.forEach(fe => {
                current_element = document.getElementById(`ft-${fe.name}`);
                if (current_element === null) {
                    alert(`element 'fe-${fe.name}' was null`);
                }
                fe.value = parseInt(current_element.value);
            });
            fields_content.push({
                name: "extras",
                content: JSON.stringify(extras)
            });
        }
        return fields_content;
    }

    const triggerCloseCallback = e => {
        if (onClose != undefined && e.target === e.currentTarget) {
            onClose();
        }
    }

    const submitForm = () => {
        if(checkFormCompletness()) {
            const fields_content= getFormContent();
            console.log(fields_content);
            const headers = new Headers();
            headers.set("X-sk", getSessionKey())
            let forma = new FormData();
            forma.append("type_name", type_name);
            fields_content.forEach(f => {
                forma.append(f.name, f.content);
            });

            const request = new Request(`${server_name}/${type_name}`, {body: forma, method: "POST", headers:headers});
            fetch(request)
                .then(promise => {
                    if (promise.ok) {
                        console.log("user-registred");
                        updateItemsContent();
                    } else if (promise.status === 403) {
                        alert("you are not allowed to create users");
                    } else {
                        console.warn("server error");
                    }
                }).catch(reason => console.warn(reason));
            clear()
        }
    } 

    const setItemSelected = item_data => {
        const item_keys = Object.keys(item_data);
        let current_element
        for(let k of item_keys) {
            if (k === "id") {
                continue;
            }
            if (k === "extras") {
                for(let extra_field of item_data.extras) {
                    current_element = document.getElementById(`ft-${extra_field.name}`)
                    current_element.value = extra_field.value
                    current_element.max = extra_field.value + current_element.max
                }
                continue;
            }
            current_element = document.getElementById(`ft-${k}`);
            if (current_element === null) {
                alert(`missing key '${k}'`);
                continue;
            }
            current_element.value = item_data[k];
        
        }
        item_id = item_data.id;
        is_form_ready = true;
        
    }

    const searchItem = search_value => {
        const headers = new Headers();
        headers.set("X-sk", getSessionKey())
        fetch(`${server_name}/search?type_name=${type_name}&value=${search_value}`, {method: "GET", headers: headers})
        .then(promise => {
            if (promise.ok) {
                promise.json().then(response => type_content = response);
            } else if (promise.status == 404) {
                console.warn(`no ${search_value} found in ${type_name}`);
            }
        })
    }

    const requestUpdate = () => {
        const item_data = getFormContent();
        const forma = new FormData();
        forma.append("id", item_id);
        for(let f of item_data) {
            forma.append(f.name, f.content);
        }
        const headers = new Headers();
        headers.set("X-sk", getSessionKey());

        const request = new Request(`${server_name}/${type_name}`, {method: "PATCH", headers: headers, body:forma});
        
        fetch(request)
            .then(promise => {
                if (promise.status === 403) {
                    alert("muppets like you dont have permission to do this action, now for trying you are soooo fired...");
                }
            })
        clear();
        updateItemsContent();
    }

    const requestItemData = e => {
        const item_uuid = e.currentTarget.getAttribute("uuid");
        const headers = new Headers();
        headers.set("X-sk", getSessionKey());

        const request = new Request(`${server_name}/${type_name}?id=${item_uuid}`, {method: "GET", headers: headers})
        fetch(request)
            .then(promise => {
                if (promise.ok) {
                    promise.json().then(setItemSelected)
                }
            })
    }

    const waitForEnter = e => {
        if(e.key.toLowerCase() === "enter") {
            e.target.blur();
            if (item_id == undefined) {
                is_form_ready = checkFormCompletness();
            }
        }
    }

</script>

<style>

    
    /*=============================================
    =            Main Component            =
    =============================================*/
    
    #data-management-modal-background {
        position: fixed;
        display: flex;
        height: 100vh;
        width: 100vw;
        justify-content: center;
        z-index: 3;
    }

    #data-management-modal {
        position: fixed;
        display: flex;
        height: 100vh;
        width: 80%;
        background: var(--dark-color);
        flex-direction: column;
        align-items: center;
        box-shadow: 0 0 10px 10px #000000ae;
        padding-top: 3%;
    }

    
    /*=============================================
    =            Controls            =
    =============================================*/
    

    #dmm-controls {
        width: 40%;
        display: flex;
        justify-content: space-around;
        margin: 1.5% auto;
    }

    
    /*=============================================
    =            type-fields            =
    =============================================*/
    
    #type-fields-container {
        display: flex;
        flex-wrap: wrap;
        justify-content: space-evenly;
    }

    :global(#type-fields-container div) {
        margin-top: 1%;
    }
    
    
    /*=============================================
    =            type items            =
    =============================================*/
    
    #type-items-container {
        width: 70%;
        margin-top: 6%;
    }

    #type-name {
        padding: .5% 1%;
        font-size: 1.2rem;
        font-weight: bold;
        text-transform: capitalize;
        border-bottom: 2px solid var(--theme-color);
    }
    
    #items-container {
        overflow-y: auto;
        height: 30vh;

    }

    .item-container {
        cursor: pointer;
        display: flex;
        width: 90%;
        height: 10vh;
        background: hsla(97, 77%, 60%, 0.589);
        color: white;
        justify-content: space-around;
        align-items: center;
        margin: 1vh auto 0;
    }

    .item-container:hover {
        filter: brightness(1.2);
    }

    .item-data1 {
        font-size: 1.3rem;
    }

</style>

<div on:click={triggerCloseCallback} id="data-management-modal-background">
    <div id="data-management-modal">
        <SearchBar name={type_name} searchCallback={searchItem}/>
        <div id="dmm-controls">
            <Button label="Create" isEnabled={is_form_ready} onClick={submitForm}/>
            <Button label="Update" isEnabled={is_form_ready} onClick={requestUpdate}/>
            <Button label="Delete" isEnabled={is_form_ready} onClick={deleteItem} isDanger={true}/>
        </div>
        <div id="type-fields-container">
            {#each type_data.fields as tf}
                {#if tf.type !== "select"}                    
                    <Input 
                        input_id={`ft-${tf['field-name']}`}
                        input_type={tf.type}
                        input_placeholder={`${tf['field-name']}...`}
                        onKeypressed={waitForEnter}
                        onBlur={() => is_form_ready = checkFormCompletness()}
                    />
                {:else}
                    <Select options={tf.options} name={tf['field-name']} />
                {/if} 
            {/each}
            {#if type_data.extras !== undefined}
            <!-- checking if the type includes extras, services should include it -->
                {#each type_data.extras as tfe}
                    <Input 
                        input_id={`ft-${tfe['name']}`}
                        input_type="number"
                        input_label={`${tfe['name']}`}
                        initial_value=0
                        min={tfe.min}
                        max={tfe.max}
                        onKeypressed={waitForEnter}
                        onBlur={() => is_form_ready = checkFormCompletness()}
                    />
                {/each}
            {/if}
        </div>
        <div id="type-items-container">
            <div id="type-name">
                {type_name}
            </div>
            <div id="items-container">
                {#each type_content as c}
                    <div on:click={requestItemData} uuid={c.id} class="item-container">
                        <span class="item-data1">{c.username}:</span>
                        <span class="item-data2">{c.name}</span>
                    </div>
                {/each}
            </div>
        </div>
    </div>
</div>
