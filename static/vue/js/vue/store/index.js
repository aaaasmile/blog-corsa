import Generic from './generic-store.js?version=105'
import Admin from './admin-store.js?version=110'

export default new Vuex.Store({
  modules: {
    gen: Generic,
    admin: Admin,
  }
})
