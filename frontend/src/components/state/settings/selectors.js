import { selector } from "recoil";
import { settingsState } from "./atoms";

export const settingsStateSelector = selector({
    key: "settingsStateSelector",
    get: ({ get }) => {
        const settings = get(settingsState);
        return settings;
    },
});