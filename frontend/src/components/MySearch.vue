<template>
  <div>
    <input
      @keyup="getOrder"
      v-model.trim="query"
      class="form-control"
      placeholder="Get order..."
    />
    <div class="mt-4">
      <div class="card">
          <div class="card-body">
            <p class="card-text">
              <pre>{{orders}}</pre>
            </p>
          </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  data() {
    return {
      query: "",
    };
  },
  computed: mapState({
    orders: (state) => state.getResult,
  }),
  methods: {
    getOrder() {
      if (this.query != this.lastQuery) {
        this.$store.dispatch("getOrder", this.query);
        this.lastQuery = this.query;
      }
    },
  },
  components: {},
};
</script>

<style lang="scss" scoped>
.card {
  margin-bottom: 1rem;
}
.card-body {
  padding: 0.5rem;
  p {
    margin-bottom: 0;
  }
}
</style>