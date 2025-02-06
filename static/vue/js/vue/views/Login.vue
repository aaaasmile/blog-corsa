<template>
  <v-card>
    <v-card-title class="subheading grey--text">Login</v-card-title>
    <v-divider></v-divider>
    <v-container>
      <v-row>
        <v-text-field v-model="User" label="User"></v-text-field>
      </v-row>
      <v-row>
        <v-text-field v-model="Password" type="password" label="Password"></v-text-field>
      </v-row>
    </v-container>
    <v-card-actions>
      <v-btn color="primary" v-on:click="doLogin">Login</v-btn>
    </v-card-actions>
  </v-card>
</template>