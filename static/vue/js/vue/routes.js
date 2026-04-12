import Dashboard from './views/Dashboard.js?version=107'
import Login from './views/Login.js?version=100'
import Races from './views/Races.js?version=100'

export default [
  { path: '/', icon: 'dashboard', title: 'Dashboard', component: Dashboard },
  { path: '/login', icon: 'folder', title: 'Login', component: Login },
  { path: '/races', icon: 'mdi-flag-checkered', title: 'Races', component: Races },
]
