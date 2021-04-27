<script>
    import Option from '../components/Option.svelte';
    import DataManagementModal from '../components/DataManagementeModal.svelte';
    import { server_name } from '../server_data';

    export let params = {};
    let type_name = "";
    let type_data = [];
    let show_management_modal = false;

    const updateType = type_label => {
        if (type_name != type_label) {
            const headers = new Headers();
            headers.set("X-sk", String(params.session));
    
            const request = new Request(`${server_name}/type?type_name=${type_label}`, {method: "GET", headers: headers});
            fetch(request,)
                .then(promise => {
                    if (promise.status == 200) {
                        promise.json().then( response => {
                            type_name = type_label;
                            type_data = response;
                            show_management_modal = true;
                        });
                    }
                });
        } else {
            show_management_modal = true;
        }
    }


    const closeModal = () =>  show_management_modal = false;
</script>

<style>

    #home-page-component {
        height: 100vh;
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
    }

    #home-header {
        margin-bottom: 20%;
        user-select: none;
    }

    header {
        width: 14%;
        padding: 1rem;
        border-bottom: 2px solid var(--theme-color);
        text-align: center;
    }

    #home-hub {
        display: flex;
        width: 80%;
        flex-direction: column;
    }

    .options-level-container{
        cursor: pointer;
        display: flex;
        justify-content: space-around;
        margin-top: 3%;
    }
</style>

<div id="home-page-component">
    {#if show_management_modal}
        <DataManagementModal onClose={closeModal} getSessionKey={() => params.session} {type_data} {type_name}/>
    {/if}
    <header id="home-header">
        <span>User&Services</span>
    </header>
    <div id="home-hub">
        <div class="options-level-container">
            <Option onClick={updateType} label="users"/>
            <Option onClick={updateType} label="clients"/>
            <Option label="services"/>
        </div>
        <div class="options-level-container">
            <Option label="cars"/>
            <Option label="warehouse"/>
        </div>
    </div>
</div>