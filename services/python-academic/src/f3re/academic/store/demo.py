import datetime as dt
import json
import os
import sqlite3

from ..model import models, constants, contact
from .json_mapper import EnhancedJSONEncoder, EnhancedJSONDecoder

DB_FILE = 'student_demo.sqlite'


def create_tables(conn):
    """Create database tables for students."""
    c = conn.cursor()
    # A simplified student table. Complex fields are stored as JSON.
    c.execute('''
        CREATE TABLE IF NOT EXISTS students (
            student_id INTEGER PRIMARY KEY,
            name TEXT NOT NULL,
            sex TEXT NOT NULL,
            birthdate TEXT NOT NULL,
            enroll_year INTEGER NOT NULL,
            major TEXT NOT NULL,
            class_id INTEGER NOT NULL,
            phone TEXT NOT NULL,
            email TEXT NOT NULL,
            address TEXT NOT NULL,
            family_members TEXT NOT NULL,
            status TEXT NOT NULL,
            grades TEXT NOT NULL
        )
    ''')
    conn.commit()


def insert_student(conn, student: models.Student):
    """Inserts a student object into the database."""
    c = conn.cursor()

    # Serialize complex objects to JSON using the corrected JSON mapper
    major_json = json.dumps(student.major, cls=EnhancedJSONEncoder)
    phone_json = json.dumps(student.phone, cls=EnhancedJSONEncoder)
    email_json = json.dumps(student.email, cls=EnhancedJSONEncoder)
    address_json = json.dumps(student.address, cls=EnhancedJSONEncoder)
    family_members_json = json.dumps(student.family_members, cls=EnhancedJSONEncoder)
    grades_json = json.dumps(student.grades, cls=EnhancedJSONEncoder)

    c.execute(
        '''
        INSERT INTO students (
            student_id, name, sex, birthdate, enroll_year, major, class_id,
            phone, email, address, family_members, status, grades
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ''',
        (
            student.student_id,
            student.name,
            student.sex.name,
            student.birthdate.isoformat(),
            student.enroll_year,
            major_json,
            student.class_id,
            phone_json,
            email_json,
            address_json,
            family_members_json,
            student.status.name,
            grades_json,
        )
    )
    conn.commit()


def get_student_by_id(conn, student_id: int) -> models.Student | None:
    """Retrieves a student from the database by student_id."""
    c = conn.cursor()
    c.execute('SELECT * FROM students WHERE student_id = ?', (student_id,))
    row = c.fetchone()

    if row:
        (
            student_id, name, sex, birthdate, enroll_year, major_json, class_id,
            phone_json, email_json, address_json, family_members_json, status, grades_json
        ) = row

        # Deserialize JSON fields and reconstruct the Student object
        student_data = {
            'student_id': student_id,
            'name': name,
            'sex': constants.Sex[sex],
            'birthdate': dt.date.fromisoformat(birthdate),
            'enroll_year': enroll_year,
            'major': tuple(json.loads(major_json)),
            'class_id': class_id,
            'phone': json.loads(phone_json, cls=EnhancedJSONDecoder),
            'email': json.loads(email_json, cls=EnhancedJSONDecoder),
            'address': json.loads(address_json, cls=EnhancedJSONDecoder),
            'family_members': json.loads(family_members_json, cls=EnhancedJSONDecoder),
            'status': constants.Status[status],
            'grades': json.loads(grades_json, cls=EnhancedJSONDecoder)
        }
        return models.Student(**student_data)
    return None


def main():
    """Demonstration of database operations for students."""
    # Create a sample student
    teacher = models.Teacher(
        teacher_id=1,
        name="Dr. Smith",
        sex=constants.Sex.MALE,
        department="Computer Science",
        phone=contact.Phone("13800138000"),
        email=contact.Email("smith@example.com")
    )
    course = models.Course(
        course_id=101,
        name="Introduction to Programming",
        teacher=teacher,
        location="Building A, Room 101",
        credit=3,
        class_id=frozenset([1, 2]),
        time_slots=(
            models.TimeSlot(
                day=constants.DayOfWeek.MONDAY,
                start_time=dt.time(9, 0),
                end_time=dt.time(11, 0),
                repetition=constants.Repetition.WEEKLY
            ),
        )
    )
    student = models.Student(
        student_id=2024001,
        name="Alice",
        sex=constants.Sex.FEMALE,
        birthdate=dt.date(2005, 5, 20),
        enroll_year=2024,
        major=(1, "Computer Science"),
        class_id=1,
        phone=contact.Phone("13912345678"),
        email=contact.Email("alice@example.com"),
        address=models.Address(province="Shanghai", city="Shanghai"),
        family_members=[
            models.FamilyMember(
                name="Bob",
                relationship="Father",
                phone=contact.Phone("13987654321")
            )
        ],
        status=constants.Status.ACTIVE,
        grades=[models.Grade(course=course, score=95.5)]
    )

    conn = sqlite3.connect(DB_FILE)

    try:
        create_tables(conn)
        print(f"Inserting student: {student.name}")
        insert_student(conn, student)
        print("Student inserted successfully.")

        print(f"Retrieving student with ID: {student.student_id}")
        retrieved_student = get_student_by_id(conn, student.student_id)

        if retrieved_student:
            print("Retrieved student successfully:")
            print(retrieved_student)
            assert student == retrieved_student
            print("Verification successful: Original and retrieved students are the same.")
        else:
            print("Failed to retrieve student.")

    finally:
        conn.close()
        if os.path.exists(DB_FILE):
            os.remove(DB_FILE)
            print(f"Cleaned up demo database file: {DB_FILE}")


if __name__ == "__main__":
    main()
