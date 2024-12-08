export type Card = {
   id: string;
   profileName: string;
   name: string;
   number: string;
   expMonth: string;
   expYear: string;
   cvv: string;
};

export type CardColumn = {
   name: string;
   uid: string;
};

export type CardTableProps = {
   columns: CardColumn[];
   cards: Card[];
};

export const cardColumns: CardColumn[] = [
   {name: 'PROFILE', uid: 'profile'},
   {name: 'NAME', uid: 'name'},
   {name: 'NUMBER', uid: 'number'},
   {name: 'EXPIRATION', uid: 'expiration'},
   {name: 'ACTIONS', uid: 'actions'},
];