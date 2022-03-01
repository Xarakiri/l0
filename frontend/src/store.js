import axios from "axios";
import { createStore } from "vuex";

const BACKEND_URL = "http://localhost:8080"
const CREATE_ORDER = "CREATE_ORDER"
const GET_ORDER = "GET_ORDER"
const GET_SUCCESS = "GET_SUCCESS"
const GET_ERROR = "GET_ERROR"

const store = createStore({
    state () {
        return {
            order: '',
            getResult: '',
        }
    },
    mutations: {
        [CREATE_ORDER](state, order) {
            state.order = order;
        },
        [GET_ORDER](state, order) {
            state.getResult = order;
        },
        [GET_SUCCESS](state, order) {
            state.getResult = order;
        },
        [GET_ERROR](state) {
            state.getResult = '';
        },
    },
    actions: {
        getOrder({commit }, data) {
            if (data.length === 0) {
                commit(GET_SUCCESS, '');
                return;
            }
            axios
                .get(`${BACKEND_URL}/orders/${data}`)
                .then(({ data }) => commit(GET_SUCCESS, data))
                .catch((err) => {
                    console.error(err);
                    commit(GET_ERROR);
                });
        },
        async createOrder({ commit }, order) {
            await axios.post(`${BACKEND_URL}/orders`, order).
            then((result) => {
                console.log(result.data);
                console.log(result);
                commit(CREATE_ORDER, result.data);
            }).catch((error) => {
                console.log(error)
            })
        },
    },
});

export default store;