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
            sessionStorage.setItem("token", token)
        },
        tokenFromCache(state){
            const token = sessionStorage.getItem("token")
            if (token){
                console.log('token from session storage')
            }
            state.token = token
        }
    }
}