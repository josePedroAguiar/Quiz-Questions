CREATE TABLE tblquizzes (
	id		 BIGINT ,
	created_time	 TIMESTAMP NOT NULL DEFAULT NOW(),
	name		 TEXT,
	tblusers_id BIGINT NOT NULL,
	PRIMARY KEY(id)
);

CREATE TABLE tblquizzes_tblquestions (
	tblquizzes_id	 BIGINT,
	questions_id BIGINT,
	PRIMARY KEY(tblquizzes_id,questions_id),
    question TEXT,
    description TEXT,
    answer_a TEXT,
    answer_b TEXT,
    answer_c TEXT,
    answer_d TEXT,
    answer_e TEXT,
    answer_f TEXT,
    multiple_correct_answers BOOLEAN,
    answer_a_correct BOOLEAN,
    answer_b_correct BOOLEAN,
    answer_c_correct BOOLEAN,
    answer_d_correct BOOLEAN,
    answer_e_correct BOOLEAN,
    answer_f_correct BOOLEAN,
    correct_answer TEXT,
    explanation TEXT,
    tip TEXT,
    category TEXT,
    difficulty TEXT
);

CREATE TABLE tblquizzes_tblusers (
	tblquizzes_id	 BIGINT,
	tblusers_id BIGINT,
	PRIMARY KEY(tblquizzes_id,tblusers_id)
);


ALTER TABLE tblquizzes_tblquestions ADD CONSTRAINT tblquizzes_tblquestions_fk1 FOREIGN KEY (tblquizzes_id) REFERENCES tblquizzes(id);
ALTER TABLE tblquizzes_tblusers ADD CONSTRAINT tblquizzes_tblusers_fk1 FOREIGN KEY (tblquizzes_id) REFERENCES tblquizzes(id);
