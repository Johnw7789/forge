export type Task = {
   id: string;
   name: string;
   email: string;
   password: string;
   proxy: string;
   status: string;
};

export type TaskColumn = {
   name: string;
   uid: string;
};

export type TaskTableProps = {
   columns: TaskColumn[];
   tasks: Task[];
};

export const taskColumns: TaskColumn[] = [
   {name: 'EMAIL', uid: 'email'},
   {name: 'PASSWORD', uid: 'password'},
   {name: 'PROXY', uid: 'proxy'},
   {name: 'STATUS', uid: 'status'},
   {name: 'ACTIONS', uid: 'actions'},
];