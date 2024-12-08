import { selector } from "recoil";
import { emailsInputState} from "./atoms";

export const emailsInputStateSelector = selector({
    key: "emailsInputStateSelector",
    get: ({ get }) => {
        const emails = get(emailsInputState);
        return emails;
    },
});