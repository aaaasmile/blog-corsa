export default {
    state: {
        server_name: '',
        loadingSync: false,
        session_id: '',
    },
    mutations: {
        serviceServer(state, servername) {
            state.server_name = servername
        }
    }
}