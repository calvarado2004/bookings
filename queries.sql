select count(id) from room_restrictions where '2022-07-18' < end_date and '2022-07-20' > start _date;



select count(id) from room_restrictions where '2022-09-20' < end_date and '2022-09-21' > start _date;



select 
    r.id, r.room_name
from
    rooms r
where
    r.id not in (select rr.room_id from room_restrictions rr where '2021-09-20' < rr.end_date and '2021-09-20' > rr.start_date )