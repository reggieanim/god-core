export class Notify {
  constructor() {
  }
  sendNotification = async (instruction) => {
    await chrome.runtime.sendMessage({
      action: "notify",
      title: "My Approval Extension",
      message: instruction.value
    });
  };
}
