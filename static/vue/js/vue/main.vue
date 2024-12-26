<template>
  <v-app class="grey lighten-4">
    <Navbar />
    <v-content class="mx-4 mb-4">
      <router-view></router-view>
    </v-content>
    <v-footer absolute>
      <v-col class="text-center caption" cols="12">
        {{ new Date().getFullYear() }} —
        <span>Buildnr: {{Buildnr}} - Server: {{ServiceServer}}</span>
      </v-col>
    </v-footer>
  </v-app>
</template>