from dataclasses import is_dataclass, asdict
import datetime as dt
from enum import Enum
import sqlite3
import sqlite3 as sql
from typing import Optional, get_type_hints

from ..model.contact import Email, Phone
from ..model.models import Student
from . import utils


STUDENT_DB = 'students.sqlite'


class StudentStore:
    def __init__(self, db_path: str = STUDENT_DB):
        self._db_path = db_path
        self._conn: Optional[sqlite3.Connection] = None

    def __enter__(self):
        self._conn = sql.connect(self._db_path)
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        if self._conn:
            self._conn.close()

    def create_table(self):
        if not self._conn:
            raise ConnectionError('Cannot create table')
        c = self._conn.cursor()
        c.execute('''
        CREATE TABLE IF NOT EXISTS students_profile (
            student_id INTEGER PRIMARY KEY,
            student_name VARCHAR(255) NOT NULL,
            sex VARCHAR(10) NOT NULL,
            birthdate DATE NOT NULL,
            phone INTEGER NOT NULL,
            email VARCHAR(255),
            status VARCHAR(20) NOT NULL
        )
        ''')
        c.execute('''
        CREATE TABLE IF NOT EXISTS majors (
            
        )
        ''')
        c.execute('''
        CREATE TABLE IF NOT EXISTS courses (
        )
        ''')
        c.execute('''
        CREATE TABLE IF NOT EXISTS students_academic (
            student_id INTEGER PRIMARY KEY,
            enroll_year INTEGER NOT NULL,
            major_id INTEGER NOT NULL,
            major_name VARCHAR(255) NOT NULL,
            class_id INTERGER NOT NULL,
            FOREIGN KEY (student_id) REFERENCES students_profile(student_id)
        )
        ''')