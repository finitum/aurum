<template>
  <div
    class="fixed z-10 inset-0 overflow-y-auto flex justify-center items-center"
  >
    <div class="fixed inset-0 transition-opacity">
      <div class="absolute inset-0 bg-gray-500 opacity-75"></div>
    </div>
    <div class="bg-white z-20 text-center p-2 my-2 shadow-md rounded">
      <h3 class="text-2xl font-medium">Login</h3>
      <form
        class="flex flex-col items-center align-middle justify-between text-center"
        @submit.prevent="ClickLogin"
      >
        <span class="label">
          <label for="username">Username:</label>
          <input
            v-model="user.username"
            id="username"
            type="text"
            class="input-box"
          />
        </span>
        <span class="flex flex-col label">
          <label for="password">Password:</label>
          <input
            v-model="user.password"
            id="password"
            type="password"
            class="input-box"
          />
        </span>
        <button class="button w-full" type="submit">
          Login
        </button>

        <span class="mt-3 text-gray-600">
          Don't have an account?
          <a class="text-blue-800 cursor-pointer" @click="router.push('/register')">
            Register
          </a>  instead
        </span>

      </form>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, reactive } from "vue";
import router from "../router";
import { client } from "@/client/client";
import { ErrorCode, User } from "aurum-client";
import { CreateNotification } from "@/components/modals/NotificationState";

export default defineComponent({
  name: "Login",
  setup() {
    const user = reactive<User>({ username: "", password: "", email: "" });

    async function ClickLogin() {
      if (user.username === "") {
        CreateNotification("Username field empty");
      } else if (user.password === "") {
        CreateNotification("Password field empty");
      } else {
        const error = await client.Login(user);
        if (error.isOk()) {
          console.log(client.Verify())

          await router.push("/");
        } else if (error.error.Code === ErrorCode.Unauthorized) {
          CreateNotification("Username or password incorrect");
        } else if (error.error.Code === ErrorCode.WeakPassword) {
          CreateNotification(error.error.Message);
        }
      }
    }

    return {
      user,
      ClickLogin,
      router,
    };
  }
});
</script>

<style scoped lang="postcss">
.label {
  @apply flex flex-col text-gray-700 font-bold p-2 m-1 w-full;
}

.button {
  @apply flex-shrink-0 bg-primary text-sm text-white py-2 px-3 rounded;
}

.button:hover {
  @apply bg-primarydark;
}

.input-box {
  @apply bg-gray-200 appearance-none border-2 border-gray-200 rounded py-2 px-4 text-gray-700 leading-tight;
}

.input-box:focus {
  @apply outline-none bg-white border-indigo-500;
}
</style>
