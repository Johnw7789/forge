import { selector } from "recoil";
import { tasksState } from "./atoms";

export const tasksStateSelector = selector({
    key: "tasksStateSelector",
    get: ({ get }) => {
        const tasks = get(tasksState);
        return tasks;
    },
});