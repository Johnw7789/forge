import { selector } from "recoil";
import { proxiesInputState} from "./atoms";

export const proxiesInputStateSelector = selector({
    key: "proxiesInputStateSelector",
    get: ({ get }) => {
        const proxies = get(proxiesInputState);
        return proxies;
    },
});