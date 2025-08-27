import datetime as dt
import json
from dataclasses import is_dataclass, asdict
from enum import Enum

from academic_data import Student, Course, Teacher, Sex, Address, Grade, Status, Phone, Email
import utils


class EnhancedJSONEncoder(json.JSONEncoder):
    def default(self, o):
        if is_dataclass(o):
            if o.__class__.__name__ in utils.CLASS_REGISTRY:
                dct = asdict(o)
                dct['__type__'] = o.__class__.__name__
                return dct
        if isinstance(o, dt.datetime) or isinstance(o, dt.date) or isinstance(o, dt.time):
            return o.isoformat()
        if isinstance(o, Enum):
            return o.name.lower()
        if isinstance(o, frozenset):
            return list(o)
        return super().default(o)


def student_to_json(student: Student) -> str:
    """
    Args:
        student: Student object
    Returns:
        str: JSON representation of the student object
    """
    if not isinstance(student, Student):
        raise TypeError('Input must be a Student object.')
    return json.dumps(student, cls=EnhancedJSONEncoder, indent=2)


if __name__ == '__main__':
    teacher = Teacher(teacher_id=101, name='Dr. Turing', sex=Sex.MALE, department='Computer Science', phone=Phone('12345678901'), email=Email('turing@example.com'))
    course = Course(
        course_id=1,
        name='Introduction to CS',
        teacher=teacher,
        location='Room 101',
        credit=3,
        class_id=frozenset([1, 2]),
        time_slots=()
    )
    sample_student = Student(
        student_id=12345,
        name='Ada Lovelace',
        sex=Sex.FEMALE,
        birthdate=dt.date(1815, 12, 10),
        enroll_year=1833,
        major=(1, 'Computer Science'),
        class_id=1,
        phone=Phone('09876543210'),
        email=Email('ada@example.com'),
        address=Address(province='Beijing', city='Beijing'),
        family_members=[],
        status=Status.GRADUATED,
        grades=[Grade(course=course, score=95.5)]
    )

    json_output = student_to_json(sample_student)
    print(json_output)
