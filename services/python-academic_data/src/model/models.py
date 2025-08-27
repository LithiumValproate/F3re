from __future__ import annotations
from dataclasses import dataclass, field
import datetime as dt

from . import constants as cst
from .contact import Phone, Email


@dataclass(frozen=True)
class Address:
    province: str
    city: str

    def __post_init__(self):
        if not isinstance(self.province, str) or not isinstance(self.city, str):
            raise TypeError('Province and city must be strings.')
        if self.province not in cst.CHINA_PROVINCES:
            raise ValueError(f"Invalid province: {self.province}.")


@dataclass(frozen=True)
class FamilyMember:
    name: str
    relationship: str
    phone: Phone



@dataclass(frozen=True)
class TimeSlot:
    day: cst.DayOfWeek
    start_time: dt.time
    end_time: dt.time
    repetition: cst.Repetition

    def __post_init__(self):
        if self.start_time >= self.end_time:
            raise ValueError("Start time must be before end time.")
        if not isinstance(self.repetition, cst.Repetition):
            raise TypeError("Repetition must be an instance of Repetition Enum.")


@dataclass(eq=True, unsafe_hash=True)
class Teacher:
    teacher_id: int
    name: str
    sex: cst.Sex
    department: str
    phone: Phone
    email: Email
    courses: set[Course] = field(default_factory=set, repr=False, hash=False, compare=False)


@dataclass(frozen=True)
class Course:
    course_id: int
    name: str
    teacher: Teacher
    location: str
    credit: int
    class_id: frozenset[int]
    time_slots: tuple[TimeSlot, ...]

    def __post_init__(self):
        if not 0 < self.credit <= 5:
            raise ValueError("Credit must be between 1 and 5.")
        if not self.time_slots:
            raise ValueError("At least one time slot must be provided.")


@dataclass(frozen=True)
class Grade:
    course: Course
    score: float

    @property
    def grade_point(self) -> float:
        return round(self.score / 20, 1)

    @property
    def quality_point(self) -> float:
        return self.grade_point * self.course.credit

    def __post_init__(self):
        if not isinstance(self.course, Course):
            raise TypeError("Course must be an instance of Course.")
        if not (0 <= self.score <= 100):
            raise ValueError("Score must be between 0 and 100.")


@dataclass
class Student:
    student_id: int
    name: str
    sex: cst.Sex
    birthdate: dt.date
    enroll_year: int
    major: tuple[int, str]
    class_id: int
    phone: Phone
    email: Email
    address: Address
    family_members: list[FamilyMember]
    status: cst.Status
    grades: list[Grade] = field(default_factory=list, repr=False)

    @property
    def age(self) -> int:
        today = dt.date.today()
        return today.year - self.birthdate.year - (
                (today.month, today.day) < (self.birthdate.month, self.birthdate.day))

    @property
    def gpa(self) -> float:
        if not self.grades:
            return 0.0
        total_quality_points = total_credits = 0.0
        for grade in self.grades:
            total_quality_points += grade.quality_point
            total_credits += grade.course.credit
        return round(total_quality_points / total_credits, 1) if total_credits > 0 else 0.0

    def __post_init__(self):
        if self.enroll_year > dt.date.today().year:
            raise ValueError("Enroll year cannot be in the future.")
        if not (1900 <= self.birthdate.year <= dt.date.today().year):
            raise ValueError(f"Birth year {self.birthdate.year} must be between 1900 and {dt.date.today().year}.")
