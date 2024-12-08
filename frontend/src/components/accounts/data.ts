export type Account = {
   id: string;
   name: string;
   email: string;
   password: string;
   phone: string;
   proxy: string;
   key2fa: string;
   arc: string;
   ard: string;
   cookies: string;
   prime: boolean;
   status: string;
};

export type AccountColumn = {
   name: string;
   uid: string;
};

export type AccountTableProps = {
   columns: AccountColumn[];
   accounts: Account[];
};

export const accountColumns: AccountColumn[] = [
   {name: 'NAME', uid: 'name'},
   {name: 'EMAIL', uid: 'email'},
   // {name: 'PHONE', uid: 'phone'},
   {name: 'PASSWORD', uid: 'password'},
   {name: 'PROXY', uid: 'proxy'},
   // {name: 'PRIME', uid: 'prime'},
   {name: 'STATUS', uid: 'status'},
   {name: 'ACTIONS', uid: 'actions'},
];