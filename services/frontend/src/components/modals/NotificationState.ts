import { reactive } from "vue";

export const notificationState = reactive({
  enabled: false,
  message: "default message",
  color: "#fff"
});

export function CreateNotification(message: string) {
  notificationState.enabled = true;
  notificationState.color = "#feb2b2";
  notificationState.message = message;
}
