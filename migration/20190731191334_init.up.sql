create table StudentCourse
(
    semester_id          smallint    not null,
    student_id           varchar(16) not null,
    course_by_teacher_id bigint      not null,
    constraint StudentCourse_pk
        unique (semester_id, student_id, course_by_teacher_id)
);
