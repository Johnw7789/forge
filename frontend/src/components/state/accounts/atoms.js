import { atom } from "recoil";

export const accountsState = atom({
    key: "accountsState",
    default: [],
});

export const selectedAccountsState = atom({
    key: "selectedAccountsState",
    default: [],
});

export const accountModalState = atom({
    key: "accountModalState",
    default: {
        id: "",
        name: "",
        email: "",
        password: "",
        phone: "",
        proxy: "",
        key2fa: "",
        arc: "",
        ard: "",
        cookies: "",
        prime: false,
        status: "Idle",
    },
});

// export type Account = {
//     id: string;
//     name: string;
//     email: string;
//     password: string;
//     phone: string;
//     proxy: string;
//     key2fa: string;
//     arc: string;
//     ard: string;
//     cookies: string;
//     prime: boolean;
//     status: string;
//  };