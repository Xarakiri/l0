<template>
    <div class="container py-5">
        <div class="row mb-4">
            <h1 class="col-12">Order Service</h1>
        </div>
        <div class="row">
            <MyGet class="col" />
            <MySearch class="col" />
        </div>
    </div>
</template>

<script>
import MyGet from '@/components/MyGet';
import MySearch from '@/components/MySearch';

export default {
    components: {
      MyGet,
      MySearch,
    },
};
</script>

<style lang="scss" scoped>
.container {
    max-width: 768px;
}
</style>
