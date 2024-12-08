import { atom } from "recoil";

export const addressesState = atom({
    key: "addressesState",
    default: [],
});

export const addressModalState = atom({
    key: "addressModalState",
    default: {
        id: "",
        profileName: "",
        name: "",
        line1: "",
        line2: "",
        city: "",
        state: "",
        zip: "",
        phone: "",
    },
});