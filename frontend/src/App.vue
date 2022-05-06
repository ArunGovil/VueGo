<script>
export default {
  name: "App",
  data() {
    return {
      version: "",
    };
  },

  async created() {
    const response = await fetch("http://localhost:8080/api/version");
    const result = await response.text();
    this.version = result;
  },
};
</script>

<template>
  <div>
    <img alt="Vue logo" src="./assets/logo.png" />
    <p>{{ version }}</p>
  </div>
</template>


<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
