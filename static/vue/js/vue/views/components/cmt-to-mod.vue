<template>
  <v-container>
    <v-card>
      <v-card-subtitle>Comments to moderate</v-card-subtitle>

      <v-divider></v-divider>

      <v-list
        lines="three"
        select-strategy="leaf"
      >
        <v-list-item
          v-for="cmt in Comments"
          :key="cmt.Id"
          :subtitle="cmt.PostId"
          :title="cmt.Comment"
        >
          <template v-slot:append>
            <v-btn
              color="grey-lighten-1"
              icon="mdi-information"
              variant="text"
            ></v-btn>
          </template>
        </v-list-item>
      </v-list>
    </v-card>
  </v-container>
</template>