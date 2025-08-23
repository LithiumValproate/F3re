# ruff: noqa: F401
"""
A collection of modules for the Sophomore project.

This package exports all the core data models, contact information types,
and constants used throughout the academic-related services.
"""

from .constants import (
    CHINA_PROVINCES,
    DayOfWeek,
    Repetition,
    Sex,
    Status,
)
from .contact import (
    Email,
    Phone,
)
from .models import (
    Address,
    Course,
    FamilyMember,
    Grade,
    Student,
    Teacher,
    TimeSlot,
)

__all__ = [
    # from constants.py
    "CHINA_PROVINCES",
    "DayOfWeek",
    "Repetition",
    "Sex",
    "Status",
    # from contact.py
    "Email",
    "Phone",
    # from models.py
    "Address",
    "Course",
    "FamilyMember",
    "Grade",
    "Student",
    "Teacher",
    "TimeSlot",
]
