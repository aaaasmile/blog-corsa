import Generic from './generic-store.js?version=105'
import Cmt from './cmt-store.js?version=110'

export default new Vuex.Store({
  modules: {
    gen: Generic,
    cmt: Cmt,
  }
})
