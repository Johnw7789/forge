import { atom } from "recoil";

export const settingsState = atom({
    key: "settingsState",
    default: {
        licenseKey: "",
        maxTasks: 3,
        limitProxyUse: true,
        persistState: true,
        nameOverride: "",
        webhooks: {
            success: "",
            fail: "",
        },
        imapConfig: {
            uniqueTaskClient: true,
            username: "",
            password: ""
        },
        smsConfig: {
            maxTries: 2,
            provider: "SMS Man",
            username: "",
            apiKey: ""
        },
        captchaKey: "",
        captchaMaxTries: 2,
        icloudConfig: {
            username: "",
            password: ""
        },
        icloudCookies: "",
        localHost: false
	},
});