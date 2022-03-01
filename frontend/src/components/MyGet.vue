<template>
    <div>
        <form v-on:submit.prevent="createOrder">
            <div class="input-group">
                <input v-model.trim="orderBody" type="text" class="form-control" placeholder="Send order to nats...">
                <div class="input-group-append">
                    <button class="btn btn-primary" type="submit">Send to NATS</button>
                </div>
            </div>
        </form>

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
import { mapState } from 'vuex';

export default {
    data() {
        return {
            orderBody: '',
        };
    },
    computed: mapState({
        orders: (state) => state.order,
    }),
    methods: {
        createOrder() {
            if (this.orderBody.length != 0) {
                this.$store.dispatch('createOrder', this.orderBody);
                this.orderBody = '';
            }
        },
    },
    components: {
    },
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