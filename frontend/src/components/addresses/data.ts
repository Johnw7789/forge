export type Address = {
   id: string;
   profileName: string;
   name: string;
   line1: string;
   line2: string;
   city: string;
   state: string;
   zip: string;
   phone: string;
};

export type AddressColumn = {
   name: string;
   uid: string;
};

export type AddressTableProps = {
   columns: AddressColumn[];
   addresses: Address[];
};

export const addressColumns: AddressColumn[] = [
   {name: 'PROFILE', uid: 'profile'},
   {name: 'NAME', uid: 'name'},
   {name: 'STREET', uid: 'street'},
   {name: 'PHONE', uid: 'phone'},
   {name: 'ACTIONS', uid: 'actions'},
];