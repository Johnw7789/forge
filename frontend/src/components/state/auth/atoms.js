import { atom, selector } from "recoil"

export const authState = atom({
    key: "authState",
    default: {
		loading: false,
		authenticated: false,
		licenseKey: "",
		msg: "",
		discordUser: "",
		discordImage: "",
	},
})