import { atom } from "recoil";

export const emailsInputState = atom({
    key: "emailsInputState",
    default: {
		emailsString: "",
	},
});
