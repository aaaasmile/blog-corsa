<template>
  <v-card>
    <v-card-title class="subheading grey--text">Login</v-card-title>
    <v-divider></v-divider>
    <v-container>
      <v-row justify="center">
        <v-col cols="12" sm="6" md="4">
          <v-text-field v-model="User" label="User"></v-text-field>
          <v-text-field v-model="Password" type="password" label="Password"></v-text-field>
        </v-col>
      </v-row>
    </v-container>
    <v-card-actions class="justify-center pb-4">
      <v-btn color="primary" v-on:click="doLogin">Login</v-btn>
    </v-card-actions>
  </v-card>
</template>