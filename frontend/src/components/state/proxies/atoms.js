import { atom } from "recoil";

export const proxiesInputState = atom({
    key: "proxiesInputState",
    default: {
		proxiesString: "",
	},
});
