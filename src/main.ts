// import { AesEncryptUtil } from "./helpers/functions/aesEncrypt";
import { CreateNewWindowOrTab, getStartingURL } from "./helpers/functions/serviceWorker";

document.getElementById("startButton")!.addEventListener("click", async () => {
  try {
    const response = await fetch("/zz.json");
    const rawInstructions = await response.text();
    const startingUrl = getStartingURL(rawInstructions);
    await chrome.storage.session.set({ startingUrl: startingUrl });
    await chrome.storage.session.set({ instructions: rawInstructions });
    await CreateNewWindowOrTab(rawInstructions);
  } catch (error) {
    console.error("Error starting instruction processor:", error);
  }
});
