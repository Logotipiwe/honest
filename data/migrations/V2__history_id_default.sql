alter table questions_history
    alter column id set default (UUID());
