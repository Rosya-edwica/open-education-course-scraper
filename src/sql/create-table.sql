create table course(
    Название_курса text not null,
    Описание_курса text not null,
    Сертификат boolean,
    Продолжительность_в_неделях int,
    Ссылка_на_страницу text,
    Дата_начала text,
    Дата_окончания text,
    Картинка_курса text,
    Цена integer,
    Навыки text,
    Что_нужно_знать text,
    Фото_преподавателя text,
    ФИО_преподавателя text,
    О_преподавателе text,
    Программа  text,

    primary key (Ссылка_на_страницу)
);