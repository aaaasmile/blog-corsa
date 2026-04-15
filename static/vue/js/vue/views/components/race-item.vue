<template>
<v-card class="mt-4">
    <v-card-title>New Race</v-card-title>
    <v-card-text>
      <v-form>
        <v-text-field
          v-model="item.Name"
          label="Name"
          outlined
          dense
        ></v-text-field>
        <v-text-field
          v-model="item.Title"
          label="Title"
          outlined
          dense
        ></v-text-field>
        <v-text-field
          v-model="item.Distance"
          label="Distance"
          outlined
          dense
        ></v-text-field>
        <v-menu
          v-model="dateMenu"
          :close-on-content-click="false"
          :nudge-right="40"
          transition="scale-transition"
          offset-y
          min-width="auto"
        >
          <template v-slot:activator="{ on, attrs }">
            <v-text-field
              v-model="item.Date"
              label="Date"
              prepend-icon="mdi-calendar"
              readonly
              outlined
              dense
              v-bind="attrs"
              v-on="on"
            ></v-text-field>
          </template>
          <v-date-picker
            v-model="item.Date"
            @input="dateMenu = false"
          ></v-date-picker>
        </v-menu>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn text @click="cancel">Cancel</v-btn>
      <v-btn color="primary" @click="save">Save</v-btn>
    </v-card-actions>
  </v-card>
  </template>