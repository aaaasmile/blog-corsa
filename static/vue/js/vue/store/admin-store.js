export default {
    state: {
        server_name: '',
        loadingSync: false,
        token: '',
    },
    mutations: {
        serviceServer(state, servername) {
            state.server_name = servername
        },
        storeToken(state, token) {
            state.token = token
        },
        tokenFromCache(state){
            
        }
    }
}