CREATE TABLE public.deck (
    id uuid primary key,
    shuffled boolean not null
    );

CREATE TABLE public.card (
    id uuid primary key,
    deck_id uuid not null,
    value varchar(16) not null,
    suit varchar(16) not null,
    isDrawn boolean null,
    CONSTRAINT fk_deck_id
      FOREIGN KEY(deck_id) 
      REFERENCES public.deck(id)
);
