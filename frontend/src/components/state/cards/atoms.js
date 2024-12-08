import { atom } from "recoil";

export const cardsState = atom({
    key: "cardsState",
    default: [],
});

export const cardModalState = atom({
    key: "cardModalState",
    default: {
        id: "",
        profileName: "",
        name: "",
        number: "",
        expMonth: "",
        expYear: "",
        cvv: "",
    },
});