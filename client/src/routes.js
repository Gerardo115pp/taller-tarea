import Login from './pages/Login.svelte';
import Home from './pages/Home.svelte';

const routes = {
    // Exact path
    '/': Login,
    '/home/:session': Home,
}

export { routes }