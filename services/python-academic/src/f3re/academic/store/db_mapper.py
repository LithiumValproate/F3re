import csv
from dataclasses import is_dataclass, asdict
import datetime as dt
from enum import Enum
import sqlite3 as sql
from typing import Optional, get_type_hints

from ..model.contact import Email, Phone
from ..model.models import Student
from . import utils

STUDENT_DB = 'academic.sqlite'
MAJOR_TABLE = '../static/major.csv'
COURSE_TABLE = '../static/course.csv'


class StudentStore:
    def __init__(self, db_path: str = STUDENT_DB):
        self._db_path = db_path
        self.conn: Optional[sql.Connection] = None

    def __enter__(self):
        self.conn = sql.connect(self._db_path)
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        if self.conn:
            self.conn.close()

    def create_table(self):
        if not self.conn:
            raise ConnectionError('Cannot create table')
        c = self.conn.cursor()
        c.execute('''
                  create table Students_Profile (
                      Student_Id   integer primary key,
                      Student_Name varchar(255) not null,
                      Sex          varchar(10)  not null,
                      Birthdate    date         not null,
                      Age          integer      not null,
                      Phone        integer      not null,
                      Email        varchar(255),
                      Status       varchar(20)  not null
                      );
                  ''')
        c.execute('''
                  create table Majors (
                      Major_Id   integer primary key,
                      Major_Name varchar(255) not null
                      )
                  ''')
        c.execute('''
                  create table Courses (
                      Course_Id integer primary key,
                      Course_Name varchar(255) not null,
                  )
                  ''')
        c.execute('''
                  create table Students_Academic (
                      Student_Id  integer primary key,
                      Enroll_Year integer      not null,
                      Major_Id    integer      not null,
                      Major_Name  varchar(255) not null,
                      Class_Id    integer      not null,
                      Course_Id   integer      not null,
                      Course_Name varchar(255) not null,
                      Score       integer      not null,
                      foreign key (Student_Id) references Students_Profile (Student_Id),
                      foreign key (Major_Id) references Majors (Major_Id),
                      foreign key (Course_Id) references Courses (Course_Id) 
                      )
                  ''')
        self.conn.commit()

    def create(self, stu: Student):
        if not self.conn:
            raise ConnectionError('Cannot create student')
        c = self.conn.cursor()
        try:
            c.execute('''
                      insert into Students_Profile (Student_Id, Student_Name, Sex, Birthdate, Phone, Email, Status) value (?, ?, ?, ?, ?, ?, ?)
                      ''', (stu.student_id, stu.name, stu.sex, stu.birthdate, stu.phone, stu.email, stu.status))
            self.conn.commit()
        except Exception as e:
            self.conn.rollback()
            raise e

    def find_by_id(self, student_id: int) -> Optional[Student]:
        if not self.conn:
            raise ConnectionError('Cannot find student')
        c = self.conn.cursor()
        c.execute('''
                  select * from Students_Profile where Student_Id = ?
                  ''', (student_id,))
        row = c.fetchone()
        if not row:
            raise sql.OperationalError('Student not found at profile table')
        dct = {
            'student_id': row[0],
            'name': row[1],
            'sex': row[2],
            'birthdate': row[3],
            'phone': row[5],
            'email': row[6],
            'status': row[7]
        }
        c.execute('''
                  select * from Students_Academic where Student_Id = ?
                  ''', (student_id,))
        row = c.fetchone()
        if not row:
            raise sql.OperationalError('Student not found at academic table')
        dct.update({'enroll_year': row[1], 'major': (row[2], row[3]), 'class_id': row[4]})

